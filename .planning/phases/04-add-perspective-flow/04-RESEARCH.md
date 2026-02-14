# Phase 4: Add Perspective Flow - Research

**Researched:** 2026-02-07
**Domain:** Form management with TanStack Form Svelte, multi-field validation, progress bar visualization
**Confidence:** MEDIUM

## Summary

Phase 4 implements a perspective creation form with four rating inputs (Quality, Agreement, Importance, Confidence), Like text, Review text, and video selection. The core challenge is integrating TanStack Form's Svelte adapter (first use in the project) with the existing shadcn-svelte components, GraphQL mutations, and toast notifications.

**Key technical areas:**
1. **TanStack Form Svelte adapter** - Form state management with `createForm` and field composition pattern using Svelte 5 snippets
2. **Rating inputs** - Number inputs (0-10000 range) paired with shadcn Progress component for visualization
3. **Video selector** - shadcn Combobox or Select component connected to existing content query
4. **Form validation** - Field-level validation with error display via toast notifications before submission
5. **User attribution** - Integration with Phase 2's user selection store (`getSelectedUserId()`)

**Primary recommendation:** Use TanStack Form with field-level validation, shadcn Progress component for rating visualization, shadcn Select/Combobox for video selection, and integrate with existing toast system for error feedback. Follow the established Svelte 5 runes patterns and type-based component organization.

## Standard Stack

The established libraries/tools for this domain:

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| @tanstack/svelte-form | 1.28.0 | Form state management | Already installed in Phase 1, type-safe, Svelte 5 compatible, field-level validation |
| shadcn-svelte Progress | latest | Progress bar visualization | Matches existing design system, accessible, animates smoothly |
| svelte-sonner | 1.0.7 | Toast notifications | Already configured (top-right, 2s auto-dismiss) for validation errors |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| shadcn-svelte Select | latest | Video dropdown selector | Simple video selection (< 50 videos) |
| shadcn-svelte Combobox | latest | Video search+select | Searchable video selection (if > 50 videos or filtering needed) |
| shadcn-svelte Textarea | latest | Review text input | Multi-line freeform text entry |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| TanStack Form | Svelte native bindings + manual validation | TanStack provides built-in validation timing, error state management, and field composition — manual approach is more verbose |
| shadcn Progress | Custom SVG progress bar | shadcn Progress is accessible, animated, and matches theme — custom is more work |
| Combobox | Native `<select>` with `<datalist>` | Native datalist has poor cross-browser support and styling limitations |

**Installation:**
```bash
# Progress component
pnpm dlx shadcn-svelte@latest add progress

# Select component (for simple dropdown)
pnpm dlx shadcn-svelte@latest add select

# Combobox component (for searchable dropdown)
pnpm dlx shadcn-svelte@latest add combobox

# Textarea component (for Review field)
pnpm dlx shadcn-svelte@latest add textarea

# TanStack Form already installed (Phase 1)
```

## Architecture Patterns

### Recommended Project Structure
```
frontend/src/lib/
├── components/
│   ├── AddPerspectiveDialog.svelte    # Main form dialog
│   ├── RatingInput.svelte             # Reusable rating input + progress bar
│   ├── VideoSelector.svelte           # Video selection dropdown/combobox
│   └── shadcn/
│       ├── progress/                  # shadcn Progress component
│       ├── select/                    # shadcn Select component
│       ├── combobox/                  # shadcn Combobox component (if needed)
│       └── textarea/                  # shadcn Textarea component
├── queries/
│   └── perspectives.ts                # CREATE_PERSPECTIVE mutation
└── utils/
    └── validation.ts                  # Form validation helpers (if needed)
```

### Pattern 1: TanStack Form with Svelte 5 Snippets
**What:** Use `createForm` to manage form state, `form.Field` component with Svelte 5 `{#snippet children(field)}` for field rendering
**When to use:** All forms in the project (establishes pattern for future forms)
**Example:**
```svelte
<script lang="ts">
  import { createForm } from '@tanstack/svelte-form';
  import { toast } from 'svelte-sonner';

  const form = createForm(() => ({
    defaultValues: {
      claim: '',
      quality: null,
      agreement: null,
      importance: null,
      confidence: null,
      like: '',
      review: '',
      contentID: null,
    },
    onSubmit: async ({ value }) => {
      // Validate that required fields are set
      if (!value.claim || value.quality === null) {
        toast.error('Please fill in all required fields');
        return;
      }

      // Submit mutation
      try {
        await mutation.mutateAsync(value);
        toast.success('Perspective created successfully');
      } catch (error) {
        toast.error('Failed to create perspective');
      }
    },
  }));
</script>

<form onsubmit={(e) => {
  e.preventDefault();
  e.stopPropagation();
  form.handleSubmit();
}}>
  <form.Field name="claim">
    {#snippet children(field)}
      <label for={field.name}>Claim *</label>
      <input
        id={field.name}
        name={field.name}
        value={field.state.value}
        onblur={field.handleBlur}
        oninput={(e) => field.handleChange(e.target.value)}
      />
      {#if field.state.meta.errors.length > 0}
        <span class="text-destructive text-sm">{field.state.meta.errors[0]}</span>
      {/if}
    {/snippet}
  </form.Field>

  <button type="submit">Submit</button>
</form>
```
**Source:** [TanStack Form Svelte Quick Start](https://tanstack.com/form/v1/docs/framework/svelte/quick-start)

### Pattern 2: Rating Input with Progress Visualization
**What:** Pair number input (0-10000) with shadcn Progress component to visualize rating value
**When to use:** All four rating fields (Quality, Agreement, Importance, Confidence)
**Example:**
```svelte
<script lang="ts">
  import { Progress } from '$lib/components/shadcn/progress';

  let { value = $bindable(null), label, name, field } = $props<{
    value: number | null;
    label: string;
    name: string;
    field?: any; // TanStack Form field API
  }>();

  const MAX_RATING = 10000;

  function handleInput(event: Event) {
    const target = event.target as HTMLInputElement;
    const numValue = parseInt(target.value, 10);
    if (Number.isNaN(numValue) || numValue < 0) {
      value = 0;
    } else if (numValue > MAX_RATING) {
      value = MAX_RATING;
    } else {
      value = numValue;
    }
    field?.handleChange(value);
  }
</script>

<div class="space-y-2">
  <label for={name} class="text-sm font-medium">{label}</label>
  <input
    id={name}
    type="number"
    min="0"
    max={MAX_RATING}
    value={value ?? 0}
    oninput={handleInput}
    onblur={field?.handleBlur}
    class="h-9 w-24 rounded-md border border-input bg-background px-3 text-sm"
  />
  <Progress value={value ?? 0} max={MAX_RATING} class="h-2" />
  <p class="text-xs text-muted-foreground">
    {value ?? 0} / {MAX_RATING}
  </p>
</div>
```
**Source:** [shadcn-svelte Progress component](https://shadcn-svelte.com/docs/components/progress)

### Pattern 3: Video Selector with Existing Content Query
**What:** Use TanStack Query to fetch content list, render in shadcn Select or Combobox
**When to use:** Video selection field in perspective form
**Example:**
```svelte
<script lang="ts">
  import { createQuery } from '@tanstack/svelte-query';
  import { graphqlClient } from '$lib/queries/client';
  import { LIST_CONTENT } from '$lib/queries/content';

  let { value = $bindable(null), field } = $props<{
    value: number | null;
    field?: any;
  }>();

  const contentQuery = createQuery(() => ({
    queryKey: ['content'],
    queryFn: () => graphqlClient.request(LIST_CONTENT),
    staleTime: 5 * 60 * 1000,
  }));

  function handleChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    const newValue = target.value ? parseInt(target.value, 10) : null;
    value = newValue;
    field?.handleChange(newValue);
  }
</script>

{#if contentQuery.isLoading}
  <select disabled class="h-9 rounded-md border border-input bg-background px-3 text-sm">
    <option>Loading videos...</option>
  </select>
{:else if contentQuery.data}
  <select
    value={value ? String(value) : ''}
    onchange={handleChange}
    onblur={field?.handleBlur}
    class="h-9 rounded-md border border-input bg-background px-3 text-sm"
  >
    <option value="">Select video...</option>
    {#each contentQuery.data.content.items as content}
      <option value={content.id}>{content.name}</option>
    {/each}
  </select>
{/if}
```
**Source:** Adapted from UserSelector.svelte (Phase 2 pattern)

### Pattern 4: Form Validation with Toast Notifications
**What:** Display validation errors via toast notifications BEFORE submission, inline errors during typing
**When to use:** All form validation (required fields, field constraints)
**Example:**
```svelte
<script lang="ts">
  import { createForm } from '@tanstack/svelte-form';
  import { toast } from 'svelte-sonner';

  const form = createForm(() => ({
    defaultValues: { /* ... */ },
    onSubmit: async ({ value }) => {
      // Pre-submission validation
      const errors: string[] = [];

      if (!value.claim?.trim()) {
        errors.push('Claim is required');
      }
      if (value.quality === null || value.quality < 0 || value.quality > 10000) {
        errors.push('Quality rating must be between 0 and 10000');
      }
      // ... other validations

      if (errors.length > 0) {
        errors.forEach(err => toast.error(err));
        return;
      }

      // Submit mutation
      // ...
    },
  }));

  // Field-level validation (optional, for inline errors)
  const validateClaim = (value: string) => {
    if (!value?.trim()) return 'Claim is required';
    if (value.length < 3) return 'Claim must be at least 3 characters';
    return undefined;
  };
</script>

<form.Field name="claim" validators={{ onChange: validateClaim }}>
  {#snippet children(field)}
    <input
      value={field.state.value}
      oninput={(e) => field.handleChange(e.target.value)}
      onblur={field.handleBlur}
    />
    {#if field.state.meta.errors.length > 0}
      <span class="text-destructive text-sm">{field.state.meta.errors[0]}</span>
    {/if}
  {/snippet}
</form.Field>
```
**Source:** [Form validation best practices](https://www.nngroup.com/articles/errors-forms-design-guidelines/) and [TanStack Form validation guide](https://tanstack.com/form/v1/docs/framework/react/guides/validation)

### Pattern 5: Dialog State Reset on Open
**What:** Reset form state when dialog opens to prevent stale data
**When to use:** All dialog-based forms
**Example:**
```svelte
<script lang="ts">
  import { Dialog } from '$lib/components/shadcn/dialog';

  let { open = $bindable(false) } = $props();

  const form = createForm(() => ({ /* ... */ }));

  // Reset form when dialog opens
  $effect(() => {
    if (open) {
      form.reset();
    }
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <form onsubmit={/* ... */}>
      <!-- Form fields -->
    </form>
  </Dialog.Content>
</Dialog.Root>
```
**Source:** Phase 3 research (dialog pattern established)

### Anti-Patterns to Avoid
- **Using `$effect()` for derived values** — Use `$derived()` instead (Svelte 5 convention)
- **Manually managing form state with `$state()`** — Use TanStack Form's `createForm` for consistency
- **Validating empty required fields on input** — Only validate on blur/submit to avoid interrupting user flow
- **Disabling submit without pending state** — Always check `mutation.isPending` to prevent double-submission
- **Showing technical error messages to users** — Map GraphQL errors to user-friendly messages

## Don't Hand-Roll

Problems that look simple but have existing solutions:

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Form state management | Custom `$state()` with manual validation | `@tanstack/svelte-form` | Built-in validation timing (onChange, onBlur, onSubmit), error state management, field composition, type safety |
| Progress bar component | Custom SVG with transform calculations | `shadcn-svelte Progress` | Accessible, animated, themed, tested across browsers |
| Video search dropdown | Native `<select>` + manual filtering | `shadcn-svelte Combobox` | Keyboard navigation, accessible, popover positioning, focus management |
| Form validation errors | `console.log` or inline-only errors | Toast notifications + inline errors | Users expect both toast feedback (immediate) and inline errors (context) |
| Rating input constraints | Unvalidated number input | Number input with min/max + validation | Prevents invalid values, provides clear constraints, validates edge cases (NaN, negative, overflow) |

**Key insight:** Form UX is deceptively complex — validation timing, error display, focus management, accessibility, and edge cases multiply quickly. TanStack Form handles these concerns, and shadcn-svelte components provide accessible, themed UI primitives.

## Common Pitfalls

### Pitfall 1: TanStack Form Schema Validation Type System Gap
**What goes wrong:** Using Zod schema with TanStack Form results in validation running but errors not displaying
**Why it happens:** TanStack Form's error array is loosely typed as `unknown` to allow flexibility, leading to type system gaps when using schema validation libraries
**How to avoid:** Use field-level validators (functions returning `string | undefined`) instead of schema validators, or use custom form hooks with pre-bound validation
**Warning signs:** Field validates (console shows validation running) but `field.state.meta.errors` is empty or incorrectly typed

**Source:** [Avoiding TanStack Form Pitfalls](https://matthuggins.com/blog/posts/avoiding-tanstack-form-pitfalls)

### Pitfall 2: Form State Not Resetting Between Dialog Opens
**What goes wrong:** User opens dialog, enters data, closes dialog, reopens → old data still in form
**Why it happens:** Dialog visibility and form state are separate; hiding dialog doesn't reset `createForm` state
**How to avoid:** Call `form.reset()` in `$effect()` watching dialog `open` prop
**Warning signs:** Stale input values when reopening dialog after previous attempt

```svelte
$effect(() => {
  if (open) {
    form.reset();
  }
});
```

### Pitfall 3: Validation Fires on Every Keystroke (Poor UX)
**What goes wrong:** User types "Qua" in claim field, sees "Must be at least 5 characters" error while typing
**Why it happens:** Using `onChange` validator without considering UX — validation fires on every input event
**How to avoid:** Use `onBlur` for field-level validation (validates when user leaves field), reserve `onChange` for simple constraints (max length)
**Warning signs:** Users complain about "annoying" or "distracting" error messages while typing

**UX best practice:** Validate on blur for semantic errors (required, format), validate on change only for constraints (max length, character whitelist)

**Source:** [Form validation UX best practices](https://www.nngroup.com/articles/errors-forms-design-guidelines/)

### Pitfall 4: Number Input Accepts Non-Numeric Values
**What goes wrong:** User types "abc" in rating input, form state becomes NaN or empty string
**Why it happens:** HTML `<input type="number">` allows non-numeric input, returns empty string for invalid input
**How to avoid:** Parse input with `parseInt()` in `oninput` handler, clamp to min/max range, handle NaN explicitly
**Warning signs:** Form submits with NaN values, backend rejects request, user confused why rating "disappeared"

```svelte
function handleInput(event: Event) {
  const target = event.target as HTMLInputElement;
  const numValue = parseInt(target.value, 10);
  if (Number.isNaN(numValue) || numValue < 0) {
    value = 0;
  } else if (numValue > 10000) {
    value = 10000;
  } else {
    value = numValue;
  }
}
```

### Pitfall 5: GraphQL Mutation Errors Are Not User-Friendly
**What goes wrong:** Backend returns "Field validation error: userID is required", user sees technical error in toast
**Why it happens:** GraphQL-request throws with raw error message, not extracted/mapped to user-friendly text
**How to avoid:** Wrap mutation in try/catch, parse error message, map to user-friendly text, show generic fallback
**Warning signs:** Toast shows constraint violation errors, null pointer exceptions, or GraphQL operation names

```typescript
try {
  await mutation.mutateAsync(value);
  toast.success('Perspective created successfully');
} catch (error) {
  const message = error.message || '';
  if (message.includes('userID') || message.includes('user')) {
    toast.error('Please select a user before submitting');
  } else if (message.includes('claim')) {
    toast.error('Claim is required');
  } else {
    toast.error('Failed to create perspective. Please try again.');
  }
}
```

### Pitfall 6: Rating Input Range (0-10000) Is Not Usable
**What goes wrong:** User tries to set rating to 5000 using number input, typing "5000" character by character triggers validation/clamping, frustrating
**Why it happens:** Number input with strict validation can interfere with typing flow
**How to avoid:** Validate on blur, not on input; OR provide slider as alternative input method alongside number input
**Warning signs:** Users struggle to enter precise values, complaints about "jumpy" inputs

**UX recommendation:** Provide both number input (for precision) AND visual progress bar (for quick approximate values). Consider adding a slider if users request it.

## Code Examples

Verified patterns from official sources and existing codebase:

### Complete Perspective Form with TanStack Form
```svelte
<script lang="ts">
  import { createForm } from '@tanstack/svelte-form';
  import { createMutation, createQuery } from '@tanstack/svelte-query';
  import { toast } from 'svelte-sonner';
  import { graphqlClient } from '$lib/queries/client';
  import { CREATE_PERSPECTIVE } from '$lib/queries/perspectives';
  import { getSelectedUserId } from '$lib/stores/userSelection.svelte';
  import { Dialog } from '$lib/components/shadcn/dialog';
  import { Button } from '$lib/components/shadcn/button';
  import { Progress } from '$lib/components/shadcn/progress';

  let { open = $bindable(false) } = $props();

  const mutation = createMutation(() => ({
    mutationFn: (input: CreatePerspectiveInput) =>
      graphqlClient.request(CREATE_PERSPECTIVE, { input }),
    onSuccess: () => {
      toast.success('Perspective created successfully');
      open = false;
    },
    onError: (error) => {
      const message = error.message || '';
      if (message.includes('userID')) {
        toast.error('Please select a user before submitting');
      } else {
        toast.error('Failed to create perspective. Please try again.');
      }
    },
  }));

  const form = createForm(() => ({
    defaultValues: {
      claim: '',
      quality: null as number | null,
      agreement: null as number | null,
      importance: null as number | null,
      confidence: null as number | null,
      like: '',
      review: '',
      contentID: null as number | null,
    },
    onSubmit: async ({ value }) => {
      // Pre-submission validation
      const errors: string[] = [];
      const userID = getSelectedUserId();

      if (!userID) {
        errors.push('Please select a user');
      }
      if (!value.claim?.trim()) {
        errors.push('Claim is required');
      }
      if (value.quality === null) {
        errors.push('Quality rating is required');
      }
      if (value.agreement === null) {
        errors.push('Agreement rating is required');
      }
      if (value.importance === null) {
        errors.push('Importance rating is required');
      }
      if (value.confidence === null) {
        errors.push('Confidence rating is required');
      }

      if (errors.length > 0) {
        errors.forEach(err => toast.error(err));
        return;
      }

      // Submit with userID from store
      await mutation.mutateAsync({
        ...value,
        userID: userID!,
      });
    },
  }));

  // Reset form when dialog opens
  $effect(() => {
    if (open) {
      form.reset();
    }
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="max-w-2xl">
    <Dialog.Header>
      <Dialog.Title>Add Perspective</Dialog.Title>
    </Dialog.Header>

    <form onsubmit={(e) => {
      e.preventDefault();
      e.stopPropagation();
      form.handleSubmit();
    }} class="space-y-4">

      {/* Form fields would go here */}

      <Dialog.Footer>
        <Button type="button" variant="outline" onclick={() => open = false}>
          Cancel
        </Button>
        <Button type="submit" disabled={mutation.isPending}>
          {mutation.isPending ? 'Creating...' : 'Create Perspective'}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
```

### Reusable Rating Input Component
```svelte
<script lang="ts">
  import { Progress } from '$lib/components/shadcn/progress';

  let {
    label,
    name,
    value = $bindable(null),
    field = null,
    required = false,
  } = $props<{
    label: string;
    name: string;
    value: number | null;
    field?: any;
    required?: boolean;
  }>();

  const MAX_RATING = 10000;

  function handleInput(event: Event) {
    const target = event.target as HTMLInputElement;
    const numValue = parseInt(target.value, 10);

    if (Number.isNaN(numValue) || numValue < 0) {
      value = 0;
    } else if (numValue > MAX_RATING) {
      value = MAX_RATING;
    } else {
      value = numValue;
    }

    field?.handleChange(value);
  }

  function handleBlur() {
    field?.handleBlur();
  }
</script>

<div class="space-y-2">
  <label for={name} class="text-sm font-medium">
    {label}
    {#if required}<span class="text-destructive">*</span>{/if}
  </label>

  <div class="flex items-center gap-4">
    <input
      id={name}
      type="number"
      min="0"
      max={MAX_RATING}
      value={value ?? 0}
      oninput={handleInput}
      onblur={handleBlur}
      class="h-9 w-28 rounded-md border border-input bg-background px-3 text-sm"
    />

    <div class="flex-1">
      <Progress value={value ?? 0} max={MAX_RATING} class="h-2" />
      <p class="text-xs text-muted-foreground mt-1">
        {value ?? 0} / {MAX_RATING}
      </p>
    </div>
  </div>

  {#if field?.state.meta.errors.length > 0}
    <span class="text-destructive text-sm">{field.state.meta.errors[0]}</span>
  {/if}
</div>
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Manual form state with `$state()` | TanStack Form `createForm` | v1.0.0 (2024) | Standardized validation timing, error handling, field composition |
| Custom progress bar with inline SVG | shadcn-svelte Progress component | shadcn-svelte release | Accessible, animated, themed, less code |
| Native `<select>` with manual filtering | shadcn-svelte Combobox | shadcn-svelte v2 | Better keyboard nav, popover positioning, focus management |
| Inline validation errors only | Toast + inline errors | UX research consensus | Immediate feedback (toast) + context (inline) improves UX |

**Deprecated/outdated:**
- **Svelte 4 `on:` event syntax**: Use Svelte 5 `onevent={}` syntax (e.g., `onclick={}` not `on:click={}`)
- **Svelte 4 `<slot />`**: Use Svelte 5 snippets `{@render children()}` and `{#snippet children(field)}`
- **Form validation on every keystroke**: UX research shows validation on blur is better for semantic errors

## Open Questions

Things that couldn't be fully resolved:

1. **Should Video Selector use Select or Combobox?**
   - What we know: Select is simpler, Combobox adds search functionality
   - What's unclear: How many videos will users have in typical usage? If > 50, Combobox is better
   - Recommendation: Start with Select (simpler), migrate to Combobox if users request search (easy swap due to same shadcn-svelte API)

2. **Should rating inputs include a slider alongside number input?**
   - What we know: Number input alone may be frustrating for large range (0-10000)
   - What's unclear: Will users prefer slider for quick approximate values or number input for precision?
   - Recommendation: Start with number input + progress bar (read-only visualization). Add slider if user testing shows frustration

3. **How should Review field differ from Like field?**
   - What we know: Requirements say "Review text (design TBD)" — both are freeform text
   - What's unclear: Should Review be multi-line? Should it have max length? Should it support markdown?
   - Recommendation: Use `<textarea>` for Review (multi-line), `<input>` for Like (single-line). Defer markdown support to future phase

4. **Should validation errors show in toasts, inline, or both?**
   - What we know: Requirements say "validation error toasts before submission" (PERSP-08), UX research says inline + toast is best
   - What's unclear: Do we show ALL errors as toasts (overwhelming if 5+ errors) or just first error?
   - Recommendation: Show first 3 errors as toasts (2s auto-dismiss stagger), show all errors inline — balances feedback without overwhelming

## Sources

### Primary (HIGH confidence)
- TanStack Form Svelte docs - [Quick Start](https://tanstack.com/form/v1/docs/framework/svelte/quick-start), [Examples](https://tanstack.com/form/v1/docs/framework/svelte/examples/simple)
- shadcn-svelte Progress component - [Documentation](https://shadcn-svelte.com/docs/components/progress)
- UserSelector.svelte (Phase 2) - Established pattern for dropdown + TanStack Query
- Phase 3 research - Dialog reset pattern, GraphQL error handling

### Secondary (MEDIUM confidence)
- NN/g Form Validation Guidelines - [10 Design Guidelines for Reporting Errors in Forms](https://www.nngroup.com/articles/errors-forms-design-guidelines/)
- Form validation UX research - [Multiple sources](https://blog.logrocket.com/ux-design/ux-form-validation-inline-after-submission/) agree on blur validation + inline errors
- TanStack Form pitfalls - [Avoiding TanStack Form Pitfalls](https://matthuggins.com/blog/posts/avoiding-tanstack-form-pitfalls)

### Tertiary (LOW confidence)
- WebSearch: "Svelte 5 number input slider range" - Community discussion of slider libraries, needs verification with testing
- TanStack Form schema validation - GitHub issues mention type system gaps, but no official statement on recommended approach

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - TanStack Form already installed, shadcn-svelte Progress documented, patterns verified
- Architecture: HIGH - Patterns adapted from Phase 2 (UserSelector), Phase 3 (dialog), official TanStack docs
- Pitfalls: MEDIUM - Some pitfalls from blog posts (not official docs), but cross-verified with UX research and GitHub issues

**Research date:** 2026-02-07
**Valid until:** 2026-03-07 (30 days - stable libraries, but TanStack Form is actively developed)

**Critical for planner:**
- Phase 4 depends on Phase 2's user selection store (USER-03 requires `getSelectedUserId()`)
- TanStack Form is installed but NOT yet used in codebase — this is the first usage, establishes pattern
- Progress component needs to be installed via shadcn CLI before use
- Number input UX (0-10000) is a known challenge — consider slider in future if user testing shows issues
