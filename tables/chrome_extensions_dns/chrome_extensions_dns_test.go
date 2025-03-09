package chrome_extensions_dns

import (
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/nachorpaez/osquery-extensions/pkg/utils"
	"github.com/stretchr/testify/assert"
)

//go:embed test_NetworkPersistentState
var testNetworkPersistentState []byte

func TestAnalyzeNetworkState(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Write the test network state file into the tempDir
	stateFilePath := filepath.Join(tempDir, "Network Persistent State")
	err := os.WriteFile(stateFilePath, testNetworkPersistentState, 0600)
	assert.NoError(t, err)

	// Build a mock ChromeProfilePath
	mockProfile := utils.ChromeProfilePath{
		UserName: "testuser",
		Type:     utils.GoogleChrome, // see what your code expects
		Value:    tempDir,            // directory containing the "Network Persistent State" file
	}

	// Call the function we want to test
	results, err := analyzeNetworkState(context.Background(), mockProfile)
	assert.NoError(t, err, "analyzeNetworkState should not return an error")

	// We expect 3 entries from the mock file:
	//   2 servers with "Active" type
	//   1 broken alt service with "Broken" type
	assert.Len(t, results, 3)

	expected := []map[string]string{
		{
			"browser_type": "chrome",
			"profile":      filepath.Base(tempDir),
			"extension_id": "aeblfdkhhhdcdjpifhhbdiojplfjncoa\x00",
			"domain":       "my.1password.com",
			"type":         "Active",
			"user":         "testuser",
		},
		{
			"browser_type": "chrome",
			"profile":      filepath.Base(tempDir),
			"extension_id": "dgjhfomjieaadpoljlnidmbgkdffpack\x00",
			"domain":       "github.com",
			"type":         "Active",
			"user":         "testuser",
		},
		{
			"browser_type": "chrome",
			"profile":      filepath.Base(tempDir),
			"extension_id": "aeblfdkhhhdcdjpifhhbdiojplfjncoa\x00",
			"domain":       "b5x-sentry.1passwordservices.com",
			"type":         "Broken",
			"user":         "testuser",
		},
	}

	assert.ElementsMatch(t, expected, results)
}
