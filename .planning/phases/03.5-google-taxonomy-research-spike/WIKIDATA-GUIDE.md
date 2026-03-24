# Wikidata Taxonomy Explorer — Quickstart Guide

**Purpose:** Explore Wikidata's knowledge graph as an alternative (or complement) to Google NL taxonomy for Perspectize content categorization. Wikidata provides the granularity Google NL lacks — specific teams, seasons, denominations, people, and temporal context.

**Time required:** ~30 minutes for a complete walkthrough

---

## Why Wikidata Over Google NL?

| Need | Google NL | Wikidata |
|------|-----------|----------|
| "Sports" | `/Sports/Team Sports/American Football` | Same, but also... |
| "Chiefs" | Cannot express | `Q223522` — Kansas City Chiefs |
| "Chiefs 2024 season" | Cannot express | `Q124553782` — 2024 Kansas City Chiefs season |
| "Presbyterian Church (USA)" | Cannot express | `Q1149100` — specific denomination |
| Updated | Static since 2022 | Community-edited in real-time |
| Total entities | 1,091 categories | 100M+ items |
| Auto-classify text | Yes (API classifies text) | No (it's a vocabulary, not a classifier) |

**Key insight:** Google NL tells you "this text is about sports." Wikidata tells you "the Kansas City Chiefs are an American football team in the NFL, founded in 1960, who play at GEHA Field, and their 2024 season is a distinct entity."

---

## Prerequisites

1. **Postman** installed (desktop or web)
2. **Import the collection:** File → Import → select `postman/wikidata-taxonomy-explorer.json`
3. **No API key needed** — Wikidata is fully open

That's it. No Google Cloud project, no credentials.

---

## Three APIs in One Collection

| API | What It Does | When to Use |
|-----|-------------|-------------|
| **Entity Search** (Action API) | Search by name → get Q-numbers | "What's the entity ID for the Chiefs?" |
| **REST API** | Get full entity details by Q-number | "Show me everything about Q223522" |
| **SPARQL Endpoint** | Structured queries across the graph | "Give me all NFL teams with their stadiums" |

---

## Phase A: Validate Access (5 min)

### Step 1: Run `0a - Entity Search`

Click Send. You should get back:

```json
{
  "search": [
    {
      "id": "Q223522",
      "label": "Kansas City Chiefs",
      "description": "National Football League franchise in Kansas City, Missouri"
    }
  ]
}
```

If this works, you have access. No auth needed.

### Step 2: Run `0b - Get Entity Details`

This fetches the full Kansas City Chiefs entity. Look at the response shape:

```json
{
  "type": "item",
  "id": "Q223522",
  "labels": { "en": "Kansas City Chiefs", "es": "Kansas City Chiefs", ... },
  "descriptions": { "en": "National Football League franchise in Kansas City, Missouri" },
  "aliases": { "en": ["Dallas Texans", "KC Chiefs", "Chiefs"] },
  "statements": {
    "P31": [ ... ],   // instance of: American football team
    "P118": [ ... ],  // league: NFL
    "P115": [ ... ],  // home venue: GEHA Field
    "P571": [ ... ],  // inception: 1960
    "P286": [ ... ],  // head coach
    ...
  },
  "sitelinks": { "enwiki": { "title": "Kansas City Chiefs" }, ... }
}
```

**Data shape:** Every Wikidata entity has the same structure — labels, descriptions, aliases, statements (property→value pairs), sitelinks (Wikipedia links).

### Step 3: Run `0c - SPARQL Hello World`

Confirms the SPARQL endpoint works. Returns 3 known entities with labels and descriptions.

**SPARQL response shape:**
```json
{
  "head": { "vars": ["item", "itemLabel", "itemDescription"] },
  "results": {
    "bindings": [
      {
        "item": { "type": "uri", "value": "http://www.wikidata.org/entity/Q223522" },
        "itemLabel": { "type": "literal", "value": "Kansas City Chiefs" },
        "itemDescription": { "type": "literal", "value": "National Football League franchise..." }
      }
    ]
  }
}
```

Every SPARQL result has the same shape: `head.vars` lists column names, `results.bindings` is an array of row objects, each value has `type` and `value`.

---

## Phase B: Sports Exploration (10 min)

### Step 4: Run `1a - All NFL Teams`

Returns every NFL team with stadium. Note:
- Historical teams appear too (Chicago Tigers, Columbus Panhandles)
- Each team has a Q-number you can use in further queries

### Step 5: Run `1b - Kansas City Chiefs Seasons`

**This is the killer feature.** Returns:
```
2026 Kansas City Chiefs season
2025 Kansas City Chiefs season
2024 Kansas City Chiefs season
2023 Kansas City Chiefs season
...
```

Each season is a separate Wikidata entity with its own properties (win/loss record, coach, etc.). Google NL can never give you this.

### Step 6: Run `1c - All Major US Sports Leagues`

Reference query — gives you Q-numbers for NFL, NBA, MLB, MLS, NHL, Premier League, La Liga. Use these to swap into query 1a for any league's teams.

### Step 7: Run `1d - Class Hierarchy`

Shows the full "is-a" chain: Kansas City Chiefs → American football team → sports team → organization → entity.

---

## Phase C: Religion Exploration (10 min)

### Step 8: Run `2a - Search Denominations`

Search "Presbyterian" — you'll see:
- `Q178169` — Presbyterianism (the tradition)
- `Q1149100` — Presbyterian Church (the US denomination)
- `Q5024548` — Calvinistic Methodists (now Presbyterian Church of Wales)

### Step 9: Run `2b - Major Denomination Families`

Reference lookup of 10 major traditions with Q-numbers.

### Step 10: Run `2c - All Presbyterian Denominations`

Finds organizations whose religion (P140) is Presbyterianism. Shows specific bodies by country.

**Try changing the query:** Replace `Q178169` with `Q93191` (Baptist) or `Q33203` (Methodist) to explore those families.

### Step 11: Run `2d - World Religions`

Top-level reference for non-Christian religions — Buddhism, Islam, Hinduism, etc. with Q-numbers for further drilling.

---

## Phase D: Entertainment & Tech (5 min)

### Step 12-14: Run folder 3 (Entertainment) and folder 4 (Science/Tech)

Quick exploration of music genres, video game genres, film genres, academic disciplines, and programming languages.

---

## Phase E: Master the Traversal Patterns (10 min)

Folder 5 contains **reusable query patterns** — these are the most important:

### Step 15: `5a - Walk UP` (entity → root classes)

**Pattern:** `wd:ENTITY wdt:P31/wdt:P279* ?class`

Replace the entity ID with anything. Builds the taxonomy breadcrumb.

### Step 16: `5b - Walk DOWN` (class → all instances)

**Pattern:** `?item wdt:P31/wdt:P279* wd:CLASS`

Find all items that are instances of a class (including subclasses).

### Step 17: `5c - Direct Subclasses` (one level)

**Pattern:** `?subclass wdt:P279 wd:CLASS`

Get immediate children only — good for building a tree level by level.

### Step 18: `5d - Entity Relationships` (all properties)

**Pattern:** `wd:ENTITY ?p ?value . ?prop wikibase:directClaim ?p`

See everything Wikidata knows about any entity. Great for discovering new properties to query on.

---

## Wikidata Cheat Sheet

### Key Properties (P-numbers)

| Property | Meaning | Example |
|----------|---------|---------|
| **P31** | instance of | Chiefs P31 American football team |
| **P279** | subclass of | American football P279 team sport |
| **P118** | league | Chiefs P118 NFL |
| **P115** | home venue | Chiefs P115 GEHA Field |
| **P140** | religion | Org P140 Presbyterianism |
| **P5138** | season of team | Season P5138 Chiefs |
| **P641** | sport | NFL P641 American football |
| **P17** | country | Entity P17 United States |
| **P571** | inception (founding date) | Chiefs P571 1960 |
| **P361** | part of | AFC West P361 AFC |

### Key Entity Types (Q-numbers)

| Entity | Q-number | Use for |
|--------|----------|---------|
| American football team | Q17156793 | NFL teams |
| Basketball team | Q13393265 | NBA teams |
| Baseball team | Q476028 | MLB teams |
| Sports season | Q27020041 | Team seasons |
| Music genre | Q188451 | Music categories |
| Film genre | Q201658 | Movie categories |
| Video game genre | Q659563 | Gaming categories |
| Academic discipline | Q11862829 | Science/education |
| Programming language | Q9143 | Tech categories |
| Christian denomination | Q1068640 | Denominations |
| YouTube channel | Q36834616 | YouTube entities |

### SPARQL Patterns

```sparql
-- Find by name (use Action API instead — faster)
-- wbsearchentities?search=Chiefs

-- All instances of a class
SELECT ?item ?itemLabel WHERE {
  ?item wdt:P31 wd:Q17156793 .  -- instance of: American football team
  SERVICE wikibase:label { bd:serviceParam wikibase:language "[AUTO_LANGUAGE],en" . }
}

-- Walk UP the hierarchy
SELECT ?class ?classLabel WHERE {
  wd:Q223522 wdt:P31/wdt:P279* ?class .
  SERVICE wikibase:label { bd:serviceParam wikibase:language "[AUTO_LANGUAGE],en" . }
}

-- Walk DOWN (all instances including subclasses)
SELECT ?item ?itemLabel WHERE {
  ?item wdt:P31/wdt:P279* wd:Q349 .  -- anything that's a sport
  SERVICE wikibase:label { bd:serviceParam wikibase:language "[AUTO_LANGUAGE],en" . }
}
LIMIT 100

-- Direct children only
SELECT ?child ?childLabel WHERE {
  ?child wdt:P279 wd:Q349 .  -- direct subclass of: sport
  SERVICE wikibase:label { bd:serviceParam wikibase:language "[AUTO_LANGUAGE],en" . }
}
```

---

## Rate Limits & Rules

| Limit | Value |
|-------|-------|
| Query timeout | 60 seconds |
| Compute budget | 60s of query time per minute |
| Parallel queries | 5 per IP |
| Error budget | 30 errors per minute |
| Authentication | None required |
| Cost | Free |

**Mandatory:** Set a `User-Agent` header (already configured in the collection). Wikidata blocks requests with no/default user agent.

---

## Wikidata vs Google NL: Decision Matrix for Perspectize

| If you need... | Use |
|---------------|-----|
| Broad content classification of text | Google NL (or Claude Haiku) |
| Specific teams, people, denominations | Wikidata |
| Temporal context (seasons, years) | Wikidata |
| Auto-classify a YouTube video | Claude Haiku (neither API does this well alone) |
| Stable category IDs that won't change | Wikidata Q-numbers (persistent identifiers) |
| Hierarchical browsing tree | Both (Google NL as ltree, Wikidata via P279 chains) |

**Recommended hybrid:** Use Google NL / Claude Haiku for broad top-level classification, then enrich with Wikidata entities for specificity. A video classified as "Sports > American Football" can be further tagged with `Q223522` (Chiefs) and linked to `Q124553782` (2024 season).

---

## Sources

- [Wikidata SPARQL Query Service](https://query.wikidata.org/)
- [Wikidata SPARQL Tutorial](https://www.wikidata.org/wiki/Wikidata:SPARQL_tutorial)
- [Wikidata Data Access](https://www.wikidata.org/wiki/Wikidata:Data_access)
- [Wikidata REST API](https://www.wikidata.org/wiki/Wikidata:REST_API)
- [SPARQL Query Examples](https://www.wikidata.org/wiki/Wikidata:SPARQL_query_service/queries/examples)
- [Query Limits](https://www.wikidata.org/wiki/Wikidata:SPARQL_query_service/query_limits)
