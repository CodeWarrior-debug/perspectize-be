-- =============================================================================
-- Perspectize Category Seed Data
-- Source: Google Cloud Natural Language API V2 taxonomy (curated subset)
-- Generated: 2026-02-20
-- Status: Initial draft — pending user review of CURATED-CATEGORIES.md
-- =============================================================================
-- Usage:
--   psql $DATABASE_URL -f categories-seed.sql
-- Prerequisites:
--   PostgreSQL with ltree extension available
--   Run AFTER your schema migration creates the categories table
-- =============================================================================

CREATE EXTENSION IF NOT EXISTS ltree;

-- CREATE TABLE categories (schema reference — use IF NOT EXISTS for idempotency)
CREATE TABLE IF NOT EXISTS categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    path        ltree NOT NULL,
    source      TEXT NOT NULL DEFAULT 'google_nl',
    google_path TEXT,
    depth       INT GENERATED ALWAYS AS (nlevel(path)) STORED,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_categories_path_gist ON categories USING GIST (path);
CREATE INDEX IF NOT EXISTS idx_categories_path_btree ON categories USING BTREE (path);
CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_path_unique ON categories (path);

BEGIN;

-- =============================================================================
-- ARTS & ENTERTAINMENT
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Arts & Entertainment',        'arts-and-entertainment',                       'arts_and_entertainment',                                                                   'google_nl', '/Arts & Entertainment');

-- Arts & Entertainment > Comics & Animation
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Comics & Animation',          'comics-and-animation',                         'arts_and_entertainment.comics_and_animation',                                              'google_nl', '/Arts & Entertainment/Comics & Animation'),
('Anime & Manga',               'anime-and-manga',                              'arts_and_entertainment.comics_and_animation.anime_and_manga',                              'google_nl', '/Arts & Entertainment/Comics & Animation/Anime & Manga'),
('Cartoons',                    'cartoons',                                     'arts_and_entertainment.comics_and_animation.cartoons',                                     'google_nl', '/Arts & Entertainment/Comics & Animation/Cartoons'),
('Comics',                      'comics',                                       'arts_and_entertainment.comics_and_animation.comics',                                       'google_nl', '/Arts & Entertainment/Comics & Animation/Comics');

-- Arts & Entertainment > Entertainment Industry
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Entertainment Industry',      'entertainment-industry',                       'arts_and_entertainment.entertainment_industry',                                            'google_nl', '/Arts & Entertainment/Entertainment Industry'),
('Film & TV Industry',          'film-and-tv-industry',                         'arts_and_entertainment.entertainment_industry.film_and_tv_industry',                       'google_nl', '/Arts & Entertainment/Entertainment Industry/Film & TV Industry'),
('Recording Industry',          'recording-industry',                           'arts_and_entertainment.entertainment_industry.recording_industry',                         'google_nl', '/Arts & Entertainment/Entertainment Industry/Recording Industry');

-- Arts & Entertainment > Events & Listings
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Events & Listings',           'events-and-listings',                          'arts_and_entertainment.events_and_listings',                                               'google_nl', '/Arts & Entertainment/Events & Listings'),
('Bars, Clubs & Nightlife',     'bars-clubs-and-nightlife',                     'arts_and_entertainment.events_and_listings.bars_clubs_and_nightlife',                      'google_nl', '/Arts & Entertainment/Events & Listings/Bars, Clubs & Nightlife'),
('Concerts & Music Festivals',  'concerts-and-music-festivals',                 'arts_and_entertainment.events_and_listings.concerts_and_music_festivals',                  'google_nl', '/Arts & Entertainment/Events & Listings/Concerts & Music Festivals');

-- Arts & Entertainment > Movies
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Movies',                      'movies',                                       'arts_and_entertainment.movies',                                                            'google_nl', '/Arts & Entertainment/Movies'),
('Action & Adventure Films',    'action-and-adventure-films',                   'arts_and_entertainment.movies.action_and_adventure_films',                                 'google_nl', '/Arts & Entertainment/Movies/Action & Adventure Films');

-- Arts & Entertainment > Music & Audio
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Music & Audio',               'music-and-audio',                              'arts_and_entertainment.music_and_audio',                                                   'google_nl', '/Arts & Entertainment/Music & Audio'),
('Classical Music',             'classical-music',                              'arts_and_entertainment.music_and_audio.classical_music',                                   'google_nl', '/Arts & Entertainment/Music & Audio/Classical Music'),
('Jazz & Blues',                'jazz-and-blues',                               'arts_and_entertainment.music_and_audio.jazz_and_blues',                                    'google_nl', '/Arts & Entertainment/Music & Audio/Jazz & Blues');

-- Arts & Entertainment > Performing Arts
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Performing Arts',             'performing-arts',                              'arts_and_entertainment.performing_arts',                                                   'google_nl', '/Arts & Entertainment/Performing Arts'),
('Dance',                       'dance',                                        'arts_and_entertainment.performing_arts.dance',                                             'google_nl', '/Arts & Entertainment/Performing Arts/Dance');

-- Arts & Entertainment > TV & Video
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('TV & Video',                  'tv-and-video',                                 'arts_and_entertainment.tv_and_video',                                                      'google_nl', '/Arts & Entertainment/TV & Video'),
('TV Shows & Programs',         'tv-shows-and-programs',                        'arts_and_entertainment.tv_and_video.tv_shows_and_programs',                                'google_nl', '/Arts & Entertainment/TV & Video/TV Shows & Programs');

-- Arts & Entertainment > Visual Art & Design
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Visual Art & Design',         'visual-art-and-design',                        'arts_and_entertainment.visual_art_and_design',                                             'google_nl', '/Arts & Entertainment/Visual Art & Design'),
('Architecture',                'architecture',                                 'arts_and_entertainment.visual_art_and_design.architecture',                                'google_nl', '/Arts & Entertainment/Visual Art & Design/Architecture');

-- =============================================================================
-- AUTOS & VEHICLES
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Autos & Vehicles',            'autos-and-vehicles',                           'autos_and_vehicles',                                                                       'google_nl', '/Autos & Vehicles');

-- =============================================================================
-- BEAUTY & FITNESS
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Beauty & Fitness',            'beauty-and-fitness',                           'beauty_and_fitness',                                                                       'google_nl', '/Beauty & Fitness');

-- =============================================================================
-- BOOKS & LITERATURE
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Books & Literature',          'books-and-literature',                         'books_and_literature',                                                                     'google_nl', '/Books & Literature');

-- =============================================================================
-- BUSINESS & INDUSTRIAL
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Business & Industrial',       'business-and-industrial',                      'business_and_industrial',                                                                  'google_nl', '/Business & Industrial');

-- Business & Industrial > Advertising & Marketing
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Advertising & Marketing',     'advertising-and-marketing',                    'business_and_industrial.advertising_and_marketing',                                        'google_nl', '/Business & Industrial/Advertising & Marketing'),
('Public Relations',            'public-relations',                             'business_and_industrial.advertising_and_marketing.public_relations',                       'google_nl', '/Business & Industrial/Advertising & Marketing/Public Relations');

-- Business & Industrial > Energy & Utilities
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Energy & Utilities',          'energy-and-utilities',                         'business_and_industrial.energy_and_utilities',                                             'google_nl', '/Business & Industrial/Energy & Utilities'),
('Oil & Gas',                   'oil-and-gas',                                  'business_and_industrial.energy_and_utilities.oil_and_gas',                                 'google_nl', '/Business & Industrial/Energy & Utilities/Oil & Gas');

-- Business & Industrial (other L2)
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Manufacturing',               'manufacturing',                                'business_and_industrial.manufacturing',                                                    'google_nl', '/Business & Industrial/Manufacturing'),
('Pharmaceuticals & Biotech',   'pharmaceuticals-and-biotech',                  'business_and_industrial.pharmaceuticals_and_biotech',                                      'google_nl', '/Business & Industrial/Pharmaceuticals & Biotech');

-- =============================================================================
-- COMPUTERS & ELECTRONICS
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Computers & Electronics',     'computers-and-electronics',                    'computers_and_electronics',                                                                'google_nl', '/Computers & Electronics');

-- =============================================================================
-- FINANCE
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Finance',                     'finance',                                      'finance',                                                                                  'google_nl', '/Finance');

-- =============================================================================
-- FOOD & DRINK
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Food & Drink',                'food-and-drink',                               'food_and_drink',                                                                           'google_nl', '/Food & Drink'),
('Cooking & Recipes',           'cooking-and-recipes',                          'food_and_drink.cooking_and_recipes',                                                       'google_nl', '/Food & Drink/Cooking & Recipes'),
('Beverages',                   'beverages',                                    'food_and_drink.beverages',                                                                 'google_nl', '/Food & Drink/Beverages'),
('Restaurants',                 'restaurants',                                  'food_and_drink.restaurants',                                                               'google_nl', '/Food & Drink/Restaurants');

-- =============================================================================
-- GAMES
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Games',                       'games',                                        'games',                                                                                    'google_nl', '/Games');

-- Games > Computer & Video Games
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Computer & Video Games',      'computer-and-video-games',                     'games.computer_and_video_games',                                                           'google_nl', '/Games/Computer & Video Games'),
('Shooter Games',               'shooter-games',                                'games.computer_and_video_games.shooter_games',                                             'google_nl', '/Games/Computer & Video Games/Shooter Games');

-- Games > Board Games
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Board Games',                 'board-games',                                  'games.board_games',                                                                        'google_nl', '/Games/Board Games'),
('Chess & Abstract Strategy Games', 'chess-and-abstract-strategy-games',        'games.board_games.chess_and_abstract_strategy_games',                                      'google_nl', '/Games/Board Games/Chess & Abstract Strategy Games');

-- Games > Card Games
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Card Games',                  'card-games',                                   'games.card_games',                                                                         'google_nl', '/Games/Card Games'),
('Poker & Casino Games',        'poker-and-casino-games',                       'games.card_games.poker_and_casino_games',                                                  'google_nl', '/Games/Card Games/Poker & Casino Games');

-- Games > Gambling
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Gambling',                    'gambling',                                     'games.gambling',                                                                           'google_nl', '/Games/Gambling'),
('Lottery',                     'lottery',                                      'games.gambling.lottery',                                                                   'google_nl', '/Games/Gambling/Lottery');

-- =============================================================================
-- HEALTH
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Health',                      'health',                                       'health',                                                                                   'google_nl', '/Health');

-- Health > Health Conditions
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Health Conditions',           'health-conditions',                            'health.health_conditions',                                                                 'google_nl', '/Health/Health Conditions'),
('Cancer',                      'cancer',                                       'health.health_conditions.cancer',                                                          'google_nl', '/Health/Health Conditions/Cancer');

-- Health > Mental Health
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Mental Health',               'mental-health',                                'health.mental_health',                                                                     'google_nl', '/Health/Mental Health'),
('Depression',                  'depression',                                   'health.mental_health.depression',                                                          'google_nl', '/Health/Mental Health/Depression');

-- Health > Medical Facilities & Services
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Medical Facilities & Services', 'medical-facilities-and-services',            'health.medical_facilities_and_services',                                                   'google_nl', '/Health/Medical Facilities & Services'),
('Hospitals & Treatment Centers', 'hospitals-and-treatment-centers',            'health.medical_facilities_and_services.hospitals_and_treatment_centers',                   'google_nl', '/Health/Medical Facilities & Services/Hospitals & Treatment Centers');

-- Health > Nutrition
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Nutrition',                   'nutrition',                                    'health.nutrition',                                                                         'google_nl', '/Health/Nutrition'),
('Vitamins & Supplements',      'vitamins-and-supplements',                     'health.nutrition.vitamins_and_supplements',                                                'google_nl', '/Health/Nutrition/Vitamins & Supplements');

-- Health > Vision Care
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Vision Care',                 'vision-care',                                  'health.vision_care',                                                                       'google_nl', '/Health/Vision Care'),
('Eyeglasses & Contacts',       'eyeglasses-and-contacts',                      'health.vision_care.eyeglasses_and_contacts',                                               'google_nl', '/Health/Vision Care/Eyeglasses & Contacts');

-- =============================================================================
-- HOBBIES & LEISURE
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Hobbies & Leisure',           'hobbies-and-leisure',                          'hobbies_and_leisure',                                                                      'google_nl', '/Hobbies & Leisure');

-- Hobbies & Leisure > Crafts
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Crafts',                      'crafts',                                       'hobbies_and_leisure.crafts',                                                               'google_nl', '/Hobbies & Leisure/Crafts'),
('Ceramics & Pottery',          'ceramics-and-pottery',                         'hobbies_and_leisure.crafts.ceramics_and_pottery',                                          'google_nl', '/Hobbies & Leisure/Crafts/Ceramics & Pottery');

-- Hobbies & Leisure > Outdoors
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Outdoors',                    'outdoors',                                     'hobbies_and_leisure.outdoors',                                                             'google_nl', '/Hobbies & Leisure/Outdoors'),
('Fishing',                     'fishing',                                      'hobbies_and_leisure.outdoors.fishing',                                                     'google_nl', '/Hobbies & Leisure/Outdoors/Fishing');

-- Hobbies & Leisure > Special Occasions
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Special Occasions',           'special-occasions',                            'hobbies_and_leisure.special_occasions',                                                    'google_nl', '/Hobbies & Leisure/Special Occasions'),
('Weddings',                    'weddings',                                     'hobbies_and_leisure.special_occasions.weddings',                                           'google_nl', '/Hobbies & Leisure/Special Occasions/Weddings');

-- Hobbies & Leisure > Water Activities
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Water Activities',            'water-activities',                             'hobbies_and_leisure.water_activities',                                                     'google_nl', '/Hobbies & Leisure/Water Activities'),
('Boating',                     'boating',                                      'hobbies_and_leisure.water_activities.boating',                                             'google_nl', '/Hobbies & Leisure/Water Activities/Boating');

-- =============================================================================
-- HOME & GARDEN
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Home & Garden',               'home-and-garden',                              'home_and_garden',                                                                          'google_nl', '/Home & Garden');

-- =============================================================================
-- JOBS & EDUCATION
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Jobs & Education',            'jobs-and-education',                           'jobs_and_education',                                                                       'google_nl', '/Jobs & Education');

-- =============================================================================
-- LAW & GOVERNMENT
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Law & Government',            'law-and-government',                           'law_and_government',                                                                       'google_nl', '/Law & Government');

-- =============================================================================
-- NEWS
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('News',                        'news',                                         'news',                                                                                     'google_nl', '/News');

-- News > Business News
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Business News',               'business-news',                                'news.business_news',                                                                       'google_nl', '/News/Business News'),
('Company News',                'company-news',                                 'news.business_news.company_news',                                                          'google_nl', '/News/Business News/Company News');

-- News (other L2)
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Health News',                 'health-news',                                  'news.health_news',                                                                         'google_nl', '/News/Health News');

-- News > Politics
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Politics',                    'politics',                                     'news.politics',                                                                            'google_nl', '/News/Politics'),
('Campaigns & Elections',       'campaigns-and-elections',                      'news.politics.campaigns_and_elections',                                                    'google_nl', '/News/Politics/Campaigns & Elections');

-- News (other L2)
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Sports News',                 'sports-news',                                  'news.sports_news',                                                                         'google_nl', '/News/Sports News'),
('Technology News',             'technology-news',                              'news.technology_news',                                                                     'google_nl', '/News/Technology News'),
('Weather',                     'weather',                                      'news.weather',                                                                             'google_nl', '/News/Weather');

-- =============================================================================
-- PEOPLE & SOCIETY
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('People & Society',            'people-and-society',                           'people_and_society',                                                                       'google_nl', '/People & Society');

-- People & Society > Family & Relationships
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Family & Relationships',      'family-and-relationships',                     'people_and_society.family_and_relationships',                                              'google_nl', '/People & Society/Family & Relationships'),
('Marriage',                    'marriage',                                     'people_and_society.family_and_relationships.marriage',                                     'google_nl', '/People & Society/Family & Relationships/Marriage');

-- People & Society > Kids & Teens
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Kids & Teens',                'kids-and-teens',                               'people_and_society.kids_and_teens',                                                        'google_nl', '/People & Society/Kids & Teens'),
("Children's Interests",        'childrens-interests',                          'people_and_society.kids_and_teens.childrens_interests',                                    'google_nl', "/People & Society/Kids & Teens/Children's Interests");

-- People & Society (other L2)
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Religion & Belief',           'religion-and-belief',                          'people_and_society.religion_and_belief',                                                   'google_nl', '/People & Society/Religion & Belief');

-- People & Society > Social Issues & Advocacy
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Social Issues & Advocacy',    'social-issues-and-advocacy',                   'people_and_society.social_issues_and_advocacy',                                            'google_nl', '/People & Society/Social Issues & Advocacy'),
('Charity & Philanthropy',      'charity-and-philanthropy',                     'people_and_society.social_issues_and_advocacy.charity_and_philanthropy',                   'google_nl', '/People & Society/Social Issues & Advocacy/Charity & Philanthropy');

-- =============================================================================
-- PETS & ANIMALS
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Pets & Animals',              'pets-and-animals',                             'pets_and_animals',                                                                         'google_nl', '/Pets & Animals');

-- =============================================================================
-- SCIENCE
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Science',                     'science',                                      'science',                                                                                  'google_nl', '/Science'),
('Astronomy',                   'astronomy',                                    'science.astronomy',                                                                        'google_nl', '/Science/Astronomy');

-- Science > Biology
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Biology',                     'biology',                                      'science.biology',                                                                          'google_nl', '/Science/Biology'),
('Neuroscience',                'neuroscience',                                 'science.biology.neuroscience',                                                             'google_nl', '/Science/Biology/Neuroscience');

-- Science (other L2)
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Chemistry',                   'chemistry',                                    'science.chemistry',                                                                        'google_nl', '/Science/Chemistry');

-- Science > Computer Science
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Computer Science',            'computer-science',                             'science.computer_science',                                                                 'google_nl', '/Science/Computer Science'),
('Machine Learning & Artificial Intelligence', 'machine-learning-and-ai',       'science.computer_science.machine_learning_and_artificial_intelligence',                    'google_nl', '/Science/Computer Science/Machine Learning & Artificial Intelligence');

-- Science (other L2)
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Physics',                     'physics',                                      'science.physics',                                                                          'google_nl', '/Science/Physics');

-- =============================================================================
-- SPORTS
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Sports',                      'sports',                                       'sports',                                                                                   'google_nl', '/Sports');

-- Sports > Team Sports
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Team Sports',                 'team-sports',                                  'sports.team_sports',                                                                       'google_nl', '/Sports/Team Sports'),
('American Football',           'american-football',                            'sports.team_sports.american_football',                                                     'google_nl', '/Sports/Team Sports/American Football'),
('Baseball',                    'baseball',                                     'sports.team_sports.baseball',                                                              'google_nl', '/Sports/Team Sports/Baseball'),
('Basketball',                  'basketball',                                   'sports.team_sports.basketball',                                                            'google_nl', '/Sports/Team Sports/Basketball'),
('Hockey',                      'hockey',                                       'sports.team_sports.hockey',                                                                'google_nl', '/Sports/Team Sports/Hockey'),
('Soccer',                      'soccer',                                       'sports.team_sports.soccer',                                                                'google_nl', '/Sports/Team Sports/Soccer');

-- Sports > Individual Sports
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Individual Sports',           'individual-sports',                            'sports.individual_sports',                                                                 'google_nl', '/Sports/Individual Sports'),
('Golf',                        'golf',                                         'sports.individual_sports.golf',                                                            'google_nl', '/Sports/Individual Sports/Golf'),
('Cycling',                     'cycling',                                      'sports.individual_sports.cycling',                                                         'google_nl', '/Sports/Individual Sports/Cycling');

-- Sports > Motor Sports
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Motor Sports',                'motor-sports',                                 'sports.motor_sports',                                                                      'google_nl', '/Sports/Motor Sports'),
('Auto Racing',                 'auto-racing',                                  'sports.motor_sports.auto_racing',                                                          'google_nl', '/Sports/Motor Sports/Auto Racing');

-- Sports > Water Sports
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Water Sports',                'water-sports',                                 'sports.water_sports',                                                                      'google_nl', '/Sports/Water Sports'),
('Surfing',                     'surfing',                                      'sports.water_sports.surfing',                                                              'google_nl', '/Sports/Water Sports/Surfing');

-- Sports > Winter Sports
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Winter Sports',               'winter-sports',                                'sports.winter_sports',                                                                     'google_nl', '/Sports/Winter Sports'),
('Skiing & Snowboarding',       'skiing-and-snowboarding',                      'sports.winter_sports.skiing_and_snowboarding',                                             'google_nl', '/Sports/Winter Sports/Skiing & Snowboarding');

-- Sports (other L2)
INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Sports Coaching & Training',  'sports-coaching-and-training',                 'sports.sports_coaching_and_training',                                                      'google_nl', '/Sports/Sports Coaching & Training'),
('Sports Fan Gear & Apparel',   'sports-fan-gear-and-apparel',                  'sports.sports_fan_gear_and_apparel',                                                       'google_nl', '/Sports/Sports Fan Gear & Apparel');

-- =============================================================================
-- TRAVEL & TRANSPORTATION
-- =============================================================================

INSERT INTO categories (name, slug, path, source, google_path) VALUES
('Travel & Transportation',     'travel-and-transportation',                    'travel_and_transportation',                                                                'google_nl', '/Travel & Transportation');

COMMIT;

-- =============================================================================
-- Verification query (run after seeding to confirm counts)
-- =============================================================================
-- SELECT
--   nlevel(path) AS depth,
--   COUNT(*) AS category_count
-- FROM categories
-- WHERE source = 'google_nl'
-- GROUP BY depth
-- ORDER BY depth;
--
-- Expected output:
--   depth | category_count
--   ------+---------------
--       1 |            20
--       2 |            ~50
--       3 |            ~41
