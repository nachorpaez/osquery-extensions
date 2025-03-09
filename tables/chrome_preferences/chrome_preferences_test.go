package chrome_preferences

import (
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/nachorpaez/osquery-extensions/pkg/utils"
	"github.com/stretchr/testify/assert"
)

//go:embed test_Preferences
var testPreferences []byte

func TestParsePreferences(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a test Chrome preferences file
	profilePath := filepath.Join(tempDir, "Preferences")
	err := os.WriteFile(profilePath, testPreferences, 0600)
	assert.NoError(t, err)

	// Build a ChromeProfilePath to reflect the new function signature
	chromeProfile := utils.ChromeProfilePath{
		UserName: "user1",            // Matches old "User" field
		Value:    tempDir,            // The containing directory
		Type:     utils.GoogleChrome, // Arbitrary choice of browser type
	}

	results, err := parsePreferences(context.Background(), chromeProfile)
	assert.NoError(t, err)
	assert.Len(t, results, 5)

	expectedRows := []map[string]string{
		{
			"category":      "geolocation",
			"url":           "https://test.com:443,*",
			"expiration":    "",
			"last_modified": "13359398812209250",
			"model":         "0",
			"setting":       "2",
			"profile_path":  tempDir,
			"user":          "user1",
			"browser_type":  "chrome", // from utils.GetChromeBrowserName(utils.GoogleChrome)
		},
		{
			"category":      "notifications",
			"url":           "https://meet.google.com:443,*",
			"expiration":    "",
			"last_modified": "13357248805228331",
			"model":         "0",
			"setting":       "1",
			"profile_path":  tempDir,
			"user":          "user1",
			"browser_type":  "chrome",
		},
		{
			"category":      "media_stream_camera",
			"url":           "https://meet.google.com:443,*",
			"expiration":    "",
			"last_modified": "13357248819111798",
			"model":         "0",
			"setting":       "1",
			"profile_path":  tempDir,
			"user":          "user1",
			"browser_type":  "chrome",
		},
		{
			"category":      "media_stream_mic",
			"url":           "https://meet.google.com:443,*",
			"expiration":    "",
			"last_modified": "13357248797611836",
			"model":         "1",
			"setting":       "1",
			"profile_path":  tempDir,
			"user":          "user1",
			"browser_type":  "chrome",
		},
		{
			"category":      "popups",
			"url":           "https://www.digicert.com:443,*",
			"expiration":    "",
			"last_modified": "",
			"model":         "0",
			"setting":       "1",
			"profile_path":  tempDir,
			"user":          "user1",
			"browser_type":  "chrome",
		},
	}

	// Use ElementsMatch if the order doesn't matter; otherwise use Equal
	assert.ElementsMatch(t, expectedRows, results)
}
