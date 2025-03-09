package chrome_extensions_dns

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/nachorpaez/osquery-extensions/pkg/utils"
	"github.com/osquery/osquery-go/plugin/table"
	"github.com/pkg/errors"
)

// NetworkState represents the Chrome network state file structure.
type NetworkState struct {
	Net struct {
		HTTPServerProperties struct {
			Servers                   []Server                   `json:"servers"`
			BrokenAlternativeServices []BrokenAlternativeService `json:"broken_alternative_services"`
		} `json:"http_server_properties"`
	} `json:"net"`
}

type Server struct {
	Server        string        `json:"server"`
	Anonymization []interface{} `json:"anonymization"`
}

type BrokenAlternativeService struct {
	Host          string        `json:"host"`
	Anonymization []interface{} `json:"anonymization"`
	BrokenUntil   string        `json:"broken_until"`
}

// ChromeExtensionDNSColumns defines the columns for the osquery table.
func ChromeExtensionsDNSColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("browser_type"),
		table.TextColumn("profile"),
		table.TextColumn("extension_id"),
		table.TextColumn("domain"),
		table.TextColumn("type"),
		table.TextColumn("user"),
	}
}

// decodeAnonymization decodes a base64 string containing a Chrome extension ID.
// It returns the extension ID if found; otherwise, an empty string.
func decodeAnonymization(anonymizationValue interface{}) string {
	if str, ok := anonymizationValue.(string); ok {
		decoded, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			return ""
		}
		decodedStr := string(decoded)
		if strings.Contains(decodedStr, "chrome-extension://") {
			parts := strings.Split(decodedStr, "chrome-extension://")
			return parts[len(parts)-1]
		}
	}
	return ""
}

// analyzeNetworkState processes a Chrome profile's network state file to extract
// information about active and broken connections. It returns a slice of maps
// containing fields such as browser_type, profile, extension_id, domain, and more.
func analyzeNetworkState(ctx context.Context, profileInfo utils.ChromeProfilePath) ([]map[string]string, error) {
	var results []map[string]string
	profileName := filepath.Base(profileInfo.Value)
	stateFile := filepath.Join(profileInfo.Value, "Network Persistent State")

	// Try alternate path if the first one doesn't exist
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		stateFile = filepath.Join(profileInfo.Value, "Network", "Network Persistent State")
		if _, err := os.Stat(stateFile); os.IsNotExist(err) {
			return nil, fmt.Errorf("network state file not found")
		}
	}

	data, err := os.ReadFile(stateFile)
	if err != nil {
		log.Printf("Error reading state file: %s", err)
		return nil, errors.Wrap(err, "reading state file")
	}

	var netState NetworkState
	if err := json.Unmarshal(data, &netState); err != nil {
		log.Printf("Error parsing JSON: %s", err)
		return nil, errors.Wrap(err, "parsing JSON")
	}

	// Process active connections.
	for _, server := range netState.Net.HTTPServerProperties.Servers {
		if len(server.Anonymization) > 0 {
			if extID := decodeAnonymization(server.Anonymization[0]); extID != "" {
				url, err := url.Parse(server.Server)
				if err != nil {
					log.Printf("Error parsing URL: %s", err)
					continue
				}
				results = append(results, map[string]string{
					"browser_type": utils.GetChromeBrowserName(profileInfo.Type),
					"profile":      profileName,
					"extension_id": extID,
					"domain":       url.Hostname(),
					"type":         "Active",
					"user":         profileInfo.UserName,
				})
			}
		}
	}

	// Process broken connections.
	for _, broken := range netState.Net.HTTPServerProperties.BrokenAlternativeServices {
		if len(broken.Anonymization) > 0 {
			if extID := decodeAnonymization(broken.Anonymization[0]); extID != "" {
				results = append(results, map[string]string{
					"browser_type": utils.GetChromeBrowserName(profileInfo.Type),
					"profile":      profileName,
					"extension_id": extID,
					"domain":       broken.Host,
					"type":         "Broken",
					"user":         profileInfo.UserName,
				})
			}
		}
	}

	return results, nil
}

// Per docs generator function has to return an array of map of strings
func ChromeExtensionsDNSGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var results []map[string]string

	profileList, err := utils.GetChromeProfilePathList()
	if err != nil {
		log.Printf("Error retrieving Chrome profiles list: %s", err)
		return nil, errors.Wrap(err, "retrieving Chrome profile list")
	}

	for _, profile := range profileList {
		res, _ := analyzeNetworkState(ctx, profile)
		results = append(results, res...)
	}
	return results, nil
}
