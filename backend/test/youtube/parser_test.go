package youtube_test

import (
	"testing"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/youtube"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- ExtractVideoID Tests ---

func TestExtractVideoID_StandardURL(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

func TestExtractVideoID_ShortURL(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://youtu.be/dQw4w9WgXcQ")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

func TestExtractVideoID_EmbedURL(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://www.youtube.com/embed/dQw4w9WgXcQ")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

func TestExtractVideoID_OldStyleURL(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://www.youtube.com/v/dQw4w9WgXcQ")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

func TestExtractVideoID_WithExtraParams(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://www.youtube.com/watch?v=dQw4w9WgXcQ&list=PLrAXtmErZgOeiKm4sgNOknGvNjby9efdf")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

func TestExtractVideoID_WithHyphenAndUnderscore(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://www.youtube.com/watch?v=a1B-c2D_e3F")
	require.NoError(t, err)
	assert.Equal(t, "a1B-c2D_e3F", id)
}

func TestExtractVideoID_InvalidURL(t *testing.T) {
	_, err := youtube.ExtractVideoID("https://www.example.com/video")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not extract video ID")
}

func TestExtractVideoID_EmptyString(t *testing.T) {
	_, err := youtube.ExtractVideoID("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not extract video ID")
}

func TestExtractVideoID_RandomText(t *testing.T) {
	_, err := youtube.ExtractVideoID("not a url at all")
	require.Error(t, err)
}

func TestExtractVideoID_HTTPWithoutS(t *testing.T) {
	id, err := youtube.ExtractVideoID("http://www.youtube.com/watch?v=dQw4w9WgXcQ")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

func TestExtractVideoID_ShortsURL(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://www.youtube.com/shorts/dQw4w9WgXcQ")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

func TestExtractVideoID_LiveURL(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://www.youtube.com/live/dQw4w9WgXcQ")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

func TestExtractVideoID_ShortEmbedURL(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://www.youtube.com/e/dQw4w9WgXcQ")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

func TestExtractVideoID_NoCookieEmbedURL(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://www.youtube-nocookie.com/embed/dQw4w9WgXcQ")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

func TestExtractVideoID_MusicYouTubeURL(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://music.youtube.com/watch?v=dQw4w9WgXcQ")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

func TestExtractVideoID_MobileURL(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://m.youtube.com/watch?v=dQw4w9WgXcQ")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

func TestExtractVideoID_ShortsWithParams(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://youtube.com/shorts/dQw4w9WgXcQ?si=abc123")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

func TestExtractVideoID_LiveWithParams(t *testing.T) {
	id, err := youtube.ExtractVideoID("https://youtube.com/live/dQw4w9WgXcQ?feature=share")
	require.NoError(t, err)
	assert.Equal(t, "dQw4w9WgXcQ", id)
}

// --- ParseISO8601Duration Tests ---

func TestParseISO8601Duration_HoursMinutesSeconds(t *testing.T) {
	seconds, err := youtube.ParseISO8601Duration("PT1H30M45S")
	require.NoError(t, err)
	assert.Equal(t, 5445, seconds) // 1*3600 + 30*60 + 45
}

func TestParseISO8601Duration_MinutesOnly(t *testing.T) {
	seconds, err := youtube.ParseISO8601Duration("PT10M")
	require.NoError(t, err)
	assert.Equal(t, 600, seconds)
}

func TestParseISO8601Duration_SecondsOnly(t *testing.T) {
	seconds, err := youtube.ParseISO8601Duration("PT45S")
	require.NoError(t, err)
	assert.Equal(t, 45, seconds)
}

func TestParseISO8601Duration_HoursOnly(t *testing.T) {
	seconds, err := youtube.ParseISO8601Duration("PT2H")
	require.NoError(t, err)
	assert.Equal(t, 7200, seconds)
}

func TestParseISO8601Duration_HoursAndMinutes(t *testing.T) {
	seconds, err := youtube.ParseISO8601Duration("PT1H15M")
	require.NoError(t, err)
	assert.Equal(t, 4500, seconds) // 1*3600 + 15*60
}

func TestParseISO8601Duration_MinutesAndSeconds(t *testing.T) {
	seconds, err := youtube.ParseISO8601Duration("PT5M30S")
	require.NoError(t, err)
	assert.Equal(t, 330, seconds) // 5*60 + 30
}

func TestParseISO8601Duration_HoursAndSeconds(t *testing.T) {
	seconds, err := youtube.ParseISO8601Duration("PT1H10S")
	require.NoError(t, err)
	assert.Equal(t, 3610, seconds) // 1*3600 + 10
}

func TestParseISO8601Duration_ZeroDuration(t *testing.T) {
	seconds, err := youtube.ParseISO8601Duration("PT0S")
	require.NoError(t, err)
	assert.Equal(t, 0, seconds)
}

func TestParseISO8601Duration_InvalidPrefix(t *testing.T) {
	_, err := youtube.ParseISO8601Duration("P1H30M")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid duration format")
}

func TestParseISO8601Duration_EmptyString(t *testing.T) {
	_, err := youtube.ParseISO8601Duration("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid duration format")
}

func TestParseISO8601Duration_NoComponents(t *testing.T) {
	// PT with nothing after it - returns 0
	seconds, err := youtube.ParseISO8601Duration("PT")
	require.NoError(t, err)
	assert.Equal(t, 0, seconds)
}

func TestParseISO8601Duration_LargeValues(t *testing.T) {
	seconds, err := youtube.ParseISO8601Duration("PT10H59M59S")
	require.NoError(t, err)
	assert.Equal(t, 39599, seconds) // 10*3600 + 59*60 + 59
}
