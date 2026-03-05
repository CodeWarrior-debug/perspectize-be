package youtube

import (
	"fmt"
	"strings"
	"testing"
)

func TestSanitizeYouTubeError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		contains    string
		notContains string
	}{
		{
			name:        "nil error returns empty string",
			err:         nil,
			contains:    "",
			notContains: "",
		},
		{
			name:        "URL with API key removed",
			err:         fmt.Errorf("Get https://www.googleapis.com/youtube/v3/videos?key=SECRET123&id=test: status code 403"),
			contains:    "YouTube API error",
			notContains: "SECRET123",
		},
		{
			name:        "Generic googleapis error",
			err:         fmt.Errorf("request to googleapis.com failed"),
			contains:    "YouTube API request failed",
			notContains: "googleapis.com",
		},
		{
			name:        "Non-YouTube error unchanged",
			err:         fmt.Errorf("network timeout"),
			contains:    "network timeout",
			notContains: "",
		},
		{
			name:        "URL with key in middle of message",
			err:         fmt.Errorf("dial tcp: lookup www.googleapis.com/youtube/v3/videos?key=MYSECRETKEY123"),
			contains:    "YouTube API request failed",
			notContains: "MYSECRETKEY123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeYouTubeError(tt.err)

			if tt.err == nil {
				if result != "" {
					t.Errorf("expected empty string for nil error, got %q", result)
				}
				return
			}

			if tt.contains != "" && !strings.Contains(result, tt.contains) {
				t.Errorf("expected result to contain %q, got %q", tt.contains, result)
			}
			if tt.notContains != "" && strings.Contains(result, tt.notContains) {
				t.Errorf("result should not contain %q, got %q", tt.notContains, result)
			}
		})
	}
}
