# Google NL Taxonomy Exploration Guide

**Purpose:** Systematically discover Google's content taxonomy and curate the category list for Perspectize Phase 13 (Content Categories) implementation.

**Time required:** ~50 minutes for a complete exploration

---

## Prerequisites

Before starting, confirm you have:

1. **Google Cloud Project** with the Natural Language API enabled
   - Go to: Google Cloud Console > APIs & Services > Library > search "Natural Language API" > Enable
2. **API Key** from Google Cloud Console
   - Go to: APIs & Services > Credentials > Create Credentials > API Key
   - Restrict the key to "Cloud Natural Language API" only (for security)
3. **Postman** installed (desktop app or web)
4. **Postman collection imported**
   - In Postman: click Import > choose file > select `postman/google-nl-taxonomy-explorer.json`
5. **API_KEY variable set**
   - In Postman: click the collection name "Google NL Taxonomy Explorer" > Variables tab
   - Set the `API_KEY` current value to your API key from step 2

---

## Phase A: Validate API Access (5 min)

Confirm your API key works and understand V1 vs V2 model behavior before deep exploration.

### Step 1: Run Curl 1 (V1 Basic)

Open the folder **"1. V1 vs V2 Comparison"** and run **"Curl 1 - V1 Basic (filtered, single category)"**.

- Check: You receive a `200 OK` response (not a 401 or 403)
- Check: The response contains a `categories` array with at least one item
- Check: The category name looks like `/Sports/Team Sports/American Football`

If you get a 403 error: verify the API key is correct and Natural Language API is enabled in your Cloud project.

### Step 2: Run Curl 2 (V2 Model)

Run **"Curl 2 - V2 Model (unfiltered, full hierarchy)"** in the same folder.

- Check: Multiple categories are returned (not just one)
- Check: Parent categories appear alongside leaf categories (e.g., `/Sports` and `/Sports/Team Sports` alongside `/Sports/Team Sports/American Football`)
- Check: A cross-category like `/News/Sports News` may also appear

### Step 3: Compare V1 vs V2 output

Side-by-side comparison:

| Aspect | V1 (Curl 1) | V2 (Curl 2) |
|--------|-------------|-------------|
| Categories returned | 1 (most specific) | Multiple (full hierarchy) |
| Parent categories | No | Yes |
| Cross-categories | No | Yes |
| Use case | Single best match | Taxonomy discovery |

**Conclusion:** All remaining exploration uses V2. V1 is shown only for comparison.

---

## Phase B: Explore Top-Level Coverage (15 min)

Run each content-type request to discover which of the 20 curated top-level categories appear in practice.

### Step 4: Run Curl 3a (Basketball)

Open folder **"2. Sports"** and run **"Curl 3a - Basketball (V2)"**.

Record:
- What top-level category appears first? (`/Sports`)
- What is the deepest path returned? (e.g., `/Sports/Team Sports/Basketball`)
- Maximum depth (count the `/` separators + 1)
- Any unexpected cross-categories?

### Step 5: Run Curl 3b (Soccer)

Run **"Curl 3b - Soccer (V2)"** from the same folder.

Record:
- Does `/Sports/Team Sports/Soccer` appear, or is it labeled differently?
- Compare confidence scores with Basketball — are they similar for the same depth?
- Note: Watch for "Football" vs "Soccer" naming in the taxonomy

### Step 6: Run Curl 4 (Entertainment)

Open folder **"3. Entertainment"** and run **"Curl 4 - Music & Entertainment (V2)"**.

Record:
- Does `/Arts & Entertainment/Music & Audio` appear?
- Does `/Arts & Entertainment/Events & Listings/Concerts & Music Festivals` appear?
- Any cross-classification into `/People & Society`?

### Step 7: Run Curl 5 (Technology)

Open folder **"4. Technology"** and run **"Curl 5 - Technology & AI (V2)"**.

Record:
- Does `/Science/Computer Science/Machine Learning & Artificial Intelligence` appear?
- Does `/Computers & Electronics` appear?
- Which branch scores higher confidence — Science or Computers?

### Step 8: Run Curl 6 (News & Politics)

Open folder **"5. News & Politics"** and run **"Curl 6 - News & Politics (V2)"**.

Record:
- Does `/News/Politics` appear?
- Does `/Law & Government` appear?
- How deeply does the `/News` branch classify political content?

### Step 9: Run Curl 7 (Gaming)

Open folder **"6. Gaming"** and run **"Curl 7 - Gaming (V2)"**.

Record:
- Does `/Games/Computer & Video Games` appear?
- Does `/Games/Computer & Video Games/Shooter Games` appear for Call of Duty?
- Note: Top-level category is `/Games`, not `/Gaming`

### Step 10: Run Curl 8 (Food & Cooking)

Open folder **"7. Food & Cooking"** and run **"Curl 8 - Food & Cooking (V2)"**.

Record:
- Does `/Food & Drink/Cooking & Recipes` appear?
- Are there cuisine-specific subcategories?
- Does `/Food & Drink/Beverages` appear for the wine mention?

### Step 11: Run Curl 9 (Health & Fitness)

Open folder **"8. Health & Fitness"** and run **"Curl 9 - Health & Fitness (V2)"**.

Record:
- Does `/Beauty & Fitness` or `/Health` appear? Which scores higher?
- Does `/Sports/Sports Coaching & Training` appear?
- Note: Fitness content sits in a surprising place — `/Beauty & Fitness` is one top-level branch, `/Health` is another

### After Phase B: Tally your top-level categories

List every unique first path segment you observed across all requests. Compare with the 20 curated categories from the research:

**Expected to appear:** Arts & Entertainment, Beauty & Fitness, Computers & Electronics, Food & Drink, Games, Health, Law & Government, News, People & Society, Science, Sports

**May not appear** from these requests: Autos & Vehicles, Books & Literature, Business & Industrial, Finance, Hobbies & Leisure, Home & Garden, Jobs & Education, Pets & Animals, Travel & Transportation

---

## Phase C: Deep-Dive Subcategories (20 min)

Pick 3-5 top-level categories most relevant to your YouTube content and explore their subcategories.

### Step 12: Choose your priority categories

Pick from this list based on what types of YouTube content you rate most:

- `/Sports` — team sports, individual sports, specific leagues
- `/Arts & Entertainment` — music, movies, TV shows, performers
- `/Games` — video games, specific genres
- `/News` — politics, business news, sports news
- `/Science` — AI, physics, biology
- `/Food & Drink` — cooking, restaurants, beverages

### Step 13: Write custom requests to explore subcategories

In Postman, duplicate **"Curl 2 - V2 Model"** (right-click > Duplicate) and modify the `content` field to target specific subcategories.

**Technique:** Write longer, more specific content about a topic to draw out deeper subcategory classifications.

**Example: Exploring /Sports/Team Sports deeper**

```json
{
  "document": {
    "type": "PLAIN_TEXT",
    "content": "The Kansas City Chiefs and Philadelphia Eagles faced off in Super Bowl LVII. Patrick Mahomes led the AFC's top seed while Jalen Hurts powered the NFC champions. The game featured an incredible fourth quarter comeback as the Chiefs defense held the Eagles offense scoreless in the final minutes."
  },
  "classificationModelOptions": {
    "v2Model": {
      "contentCategoriesVersion": "V2"
    }
  }
}
```

Does `/Sports/Team Sports/American Football` appear? Does any NFL-specific subcategory exist beyond that?

**Example: Exploring /Games deeper**

```json
{
  "document": {
    "type": "PLAIN_TEXT",
    "content": "The World of Warcraft expansion features new dungeons, raids, and a revamped crafting system. Players level up their characters through the new story campaign and unlock gear through weekly mythic plus dungeons. The guild community aspect remains central to the endgame experience."
  },
  "classificationModelOptions": {
    "v2Model": {
      "contentCategoriesVersion": "V2"
    }
  }
}
```

Does `/Games/Computer & Video Games` split into MMORPG or RPG subcategories?

### Step 14: Record the deepest paths you find

For each priority category, note the deepest path returned, its confidence score, and whether the taxonomy stops being useful at that depth.

**Useful (keep in seed data):** Specific enough to be meaningful, general enough to apply to multiple videos

**Too generic (skip):** `/Sports` alone — not useful for navigation

**Too niche (skip or defer):** Anything that applies to only 1-2 videos — users can create these themselves

---

## Phase D: YouTube Metadata Test (10 min)

Validate that the classification works well with actual YouTube video metadata format.

### Step 15: Run Curl 10 (YouTube Metadata Format)

Open folder **"9. YouTube Video Metadata"** and run **"Curl 10 - YouTube Metadata Format (V2)"**.

- Check: The response matches `/Sports/Team Sports/American Football`
- Check: Confidence scores on the primary category — does 0.5 feel right as the threshold?

### Step 16: Test with your own YouTube videos

Duplicate **"Curl 10 - YouTube Metadata Format (V2)"** and replace the content with your own video metadata.

**Template:**
```
TITLE: [your video title here]

DESCRIPTION: [your video description here]

TAGS: [comma-separated tags]
```

Requirements:
- The combined text must be at least 20 words
- Include title, description, and tags for best results

**For each video, evaluate:**

1. Does the primary category (highest confidence) match what you would manually assign?
2. Are secondary categories (0.3-0.5 confidence) useful as labels or tags?
3. Are any categories returned that are clearly wrong?
4. What confidence threshold feels right?
   - `>= 0.7` — very confident, use as primary category
   - `0.5-0.7` — confident, still use as primary
   - `0.3-0.5` — related, use as secondary label
   - `< 0.3` — tangential, ignore

### Step 17: Document your confidence threshold decision

After testing 3+ videos, decide:

- **Primary category threshold:** The minimum confidence to assign a primary category (suggested: 0.5)
- **Secondary category threshold:** The minimum confidence to assign secondary labels (suggested: 0.3)

This decision feeds directly into Phase 13 implementation.

---

## Phase E: Record Findings (varies)

Consolidate everything you discovered into the findings template below.

### Step 18: Compile all unique category paths

From all requests run in Phases B and C, copy every unique category path from the responses. Group them by top-level category.

### Step 19: Assess relevance

For each path, decide:

- **KEEP:** Include in seed data (relevant to Perspectize YouTube content)
- **SKIP:** Don't seed (irrelevant, too niche, or users won't find it)
- **CUSTOM NEEDED:** A gap where user-created categories should go

### Step 20: Decide seed depth

Based on exploration, decide how many levels deep to seed:

- **2 levels:** `/Sports/Team Sports` — broad but navigable
- **3 levels:** `/Sports/Team Sports/American Football` — specific and useful
- **4-5 levels:** Usually too niche for seed data — defer to user-created categories

**Recommendation from research:** Seed 2-3 levels for most categories. Sport-specific and genre-specific content often warrants 3 levels. Use user-created categories for anything deeper.

---

## Findings Template

Copy this section and fill it in as you explore. This is the primary deliverable that feeds into Plan 02 (curated list and seed data).

```markdown
## My Findings

### Categories Discovered (copy paths from responses)

**Sports:**
- /Sports (confidence: X.XX)
- /Sports/Team Sports (confidence: X.XX)
- /Sports/Team Sports/American Football (confidence: X.XX)
- /Sports/Team Sports/Basketball (confidence: X.XX)
- /Sports/Team Sports/Soccer (confidence: X.XX)
- [add more as discovered]

**Arts & Entertainment:**
- /Arts & Entertainment (confidence: X.XX)
- /Arts & Entertainment/Music & Audio (confidence: X.XX)
- [add more as discovered]

**Technology:**
- /Computers & Electronics (confidence: X.XX)
- /Science/Computer Science (confidence: X.XX)
- /Science/Computer Science/Machine Learning & Artificial Intelligence (confidence: X.XX)
- [add more as discovered]

**News:**
- /News (confidence: X.XX)
- /News/Politics (confidence: X.XX)
- /News/Sports News (confidence: X.XX)
- [add more as discovered]

**Gaming:**
- /Games (confidence: X.XX)
- /Games/Computer & Video Games (confidence: X.XX)
- /Games/Computer & Video Games/Shooter Games (confidence: X.XX)
- [add more as discovered]

**Food & Drink:**
- /Food & Drink (confidence: X.XX)
- /Food & Drink/Cooking & Recipes (confidence: X.XX)
- [add more as discovered]

**Health & Fitness:**
- /Health (confidence: X.XX)
- /Beauty & Fitness (confidence: X.XX)
- [add more as discovered]

**Other top-level categories observed:**
- [any other top-level categories seen in responses]

### Relevance Assessment

**KEEP (include in seed data):**
- [paths to include, with brief rationale]

**SKIP (exclude from seed data):**
- [paths too niche, irrelevant, or better left to users]

**CUSTOM NEEDED (gaps where user categories should go):**
- [describe what's missing from the Google taxonomy]

### Depth Decision

- Seed depth: [2 or 3 levels]
- Reasoning: [why — e.g., "3 levels for sports and gaming where specificity matters; 2 levels for news and entertainment where subcategories are less distinct"]

### Confidence Threshold

- Primary category: [threshold, e.g., 0.5]
- Secondary/label: [threshold, e.g., 0.3]
- Reasoning: [what you observed that led to this decision]

### Surprising Discoveries

- [Anything unexpected — cross-categories, missing categories, surprising depth]

### Category Count

- Total unique paths discovered: [N]
- Paths to include in seed data: [N]
- Estimated seed data rows: [N] (each path = one row in categories table)
```

---

## Quick Reference: Request Names

When referencing requests in your notes, use these names:

| Request Name | Folder | Tests |
|---|---|---|
| Curl 1 - V1 Basic (filtered, single category) | 1. V1 vs V2 Comparison | V1 model behavior |
| Curl 2 - V2 Model (unfiltered, full hierarchy) | 1. V1 vs V2 Comparison | V2 hierarchy output |
| Curl 3a - Basketball (V2) | 2. Sports | /Sports/Team Sports/Basketball |
| Curl 3b - Soccer (V2) | 2. Sports | /Sports/Team Sports/Soccer |
| Curl 4 - Music & Entertainment (V2) | 3. Entertainment | /Arts & Entertainment/Music & Audio |
| Curl 5 - Technology & AI (V2) | 4. Technology | /Science/Computer Science/Machine Learning... |
| Curl 6 - News & Politics (V2) | 5. News & Politics | /News/Politics |
| Curl 7 - Gaming (V2) | 6. Gaming | /Games/Computer & Video Games |
| Curl 8 - Food & Cooking (V2) | 7. Food & Cooking | /Food & Drink/Cooking & Recipes |
| Curl 9 - Health & Fitness (V2) | 8. Health & Fitness | /Health, /Beauty & Fitness |
| Curl 10 - YouTube Metadata Format (V2) | 9. YouTube Video Metadata | Realistic video classification |

---

## What Comes Next

After completing this guide:

1. **Fill in the Findings Template** with paths discovered and your relevance assessment
2. **Share your findings** — the completed template feeds directly into Plan 02 (curated category list and seed data SQL)
3. **Plan 02 will generate:**
   - The final curated list of ~50-100 categories
   - SQL seed data with ltree paths for the categories table
   - The Go transformation function for converting Google NL paths to ltree
4. **Phase 13 implementation** uses the seed data to populate the PostgreSQL categories table and the transformation function for auto-classification

**Free tier note:** The Google Cloud NL API free tier allows 30,000 units/month. Each request in this guide uses approximately 1 unit. You can run hundreds of custom requests without incurring charges.
