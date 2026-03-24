# Perspectize Content Categories -- Curated Taxonomy

**Version:** 1.0 (initial draft — pending user review)
**Source:** Google Cloud Natural Language API V2 taxonomy (1091 categories)
**Curated:** 2026-02-20
**Status:** Awaiting user verification against Postman exploration findings

---

## Overview

This document defines the curated subset of Google NL categories selected for Perspectize. The full taxonomy has 1091 categories across 27 top-level buckets. We curate to 20 top-level categories and approximately 2–3 levels deep, targeting the content types users actually rate on YouTube.

| Property | Value |
|----------|-------|
| Top-level categories | 20 (of 27 total in Google NL) |
| Excluded categories | 7 (Adult, Internet & Telecom, Online Communities, Real Estate, Reference, Sensitive Subjects, Shopping) |
| Total curated categories | ~70 |
| Max depth seeded | 3 levels (top + 2 sub-levels) |
| Source | google_nl |

---

## Depth Decision

**Recommendation: Seed 3 levels deep (top-level + 2 sub-levels).**

Rationale:
- Level 1 (e.g., `/Sports`) — too broad for meaningful ratings
- Level 2 (e.g., `/Sports/Team Sports`) — useful grouping, good filter unit
- Level 3 (e.g., `/Sports/Team Sports/American Football`) — specific enough to be meaningful, broad enough to have many videos
- Level 4+ (e.g., NFL teams, specific athletes) — too narrow for seed data; better as user-created extensions

The full Google NL taxonomy goes to 5 levels but levels 4–5 are typically too granular for initial seeding. Users can extend the tree at any depth via the `user` source.

---

## Category Tree

Format: `Display Name` → `ltree path` (source: google_nl)

### 1. Arts & Entertainment

```
Arts & Entertainment                         → arts_and_entertainment
  Comics & Animation                         → arts_and_entertainment.comics_and_animation
    Anime & Manga                            → arts_and_entertainment.comics_and_animation.anime_and_manga
    Cartoons                                 → arts_and_entertainment.comics_and_animation.cartoons
    Comics                                   → arts_and_entertainment.comics_and_animation.comics
  Entertainment Industry                     → arts_and_entertainment.entertainment_industry
    Film & TV Industry                       → arts_and_entertainment.entertainment_industry.film_and_tv_industry
    Recording Industry                       → arts_and_entertainment.entertainment_industry.recording_industry
  Events & Listings                          → arts_and_entertainment.events_and_listings
    Bars, Clubs & Nightlife                  → arts_and_entertainment.events_and_listings.bars_clubs_and_nightlife
    Concerts & Music Festivals               → arts_and_entertainment.events_and_listings.concerts_and_music_festivals
  Movies                                     → arts_and_entertainment.movies
    Action & Adventure Films                 → arts_and_entertainment.movies.action_and_adventure_films
  Music & Audio                              → arts_and_entertainment.music_and_audio
    Classical Music                          → arts_and_entertainment.music_and_audio.classical_music
    Jazz & Blues                             → arts_and_entertainment.music_and_audio.jazz_and_blues
  Performing Arts                            → arts_and_entertainment.performing_arts
    Dance                                    → arts_and_entertainment.performing_arts.dance
  TV & Video                                 → arts_and_entertainment.tv_and_video
    TV Shows & Programs                      → arts_and_entertainment.tv_and_video.tv_shows_and_programs
  Visual Art & Design                        → arts_and_entertainment.visual_art_and_design
    Architecture                             → arts_and_entertainment.visual_art_and_design.architecture
```

### 2. Autos & Vehicles

```
Autos & Vehicles                             → autos_and_vehicles
```

*(Subcategories to be added after Postman exploration)*

### 3. Beauty & Fitness

```
Beauty & Fitness                             → beauty_and_fitness
```

*(Subcategories to be added after Postman exploration)*

### 4. Books & Literature

```
Books & Literature                           → books_and_literature
```

*(Subcategories to be added after Postman exploration)*

### 5. Business & Industrial

```
Business & Industrial                        → business_and_industrial
  Advertising & Marketing                    → business_and_industrial.advertising_and_marketing
    Public Relations                         → business_and_industrial.advertising_and_marketing.public_relations
  Energy & Utilities                         → business_and_industrial.energy_and_utilities
    Oil & Gas                                → business_and_industrial.energy_and_utilities.oil_and_gas
  Manufacturing                              → business_and_industrial.manufacturing
  Pharmaceuticals & Biotech                  → business_and_industrial.pharmaceuticals_and_biotech
```

### 6. Computers & Electronics

```
Computers & Electronics                      → computers_and_electronics
```

*(Subcategories to be added after Postman exploration)*

### 7. Finance

```
Finance                                      → finance
```

*(Subcategories to be added after Postman exploration)*

### 8. Food & Drink

```
Food & Drink                                 → food_and_drink
  Cooking & Recipes                          → food_and_drink.cooking_and_recipes
  Beverages                                  → food_and_drink.beverages
  Restaurants                                → food_and_drink.restaurants
```

### 9. Games

```
Games                                        → games
  Computer & Video Games                     → games.computer_and_video_games
    Shooter Games                            → games.computer_and_video_games.shooter_games
  Board Games                                → games.board_games
    Chess & Abstract Strategy Games          → games.board_games.chess_and_abstract_strategy_games
  Card Games                                 → games.card_games
    Poker & Casino Games                     → games.card_games.poker_and_casino_games
  Gambling                                   → games.gambling
    Lottery                                  → games.gambling.lottery
```

### 10. Health

```
Health                                       → health
  Health Conditions                          → health.health_conditions
    Cancer                                   → health.health_conditions.cancer
  Mental Health                              → health.mental_health
    Depression                               → health.mental_health.depression
  Medical Facilities & Services              → health.medical_facilities_and_services
    Hospitals & Treatment Centers            → health.medical_facilities_and_services.hospitals_and_treatment_centers
  Nutrition                                  → health.nutrition
    Vitamins & Supplements                   → health.nutrition.vitamins_and_supplements
  Vision Care                                → health.vision_care
    Eyeglasses & Contacts                    → health.vision_care.eyeglasses_and_contacts
```

### 11. Hobbies & Leisure

```
Hobbies & Leisure                            → hobbies_and_leisure
  Crafts                                     → hobbies_and_leisure.crafts
    Ceramics & Pottery                       → hobbies_and_leisure.crafts.ceramics_and_pottery
  Outdoors                                   → hobbies_and_leisure.outdoors
    Fishing                                  → hobbies_and_leisure.outdoors.fishing
  Special Occasions                          → hobbies_and_leisure.special_occasions
    Weddings                                 → hobbies_and_leisure.special_occasions.weddings
  Water Activities                           → hobbies_and_leisure.water_activities
    Boating                                  → hobbies_and_leisure.water_activities.boating
```

### 12. Home & Garden

```
Home & Garden                                → home_and_garden
```

*(Subcategories to be added after Postman exploration)*

### 13. Jobs & Education

```
Jobs & Education                             → jobs_and_education
```

*(Subcategories to be added after Postman exploration)*

### 14. Law & Government

```
Law & Government                             → law_and_government
```

*(Subcategories to be added after Postman exploration)*

### 15. News

```
News                                         → news
  Business News                              → news.business_news
    Company News                             → news.business_news.company_news
  Health News                                → news.health_news
  Politics                                   → news.politics
    Campaigns & Elections                    → news.politics.campaigns_and_elections
  Sports News                                → news.sports_news
  Technology News                            → news.technology_news
  Weather                                    → news.weather
```

### 16. People & Society

```
People & Society                             → people_and_society
  Family & Relationships                     → people_and_society.family_and_relationships
    Marriage                                 → people_and_society.family_and_relationships.marriage
  Kids & Teens                               → people_and_society.kids_and_teens
    Children's Interests                     → people_and_society.kids_and_teens.childrens_interests
  Religion & Belief                          → people_and_society.religion_and_belief
  Social Issues & Advocacy                   → people_and_society.social_issues_and_advocacy
    Charity & Philanthropy                   → people_and_society.social_issues_and_advocacy.charity_and_philanthropy
```

### 17. Pets & Animals

```
Pets & Animals                               → pets_and_animals
```

*(Subcategories to be added after Postman exploration)*

### 18. Science

```
Science                                      → science
  Astronomy                                  → science.astronomy
  Biology                                    → science.biology
    Neuroscience                             → science.biology.neuroscience
  Chemistry                                  → science.chemistry
  Computer Science                           → science.computer_science
    Machine Learning & Artificial Intelligence → science.computer_science.machine_learning_and_artificial_intelligence
  Physics                                    → science.physics
```

### 19. Sports

```
Sports                                       → sports
  Team Sports                                → sports.team_sports
    American Football                        → sports.team_sports.american_football
    Baseball                                 → sports.team_sports.baseball
    Basketball                               → sports.team_sports.basketball
    Hockey                                   → sports.team_sports.hockey
    Soccer                                   → sports.team_sports.soccer
  Individual Sports                          → sports.individual_sports
    Golf                                     → sports.individual_sports.golf
    Cycling                                  → sports.individual_sports.cycling
  Motor Sports                               → sports.motor_sports
    Auto Racing                              → sports.motor_sports.auto_racing
  Water Sports                               → sports.water_sports
    Surfing                                  → sports.water_sports.surfing
  Winter Sports                              → sports.winter_sports
    Skiing & Snowboarding                    → sports.winter_sports.skiing_and_snowboarding
  Sports Coaching & Training                 → sports.sports_coaching_and_training
  Sports Fan Gear & Apparel                  → sports.sports_fan_gear_and_apparel
```

### 20. Travel & Transportation

```
Travel & Transportation                      → travel_and_transportation
```

*(Subcategories to be added after Postman exploration)*

---

## Category Count Summary

| Category | L1 | L2 | L3 | Total |
|----------|----|----|-----|-------|
| Arts & Entertainment | 1 | 8 | 11 | 20 |
| Autos & Vehicles | 1 | 0 | 0 | 1 |
| Beauty & Fitness | 1 | 0 | 0 | 1 |
| Books & Literature | 1 | 0 | 0 | 1 |
| Business & Industrial | 1 | 4 | 2 | 7 |
| Computers & Electronics | 1 | 0 | 0 | 1 |
| Finance | 1 | 0 | 0 | 1 |
| Food & Drink | 1 | 3 | 0 | 4 |
| Games | 1 | 4 | 4 | 9 |
| Health | 1 | 5 | 5 | 11 |
| Hobbies & Leisure | 1 | 4 | 4 | 9 |
| Home & Garden | 1 | 0 | 0 | 1 |
| Jobs & Education | 1 | 0 | 0 | 1 |
| Law & Government | 1 | 0 | 0 | 1 |
| News | 1 | 6 | 2 | 9 |
| People & Society | 1 | 4 | 4 | 9 |
| Pets & Animals | 1 | 0 | 0 | 1 |
| Science | 1 | 5 | 2 | 8 |
| Sports | 1 | 7 | 7 | 15 |
| Travel & Transportation | 1 | 0 | 0 | 1 |
| **TOTAL** | **20** | **50** | **41** | **111** |

**Note:** Categories marked "(Subcategories to be added after Postman exploration)" will be expanded after user confirms findings.

---

## Custom Category Slots

User-created categories coexist in the same `categories` table with `source = 'user'`. They can:

1. **Extend an existing branch** — e.g., add `sports.team_sports.american_football.nfl` or `sports.team_sports.american_football.nfl.chiefs` for team-level granularity
2. **Add cross-cutting categories** — e.g., `reactions`, `reactions.hot_takes` for meta-commentary content
3. **Add missing niches** — e.g., `arts_and_entertainment.podcasts` if Perspectize expands beyond YouTube

Example custom category inserts (after seed):
```sql
INSERT INTO categories (name, slug, path, source) VALUES
('NFL', 'nfl', 'sports.team_sports.american_football.nfl', 'user'),
('Kansas City Chiefs', 'kansas-city-chiefs', 'sports.team_sports.american_football.nfl.kansas_city_chiefs', 'user'),
('Reactions', 'reactions', 'reactions', 'user');
```

---

## ltree Path Transformation Rules

Google NL path → ltree path:
1. Remove leading `/`
2. Replace `/` with `.`
3. Replace ` & ` with `_and_`
4. Replace spaces with `_`
5. Remove commas
6. Lowercase
7. Remove any remaining non-alphanumeric characters except `_` and `-`

| Google NL Path | ltree Path |
|----------------|------------|
| `/Sports` | `sports` |
| `/Sports/Team Sports/American Football` | `sports.team_sports.american_football` |
| `/Arts & Entertainment/Music & Audio` | `arts_and_entertainment.music_and_audio` |
| `/Science/Computer Science/Machine Learning & Artificial Intelligence` | `science.computer_science.machine_learning_and_artificial_intelligence` |
| `/Arts & Entertainment/Events & Listings/Bars, Clubs & Nightlife` | `arts_and_entertainment.events_and_listings.bars_clubs_and_nightlife` |
