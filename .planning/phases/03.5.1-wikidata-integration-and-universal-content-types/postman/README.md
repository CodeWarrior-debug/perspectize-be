# Wikidata Integration — Postman collection

`wikidata-integration.postman_collection.json` is a Postman Collection v2.1
that exercises the Phase 03.5.1 Wikidata integration end to end:

1. **Raw Wikidata Action API** — the exact `wbsearchentities` request the
   backend `wikidata.Client` adapter makes
   (`backend/internal/adapters/wikidata/client.go`).
2. **Perspectize GraphQL API** — the `wikidataSearch` query and
   `setPrimaryCategory` mutation that sit on top, run against a local server.

## Import

Postman → **Import** → drop in
`wikidata-integration.postman_collection.json` (or paste its raw URL).
No environment file is needed — everything is driven by
**collection variables** (Collection → **Variables** tab).

## Variables to set

| Variable | Default | When to change |
|----------|---------|----------------|
| `wikidata_action_api` | `https://www.wikidata.org/w/api.php` | Never (mirror of the adapter's `defaultBaseURL`). |
| `user_agent` | `Perspectize/1.0 (https://github.com/CodeWarrior-debug/perspectize)` | Never — Wikidata blocks generic/blank agents. |
| `search_query` | `Kansas City Chiefs` | Any term you want to search for. |
| `language` | `en` | Any Wikidata language code. |
| `limit` | `10` | Result cap. Keep it a bare integer (GraphQL `Int`). |
| `graphql_url` | `http://localhost:8080/graphql` | Point at your running backend. |
| `content_id` | `1` | **Required for Folder 2b** — a real `Content` row id. |
| `chiefs_qid` | `Q223522` | Known QID for ad-hoc `wbgetentities` / mutation tests. |
| `category_qid` / `category_label` | `Q223522` / `Kansas City Chiefs` | Passed to `setPrimaryCategory`. Auto-populated by request **2a**'s test script if you run it first. |
| `last_qid` | _(empty)_ | Scratch — top QID from the most recent search. |

## Folders

- **0. Quickstart — Validate Access**
  - `0a` `wbsearchentities` for "Kansas City Chiefs" (the adapter's call).
  - `0b` `wbgetentities` detail for `Q223522`
    (`props=labels|descriptions|claims`).
- **1. Entity Search (as the backend calls it)** — `wbsearchentities`
  parameterised with `{{search_query}}` / `{{language}}` / `{{limit}}`:
  happy path, gibberish query (empty results, still HTTP 200), explicit
  `limit=3`.
- **2. Perspectize GraphQL API** — `POST {{graphql_url}}`:
  - `2a` `wikidataSearch(query, language, limit)` →
    `{ qid label description entityType }`.
  - `2b` `setPrimaryCategory(input: SetPrimaryCategoryInput!)` → `Content`
    with its full `primaryCategory { id wikidataQid label description
    entityType createdAt updatedAt }`.

Every request carries `pm.test(...)` assertions. Wikidata requests assert
HTTP 200, a `search` array, and `id`/`label` on each hit (what the adapter
maps to `WikidataSearchResult`). GraphQL requests assert `data` is present
and `errors` is absent.

## Running Folder 2

```bash
cd backend
make run        # server on :8080, GraphQL at /graphql
```

The GraphQL route is `/graphql` (from `backend/cmd/server/main.go`:
`r.Handle("/graphql", srv)`; port `8080` from
`backend/config/config.example.json`). There is no `/query` route.

Run **2a** before **2b** so the search result's QID/label flow into the
mutation via collection variables. `2b` also needs `content_id` pointed at a
real row.

## Suggested order

Run the collection top to bottom (Postman **Collection Runner** works): the
test scripts chain state forward (`last_qid`, `category_qid`,
`category_label`).
