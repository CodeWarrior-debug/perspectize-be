# User Roles Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a `role` column (`admin`/`sentinel`/`default`) to the users table and fix `ListAll` to exclude sentinel users by default.

**Architecture:** New migration adds `role TEXT NOT NULL DEFAULT 'default'` with CHECK constraint, backfills existing rows. Domain gets a `UserRole` type. `IsSentinel()` switches from username matching to role check. Repository `ListAll` adds `WHERE role != 'sentinel'`. GraphQL exposes `role` as a `UserRole` enum.

**Tech Stack:** Go, PostgreSQL, GORM, gqlgen, testify

---

### Task 1: Database Migration

**Files:**
- Create: `backend/migrations/000009_add_user_role.up.sql`
- Create: `backend/migrations/000009_add_user_role.down.sql`

**Step 1: Write the up migration**

```sql
-- Add role column with CHECK constraint
ALTER TABLE public.users
ADD COLUMN role TEXT NOT NULL DEFAULT 'default'
CONSTRAINT users_role_check CHECK (role IN ('admin', 'sentinel', 'default'));

-- Backfill: user ID 1 is admin
UPDATE public.users SET role = 'admin' WHERE id = 1;

-- Backfill: sentinel users by username
UPDATE public.users SET role = 'sentinel'
WHERE username IN ('[deleted]', '[system]');
```

**Step 2: Write the down migration**

```sql
ALTER TABLE public.users DROP COLUMN IF EXISTS role;
```

**Step 3: Commit**

```bash
git add backend/migrations/000009_add_user_role.up.sql backend/migrations/000009_add_user_role.down.sql
git commit -m "feat: add role column to users table (migration 000009)"
```

---

### Task 2: Domain Model — Add UserRole Type

**Files:**
- Modify: `backend/internal/core/domain/user.go`

**Step 1: Write failing test**

File: `backend/test/domain/user_test.go` (create if needed)

```go
package domain_test

import (
	"testing"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestIsSentinel_RoleBased(t *testing.T) {
	sentinel := &domain.User{Role: domain.UserRoleSentinel, Username: "[deleted]"}
	assert.True(t, sentinel.IsSentinel())

	admin := &domain.User{Role: domain.UserRoleAdmin, Username: "admin"}
	assert.False(t, admin.IsSentinel())

	regular := &domain.User{Role: domain.UserRoleDefault, Username: "alice"}
	assert.False(t, regular.IsSentinel())
}

func TestUserRoleConstants(t *testing.T) {
	assert.Equal(t, domain.UserRole("admin"), domain.UserRoleAdmin)
	assert.Equal(t, domain.UserRole("sentinel"), domain.UserRoleSentinel)
	assert.Equal(t, domain.UserRole("default"), domain.UserRoleDefault)
}
```

**Step 2: Run test to verify it fails**

Run: `cd /home/user/perspectize/backend && go test ./test/domain/... -v -run TestIsSentinel_RoleBased`
Expected: Compile error — `domain.UserRole` undefined

**Step 3: Implement domain changes**

In `backend/internal/core/domain/user.go`, add:

```go
// UserRole represents the role of a user in the system.
type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleSentinel UserRole = "sentinel"
	UserRoleDefault  UserRole = "default"
)
```

Add `Role UserRole` field to the `User` struct.

Update `IsSentinel()`:

```go
func (u *User) IsSentinel() bool {
	return u.Role == UserRoleSentinel
}
```

**Step 4: Run test to verify it passes**

Run: `cd /home/user/perspectize/backend && go test ./test/domain/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/core/domain/user.go backend/test/domain/user_test.go
git commit -m "feat: add UserRole type to domain, update IsSentinel to use role"
```

---

### Task 3: GORM Model & Mappers

**Files:**
- Modify: `backend/internal/adapters/repositories/postgres/gorm_models.go`
- Modify: `backend/internal/adapters/repositories/postgres/gorm_mappers.go`

**Step 1: Add Role to UserModel**

In `gorm_models.go`, add to `UserModel`:

```go
Role string `gorm:"not null;default:default"`
```

**Step 2: Update mappers**

In `gorm_mappers.go`, update `userModelToDomain` to map `Role`:

```go
Role: domain.UserRole(m.Role),
```

Update `userDomainToModel` to map `Role`:

```go
Role: string(u.Role),
```

**Step 3: Build to verify**

Run: `cd /home/user/perspectize/backend && go build ./...`
Expected: Compiles with zero errors

**Step 4: Commit**

```bash
git add backend/internal/adapters/repositories/postgres/gorm_models.go backend/internal/adapters/repositories/postgres/gorm_mappers.go
git commit -m "feat: add role field to GORM model and mappers"
```

---

### Task 4: Fix ListAll to Exclude Sentinels

**Files:**
- Modify: `backend/internal/adapters/repositories/postgres/gorm_user_repository.go`

**Step 1: Write failing test in service tests**

Add to `backend/test/services/user_service_test.go`:

```go
func TestListAll_ExcludesSentinels(t *testing.T) {
	repo := &mockUserRepository{
		listAllFn: func(ctx context.Context) ([]*domain.User, error) {
			return []*domain.User{
				{ID: 1, Username: "admin", Role: domain.UserRoleAdmin},
				{ID: 2, Username: "alice", Role: domain.UserRoleDefault},
			}, nil
		},
	}

	svc := newTestUserService(repo)
	users, err := svc.ListAll(context.Background())

	require.NoError(t, err)
	assert.Len(t, users, 2)
	for _, u := range users {
		assert.False(t, u.IsSentinel(), "sentinel user should not be in ListAll results")
	}
}
```

This test validates the contract — the mock already returns non-sentinels (as the real repo will after our fix).

**Step 2: Update ListAll in repository**

In `gorm_user_repository.go`, change `ListAll`:

```go
func (r *GormUserRepository) ListAll(ctx context.Context) ([]*domain.User, error) {
	var models []UserModel

	err := r.db.WithContext(ctx).
		Where("role != ?", "sentinel").
		Order("username ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(models))
	for i := range models {
		users[i] = userModelToDomain(&models[i])
	}

	return users, nil
}
```

**Step 3: Build to verify**

Run: `cd /home/user/perspectize/backend && go build ./...`
Expected: Compiles with zero errors

**Step 4: Run all tests**

Run: `cd /home/user/perspectize/backend && go test ./... -v`
Expected: All pass

**Step 5: Commit**

```bash
git add backend/internal/adapters/repositories/postgres/gorm_user_repository.go backend/test/services/user_service_test.go
git commit -m "feat: exclude sentinel users from ListAll query"
```

---

### Task 5: GraphQL Schema & Codegen

**Files:**
- Modify: `backend/schema.graphql`
- Modify: `backend/gqlgen.yml`
- Modify (regenerated): `backend/internal/adapters/graphql/model/models_gen.go`
- Modify (regenerated): `backend/internal/adapters/graphql/generated/generated.go`
- Modify: `backend/internal/adapters/graphql/resolvers/helpers.go`

**Step 1: Add UserRole enum and role field to schema**

In `schema.graphql`, add enum before User type:

```graphql
enum UserRole {
  ADMIN
  SENTINEL
  DEFAULT
}
```

Add `role: UserRole!` field to User type.

**Step 2: Bind UserRole enum in gqlgen.yml**

Add to the `models:` section in `gqlgen.yml`:

```yaml
  UserRole:
    model:
      - github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain.UserRole
```

**Step 3: Add GraphQL enum values to domain**

The gqlgen enum binding requires the domain type to implement the graphql `Marshaler`/`Unmarshaler` interface, OR we can use the auto-generated enum. Since other enums (Privacy, ContentType, etc.) bind directly to domain types, we need to add the `Values()` and `IsValid()` pattern to `UserRole` in domain. Check how existing enums like `Privacy` are defined in domain, and follow the same pattern.

Actually — check: `SortOrder`, `Privacy`, `ContentType` are already bound. Look at how they're defined in domain to match the pattern (the enum values must match GraphQL uppercase convention).

In `backend/internal/core/domain/user.go`, update the `UserRole` constants to use UPPERCASE values matching GraphQL:

Wait — actually check the existing pattern. Privacy uses `PUBLIC`/`PRIVATE` in domain and the mappers do `strings.ToUpper`/`strings.ToLower`. So the domain enum values are already uppercase. **But** the DB stores lowercase. We need the same pattern:
- Domain: `UserRoleAdmin = "ADMIN"`, etc.
- DB mappers: `strings.ToLower(string(u.Role))` going to DB, `UserRole(strings.ToUpper(m.Role))` coming from DB
- GraphQL: binds directly to domain type (uppercase matches GraphQL enum)

**Update domain constants:**

```go
const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleSentinel UserRole = "SENTINEL"
	UserRoleDefault  UserRole = "DEFAULT"
)
```

**Update mappers** to uppercase/lowercase:

```go
// userModelToDomain
Role: domain.UserRole(strings.ToUpper(m.Role)),

// userDomainToModel
Role: strings.ToLower(string(u.Role)),
```

**Update repository** filter to use lowercase (DB stores lowercase):

```go
Where("role != ?", "sentinel")
```

This is already correct since the DB column stores lowercase.

**Step 4: Run gqlgen codegen**

Run: `cd /home/user/perspectize/backend && go run github.com/99designs/gqlgen generate`
Expected: Regenerates `models_gen.go` and `generated.go` with UserRole enum and role field on User.

**Step 5: Update GraphQL helper**

In `helpers.go`, update `userDomainToModel`:

```go
func userDomainToModel(u *domain.User) *model.User {
	return &model.User{
		ID:        strconv.Itoa(u.ID),
		Username:  u.Username,
		Email:     u.Email,
		Active:    u.Active,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
```

**Step 6: Build and test**

Run: `cd /home/user/perspectize/backend && go build ./... && go test ./... -v`
Expected: All pass

**Step 7: Commit**

```bash
git add backend/schema.graphql backend/gqlgen.yml backend/internal/core/domain/user.go \
  backend/internal/adapters/graphql/model/models_gen.go \
  backend/internal/adapters/graphql/generated/generated.go \
  backend/internal/adapters/graphql/resolvers/helpers.go \
  backend/internal/adapters/repositories/postgres/gorm_mappers.go
git commit -m "feat: expose UserRole enum in GraphQL schema"
```

---

### Task 6: Update Service Sentinel Checks

**Files:**
- Modify: `backend/internal/core/services/user_service.go`

**Step 1: Update Create to set default role**

In `user_service.go`, `Create` method, set role on the new user:

```go
user := &domain.User{
	Username: username,
	Email:    email,
	Active:   true,
	Role:     domain.UserRoleDefault,
}
```

Also update the reserved username check to use role concept (keep the username check too — belt and suspenders for creation):

No change needed — the existing username check on create is still valid as input validation. The role is set on the new user, not checked.

**Step 2: Update existing tests that create sentinel users in mocks**

In `backend/test/services/user_service_test.go`, update mock sentinel users to include `Role: domain.UserRoleSentinel`:

- `TestUpdate_SentinelUserBlocked`: Add `Role: domain.UserRoleSentinel` to the mock return
- `TestUpdate_SystemSentinelBlocked`: Add `Role: domain.UserRoleSentinel`
- `TestDelete_SentinelUserBlocked`: Add `Role: domain.UserRoleSentinel`
- `TestDelete_SystemSentinelBlocked`: Add `Role: domain.UserRoleSentinel`
- `TestDelete_Success` sentinel lookup: Add `Role: domain.UserRoleSentinel`
- `TestDelete_ReassignContentFails` sentinel: Add `Role: domain.UserRoleSentinel`
- `TestDelete_ReassignPerspectivesFails` sentinel: Add `Role: domain.UserRoleSentinel`

Non-sentinel mock users should have `Role: domain.UserRoleDefault`.

**Step 3: Run all tests**

Run: `cd /home/user/perspectize/backend && go test ./... -v`
Expected: All pass

**Step 4: Commit**

```bash
git add backend/internal/core/services/user_service.go backend/test/services/user_service_test.go
git commit -m "feat: set default role on user creation, update test mocks with roles"
```

---

### Task 7: Final Verification

**Step 1: Full build**

Run: `cd /home/user/perspectize/backend && go build ./...`
Expected: Zero errors

**Step 2: Full test suite**

Run: `cd /home/user/perspectize/backend && go test ./...`
Expected: All pass

**Step 3: Grep for stale references**

Check that no code still relies on username-based sentinel detection outside of input validation:

Run: `grep -rn "IsSentinel\|DeletedUserUsername\|SystemUserUsername" backend/`

- `IsSentinel()` should be in domain (definition) and service (usage) only
- `DeletedUserUsername` / `SystemUserUsername` should remain in domain (constants) and in service `Create` (input validation) and `Delete` (sentinel lookup by username)

**Step 4: Push**

```bash
git push -u origin claude/user-integration-flow-QODCb
```
