package chrome_preferences

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/nachorpaez/osquery-extensions/pkg/utils"
	"github.com/osquery/osquery-go/plugin/table"
	"github.com/pkg/errors"
)

type Settings struct {
	Expiration   string `json:"expiration,omitempty"`
	LastModified string `json:"last_modified,omitempty"`
	Model        int    `json:"model,omitempty"`
	Setting      int    `json:"setting"`
}

type Profile struct {
	ContentSettings struct {
		Exceptions struct {
			Notifications     map[string]Settings `json:"notifications"`
			Geolocation       map[string]Settings `json:"geolocation"`
			MediaStreamCamera map[string]Settings `json:"media_stream_camera"`
			MediaStreamMic    map[string]Settings `json:"media_stream_mic"`
			Popups            map[string]Settings `json:"popups"`
		} `json:"exceptions"`
	} `json:"content_settings"`
}

type Data struct {
	Profile Profile `json:"profile"`
}

func GoogleChromePreferencesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("category"),
		table.TextColumn("url"),
		table.TextColumn("expiration"),
		table.BigIntColumn("last_modified"),
		table.IntegerColumn("model"),
		table.IntegerColumn("setting"),
		table.TextColumn("profile_path"),
		table.TextColumn("user"),
		table.TextColumn("browser_type"),
	}
}

func parsePreferences(ctx context.Context, chromeProfile utils.ChromeProfilePath) ([]map[string]string, error) {
	var results []map[string]string
	preferenceFile := filepath.Join(chromeProfile.Value, "Preferences")
	fileContent, err := os.ReadFile(preferenceFile)

	if err != nil {
		return nil, errors.Wrap(err, "reading preferences file")
	}
	var data Data
	if err := json.Unmarshal(fileContent, &data); err != nil {
		return nil, errors.Wrap(err, "unmarshalling preferences file")
	}

	for url, preference := range data.Profile.ContentSettings.Exceptions.Geolocation {
		results = append(results, map[string]string{
			"category":      "geolocation",
			"url":           url,
			"expiration":    preference.Expiration,
			"last_modified": preference.LastModified,
			"model":         strconv.Itoa(preference.Model),
			"setting":       strconv.Itoa(preference.Setting),
			"profile_path":  chromeProfile.Value,
			"user":          chromeProfile.UserName,
			"browser_type":  utils.GetChromeBrowserName(chromeProfile.Type),
		})
	}

	for url, preference := range data.Profile.ContentSettings.Exceptions.Notifications {
		results = append(results, map[string]string{
			"category":      "notifications",
			"url":           url,
			"expiration":    preference.Expiration,
			"last_modified": preference.LastModified,
			"model":         strconv.Itoa(preference.Model),
			"setting":       strconv.Itoa(preference.Setting),
			"profile_path":  chromeProfile.Value,
			"user":          chromeProfile.UserName,
			"browser_type":  utils.GetChromeBrowserName(chromeProfile.Type),
		})
	}

	for url, preference := range data.Profile.ContentSettings.Exceptions.MediaStreamCamera {
		results = append(results, map[string]string{
			"category":      "media_stream_camera",
			"url":           url,
			"expiration":    preference.Expiration,
			"last_modified": preference.LastModified,
			"model":         strconv.Itoa(preference.Model),
			"setting":       strconv.Itoa(preference.Setting),
			"profile_path":  chromeProfile.Value,
			"user":          chromeProfile.UserName,
			"browser_type":  utils.GetChromeBrowserName(chromeProfile.Type),
		})
	}

	for url, preference := range data.Profile.ContentSettings.Exceptions.MediaStreamMic {
		results = append(results, map[string]string{
			"category":      "media_stream_mic",
			"url":           url,
			"expiration":    preference.Expiration,
			"last_modified": preference.LastModified,
			"model":         strconv.Itoa(preference.Model),
			"setting":       strconv.Itoa(preference.Setting),
			"profile_path":  chromeProfile.Value,
			"user":          chromeProfile.UserName,
			"browser_type":  utils.GetChromeBrowserName(chromeProfile.Type),
		})
	}
	for url, preference := range data.Profile.ContentSettings.Exceptions.Popups {
		results = append(results, map[string]string{
			"category":      "popups",
			"url":           url,
			"expiration":    preference.Expiration,
			"last_modified": preference.LastModified,
			"model":         strconv.Itoa(preference.Model),
			"setting":       strconv.Itoa(preference.Setting),
			"profile_path":  chromeProfile.Value,
			"user":          chromeProfile.UserName,
			"browser_type":  utils.GetChromeBrowserName(chromeProfile.Type),
		})
	}

	return results, nil
}

// Per docs generator function has to return an array of map of strings
func GoogleChromePreferencesGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var results []map[string]string

	profileList, err := utils.GetChromeProfilePathList()
	if err != nil {
		log.Printf("Error retrieving Chrome profiles list: %s", err)
		return nil, errors.Wrap(err, "retrieving Chrome profile list")
	}

	for _, profile := range profileList {
		res, _ := parsePreferences(ctx, profile)
		results = append(results, res...)
	}
	return results, nil
}
