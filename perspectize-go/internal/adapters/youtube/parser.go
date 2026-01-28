package youtube

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ExtractVideoID extracts the video ID from various YouTube URL formats
func ExtractVideoID(url string) (string, error) {
	patterns := []string{
		`(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/embed/|youtube\.com/v/)([a-zA-Z0-9_-]{11})`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("could not extract video ID from URL: %s", url)
}

// ParseISO8601Duration converts ISO 8601 duration string to seconds
// Examples: PT1H30M45S -> 5445, PT10M -> 600, PT45S -> 45
func ParseISO8601Duration(duration string) (int, error) {
	if !strings.HasPrefix(duration, "PT") {
		return 0, fmt.Errorf("invalid duration format: %s", duration)
	}

	duration = strings.TrimPrefix(duration, "PT")

	var hours, minutes, seconds int
	var err error

	// Parse hours
	if idx := strings.Index(duration, "H"); idx != -1 {
		hours, err = strconv.Atoi(duration[:idx])
		if err != nil {
			return 0, fmt.Errorf("invalid hours in duration: %w", err)
		}
		duration = duration[idx+1:]
	}

	// Parse minutes
	if idx := strings.Index(duration, "M"); idx != -1 {
		minutes, err = strconv.Atoi(duration[:idx])
		if err != nil {
			return 0, fmt.Errorf("invalid minutes in duration: %w", err)
		}
		duration = duration[idx+1:]
	}

	// Parse seconds
	if idx := strings.Index(duration, "S"); idx != -1 {
		seconds, err = strconv.Atoi(duration[:idx])
		if err != nil {
			return 0, fmt.Errorf("invalid seconds in duration: %w", err)
		}
	}

	return hours*3600 + minutes*60 + seconds, nil
}
