# YouTube Category to Google NL Category Mapping

**Version:** 1.0 (initial draft — pending user review)
**Generated:** 2026-02-20
**Status:** Awaiting user verification

---

## Overview

YouTube uses a flat taxonomy of ~32 content categories (via the YouTube Data API v3 `videoCategories` resource). This document maps each YouTube category to the corresponding Google NL V2 taxonomy path(s) used in Perspectize.

This mapping serves two purposes:
1. **Fallback classification** — When Google NL API confidence < 0.5 for a video, fall back to the YouTube category ID mapping
2. **Bulk import reference** — When importing existing video libraries, use the YouTube category to pre-assign a starting Perspectize category

---

## Mapping Table

| YT ID | YouTube Category | Primary Google NL Path | Secondary Google NL Paths | Notes |
|-------|-----------------|------------------------|---------------------------|-------|
| 1 | Film & Animation | `/Arts & Entertainment/Movies` | `/Arts & Entertainment/Comics & Animation` | Animation maps to Comics & Animation; live-action to Movies |
| 2 | Autos & Vehicles | `/Autos & Vehicles` | — | Direct 1:1 match |
| 10 | Music | `/Arts & Entertainment/Music & Audio` | — | Direct match; subcategories (Classical, Jazz) require Google NL classification |
| 15 | Pets & Animals | `/Pets & Animals` | — | Direct 1:1 match |
| 17 | Sports | `/Sports` | — | Direct match; use Google NL to narrow to Team Sports, Individual Sports, etc. |
| 18 | Short Movies | `/Arts & Entertainment/Movies` | — | Short-form film content |
| 19 | Travel & Events | `/Travel & Transportation` | `/Arts & Entertainment/Events & Listings` | Split: travel vlogs vs event coverage |
| 20 | Gaming | `/Games/Computer & Video Games` | `/Games` | Gaming is almost always video games on YouTube |
| 21 | Videoblogging | `/People & Society` | — | Personal vlogs map to People & Society |
| 22 | People & Blogs | `/People & Society` | `/People & Society/Family & Relationships` | Broad category; Google NL will narrow |
| 23 | Comedy | `/Arts & Entertainment` | `/Arts & Entertainment/TV & Video` | No direct Google NL comedy category; maps to Arts & Entertainment |
| 24 | Entertainment | `/Arts & Entertainment` | — | Broad match; use Google NL for specificity |
| 25 | News & Politics | `/News` | `/Law & Government`, `/News/Politics` | Split: news → /News; politics → /News/Politics or /Law & Government |
| 26 | Howto & Style | `/Hobbies & Leisure` | `/Beauty & Fitness`, `/Home & Garden` | Style → Beauty & Fitness; howto → Hobbies & Leisure or Home & Garden |
| 27 | Education | `/Jobs & Education` | `/Science` | Academic education → Jobs & Education; science education → /Science |
| 28 | Science & Technology | `/Science` | `/Computers & Electronics`, `/Science/Computer Science` | Tech reviews → Computers & Electronics; science explainers → /Science |
| 29 | Nonprofits & Activism | `/People & Society/Social Issues & Advocacy` | `/People & Society/Social Issues & Advocacy/Charity & Philanthropy` | Direct match at L2/L3 |
| 30 | Movies | `/Arts & Entertainment/Movies` | — | Official movie uploads |
| 31 | Anime/Animation | `/Arts & Entertainment/Comics & Animation/Anime & Manga` | `/Arts & Entertainment/Comics & Animation/Cartoons` | Anime → Anime & Manga; western animation → Cartoons |
| 32 | Action/Adventure | `/Arts & Entertainment/Movies/Action & Adventure Films` | — | Film genre |
| 33 | Classics | `/Arts & Entertainment/Movies` | — | Classic films |
| 34 | Comedy (Film) | `/Arts & Entertainment/Movies` | — | Comedy genre films |
| 35 | Documentary | `/Arts & Entertainment/Movies` | `/News` | Documentaries may also match /News |
| 36 | Drama | `/Arts & Entertainment/Movies` | `/Arts & Entertainment/TV & Video/TV Shows & Programs` | Drama films and TV dramas |
| 37 | Family | `/Arts & Entertainment/Movies` | `/People & Society/Family & Relationships` | Family films |
| 38 | Foreign | `/Arts & Entertainment/Movies` | — | International cinema |
| 39 | Horror | `/Arts & Entertainment/Movies` | — | Horror films |
| 40 | Sci-Fi/Fantasy | `/Arts & Entertainment/Movies` | `/Science` | Sci-fi may touch /Science |
| 41 | Thriller | `/Arts & Entertainment/Movies` | — | Thriller films |
| 42 | Shorts | `/Arts & Entertainment/TV & Video` | — | Short-form content |
| 43 | Shows | `/Arts & Entertainment/TV & Video/TV Shows & Programs` | — | TV show uploads |
| 44 | Trailers | `/Arts & Entertainment/Entertainment Industry/Film & TV Industry` | — | Movie/show trailers |

---

## Classification Strategy

### Primary: Google NL Auto-Classification

For every YouTube video added to Perspectize, classify by concatenating metadata:

```
{title}\n{description}\n{tags joined by commas}
```

**Requirements:**
- Minimum 20 tokens (words) — pad with channel name if needed
- Use V2 model: `classificationModelOptions.v2Model.contentCategoriesVersion = "V2"`
- Send to: `POST https://language.googleapis.com/v1/documents:classifyText`

**Category selection from response:**
1. Filter to categories present in the curated Perspectize taxonomy (exact path match or ancestor match)
2. **Primary category** — highest confidence score among curated categories, must be >= 0.5
3. **Secondary categories** — all curated categories with confidence >= 0.3 (stored as labels for filtering)
4. **No match** — if no curated category exceeds 0.5 confidence, fall back to YouTube category mapping

### Fallback: YouTube Category ID Mapping

When Google NL classification fails (confidence < 0.5 for all categories, or API error):

1. Look up the video's `categoryId` from YouTube Data API response
2. Map to primary Google NL path using the table above
3. Assign the mapped category with `confidence = null` (marking it as a fallback, not ML-assigned)
4. Flag the video for manual re-classification or future reclassification

### Confidence Thresholds

| Confidence Range | Interpretation | Action |
|-----------------|----------------|--------|
| >= 0.7 | Strong match | Assign as primary, high confidence |
| 0.5 – 0.69 | Good match | Assign as primary |
| 0.3 – 0.49 | Weak match | Assign as secondary/tag only |
| < 0.3 | Noise | Discard |
| null | YouTube fallback | Flag for manual review |

### Implementation Notes

- **Production auto-classification:** Use Claude Haiku (not Google NL API) — cheaper at scale (~$0.031/month for 100 videos vs. Google NL's per-unit pricing). Use the Google NL taxonomy paths to guide Claude's output.
- **Google NL API:** Use only for taxonomy exploration and initial seed data generation. Not for production classification.
- **ltree matching:** When Google NL returns `/Sports/Team Sports/American Football`, check whether `sports.team_sports.american_football` exists in the categories table. If it does, use it. If not, walk up the path to find the deepest matching ancestor (`sports.team_sports`, then `sports`).

---

## Path Resolution Algorithm

```
function resolveCategory(googleNLPath):
  ltreePath = transformToLtree(googleNLPath)

  // Try exact match first
  category = db.findByPath(ltreePath)
  if category: return category

  // Walk up ancestors
  parts = ltreePath.split('.')
  while parts.length > 1:
    parts.pop()
    ancestorPath = parts.join('.')
    category = db.findByPath(ancestorPath)
    if category: return category

  // No match in curated taxonomy
  return null
```

---

## Gaps and Limitations

### YouTube Categories Without Good Google NL Matches

| YouTube Category | Issue | Recommendation |
|-----------------|-------|----------------|
| Comedy (ID 23) | No Google NL comedy category | Map to `/Arts & Entertainment`; consider adding `arts_and_entertainment.comedy` as user category |
| Videoblogging (ID 21) | Vlogs span many topics | Map to `/People & Society`; rely on Google NL for specificity |
| Shorts (ID 42) | Format, not content | Map based on content; ignore format categorization |

### Google NL Categories Not Mapped from YouTube

These curated categories require Google NL classification to reach (no YouTube category maps directly):
- `/Health` and subcategories (Health Conditions, Mental Health, Nutrition)
- `/Finance`
- `/Books & Literature`
- `/Home & Garden`
- `/Business & Industrial` subtopics
- `/Science/Astronomy`, `/Science/Biology`, `/Science/Physics`

For videos in these categories, rely exclusively on Google NL classification or user manual assignment.

---

## YouTube API Reference

Retrieve category list for a region:
```
GET https://www.googleapis.com/youtube/v3/videoCategories?part=snippet&regionCode=US&key={API_KEY}
```

Note: Not all 44 category IDs are assignable by uploaders. Categories 30–44 are primarily used by YouTube's official content (movie rentals, etc.) and not typical user uploads. The most common upload categories are IDs 1, 2, 10, 15, 17, 20, 22, 24, 25, 26, 27, 28, 29.
