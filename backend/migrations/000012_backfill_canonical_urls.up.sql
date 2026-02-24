-- Normalize existing YouTube URLs to canonical form.
-- Extracts the video ID using regex and rebuilds as canonical watch URL.
-- Only updates rows where the URL is not already in canonical form.
UPDATE content
SET url = 'https://www.youtube.com/watch?v=' || (
    CASE
        -- youtu.be/ID or youtu.be/ID?params
        WHEN url ~ 'youtu\.be/([a-zA-Z0-9_-]{11})' THEN
            (regexp_match(url, 'youtu\.be/([a-zA-Z0-9_-]{11})'))[1]
        -- youtube.com/watch?v=ID (any subdomain)
        WHEN url ~ 'youtube\.com/watch\?v=([a-zA-Z0-9_-]{11})' THEN
            (regexp_match(url, 'youtube\.com/watch\?v=([a-zA-Z0-9_-]{11})'))[1]
        -- embed, v, e, shorts, live paths
        WHEN url ~ 'youtube(?:-nocookie)?\.com/(?:embed|v|e|shorts|live)/([a-zA-Z0-9_-]{11})' THEN
            (regexp_match(url, 'youtube(?:-nocookie)?\.com/(?:embed|v|e|shorts|live)/([a-zA-Z0-9_-]{11})'))[1]
    END
)
WHERE content_type = 'youtube'
  AND url IS NOT NULL
  AND url NOT LIKE 'https://www.youtube.com/watch?v=___________'
  AND url ~ '(youtu\.be/|youtube\.com/)';
