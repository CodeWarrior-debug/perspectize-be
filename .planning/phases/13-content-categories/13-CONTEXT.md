# Phase 13: Content Categories — Context

## Phase Goal

Enable content categorization using Google NL taxonomy. Add Claude-powered auto-suggestion.

## Problem Statement

From FEATURE_BACKLOG.md:

"Categorize content using Google's Cloud NL Content Taxonomy — 27 top-level categories designed for digital content classification."

Currently content has no categorization beyond `content_type` (which is just 'youtube_video'). Users cannot browse or filter by topic.

## Research Summary

See `.planning/v1.1-research/CONTENT-CATEGORIZATION.md` for full research.

**Taxonomy decision:** Google NL taxonomy (20 of 27 categories)
- Arts & Entertainment, Autos & Vehicles, Beauty & Fitness, Books & Literature
- Business & Industrial, Computers & Electronics, Finance, Food & Drink
- Games, Health, Hobbies & Leisure, Home & Garden
- Jobs & Education, Law & Government, News, People & Society
- Pets & Animals, Science, Sports, Travel & Transportation

**Dropped for v1:** Adult, Internet & Telecom, Online Communities, Real Estate, Reference, Sensitive Subjects, Shopping

**Auto-categorization:** Claude Haiku (5x cheaper than Google Cloud NL)
- Cost: $0.031/month for 100 videos vs $0.14/month for Google Cloud NL
- Accuracy: Comparable for classification tasks
- Flexibility: Custom taxonomy, no API quota limits

**Storage decision:** Lookup table with ltree (not enum)
- Supports hierarchical categories (future)
- Allows user-defined categories (future)
- Easier to add/modify categories

## Database Changes

```sql
-- Categories lookup table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    path LTREE NOT NULL,  -- For future hierarchy
    description TEXT,
    icon TEXT,            -- Emoji or icon class
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed initial categories
INSERT INTO categories (name, slug, path, icon, display_order) VALUES
('Arts & Entertainment', 'arts-entertainment', 'arts_entertainment', '🎭', 1),
('Autos & Vehicles', 'autos-vehicles', 'autos_vehicles', '🚗', 2),
('Beauty & Fitness', 'beauty-fitness', 'beauty_fitness', '💄', 3),
-- ... etc

-- Add category to content
ALTER TABLE content ADD COLUMN category_id INT REFERENCES categories(id);
CREATE INDEX idx_content_category_id ON content (category_id);
```

## GraphQL Schema Changes

```graphql
type Category {
  id: IntID!
  name: String!
  slug: String!
  icon: String
  contentCount: Int!
}

type Content {
  # ... existing fields
  category: Category
}

input CreateContentFromYouTubeInput {
  url: String!
  categoryId: IntID  # Optional, can auto-suggest
}

input UpdateContentInput {
  id: IntID!
  categoryId: IntID
}

type CategorySuggestion {
  category: Category!
  confidence: Float!  # 0-1
  reasoning: String
}

type Query {
  categories: [Category!]!
  suggestCategory(contentId: IntID!): CategorySuggestion!
}
```

## Claude Integration Pattern

```go
// Category suggestion service
type CategorySuggestionService struct {
    client *anthropic.Client
    repo   ports.CategoryRepository
}

func (s *CategorySuggestionService) SuggestCategory(ctx context.Context, content *domain.Content) (*CategorySuggestion, error) {
    categories, _ := s.repo.ListAll(ctx)

    prompt := fmt.Sprintf(`Categorize this YouTube video into one of these categories:
%s

Video title: %s
Video description: %s

Respond with JSON: {"category_slug": "...", "confidence": 0.0-1.0, "reasoning": "..."}`,
        formatCategories(categories),
        content.Name,
        content.Description,
    )

    resp, err := s.client.Messages.Create(ctx, anthropic.MessageCreateParams{
        Model:     anthropic.ModelClaudeHaiku,
        MaxTokens: 200,
        Messages: []anthropic.Message{{
            Role:    "user",
            Content: prompt,
        }},
    })

    // Parse JSON response
    // ...
}
```

## Frontend Components

```svelte
<!-- CategoryPicker.svelte -->
<script lang="ts">
    import { createQuery } from '@tanstack/svelte-query';
    import { Select } from '$lib/shadcn';

    let { value = $bindable(), onSuggest } = $props<{
        value?: number;
        onSuggest?: () => void;
    }>();

    const categories = createQuery(() => ({
        queryKey: queryKeys.categories.list(),
        queryFn: () => graphqlClient.request(LIST_CATEGORIES),
    }));
</script>

<Select bind:value>
    <SelectTrigger>
        <SelectValue placeholder="Select category" />
    </SelectTrigger>
    <SelectContent>
        {#each categories.data?.categories ?? [] as category}
            <SelectItem value={category.id}>
                <span class="mr-2">{category.icon}</span>
                {category.name}
            </SelectItem>
        {/each}
    </SelectContent>
</Select>

{#if onSuggest}
    <Button variant="outline" size="sm" on:click={onSuggest}>
        ✨ Suggest
    </Button>
{/if}
```

## Requirements Covered

- CAT-01 through CAT-08 (all category requirements)

## Success Metrics

| Metric | Current | Target |
|--------|---------|--------|
| Content with category | 0% | 100% (new) |
| Category count | 0 | 20 |
| Auto-suggestion accuracy | N/A | > 80% |
| Filter by category | No | Yes |

## Dependencies

- Phase 12 (auth for user-specific categories later)
- Claude API key (Anthropic account required)

## Risks

- **Claude API costs:** Monitor usage, set budget alerts
- **Category accuracy:** Users may disagree with suggestions
- **Schema migration:** Content backfill needs default or NULL handling

## Open Questions

1. Should uncategorized content be allowed, or require category selection?
2. Should category suggestion be automatic on add, or user-triggered?
3. Should we track category suggestion accuracy for improvement?

---

*Context gathered: 2026-02-16*
