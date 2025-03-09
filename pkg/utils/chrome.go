package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// Preferences file included in each profile
const ProfilePreferencesFile = "Preferences"

// Alternative 'Secure Preferences' file included in each profile
const SecureProfilePreferencesFile = "Secure Preferences"

// Possible configuration file names
var possibleConfigFileNames = []string{
	ProfilePreferencesFile,
	SecureProfilePreferencesFile,
}

// ChromeProfilePath represents a Chrome profile path with user and browser type information
type ChromeProfilePath struct {
	UserName string
	Type     ChromeBrowserType
	Value    string
}

// ChromeBrowserType represents different types of Chrome-based browsers
type ChromeBrowserType int

const (
	GoogleChrome ChromeBrowserType = iota
	GoogleChromeBeta
	GoogleChromeDev
	GoogleChromeCanary
	Brave
	Chromium
	Yandex
	Edge
	EdgeBeta
	Opera
	Vivaldi
	Arc
)

// WindowsPathList maps browser types to their Windows installation paths
var WindowsPathList = map[ChromeBrowserType]string{
	GoogleChrome:       "AppData\\Local\\Google\\Chrome\\User Data",
	GoogleChromeBeta:   "AppData\\Local\\Google\\Chrome Beta\\User Data",
	GoogleChromeDev:    "AppData\\Local\\Google\\Chrome Dev\\User Data",
	GoogleChromeCanary: "AppData\\Local\\Google\\Chrome SxS\\User Data",
	Brave:              "AppData\\Roaming\\brave",
	Chromium:           "AppData\\Local\\Chromium",
	Yandex:             "AppData\\Local\\Yandex\\YandexBrowser\\User Data",
	Edge:               "AppData\\Local\\Microsoft\\Edge\\User Data",
	EdgeBeta:           "AppData\\Local\\Microsoft\\Edge Beta\\User Data",
	Opera:              "AppData\\Roaming\\Opera Software\\Opera Stable",
	Vivaldi:            "AppData\\Local\\Vivaldi\\User Data",
}

// MacOSPathList maps browser types to their macOS installation paths
var MacOSPathList = map[ChromeBrowserType]string{
	GoogleChrome:       "Library/Application Support/Google/Chrome",
	GoogleChromeBeta:   "Library/Application Support/Google/Chrome Beta",
	GoogleChromeDev:    "Library/Application Support/Google/Chrome Dev",
	GoogleChromeCanary: "Library/Application Support/Google/Chrome Canary",
	Brave:              "Library/Application Support/BraveSoftware/Brave-Browser",
	Chromium:           "Library/Application Support/Chromium",
	Yandex:             "Library/Application Support/Yandex/YandexBrowser",
	Edge:               "Library/Application Support/Microsoft Edge",
	EdgeBeta:           "Library/Application Support/Microsoft Edge Beta",
	Opera:              "Library/Application Support/com.operasoftware.Opera",
	Vivaldi:            "Library/Application Support/Vivaldi",
	Arc:                "Library/Application Support/Arc/User Data",
}

// LinuxPathList maps browser types to their Linux installation paths
var LinuxPathList = map[ChromeBrowserType]string{
	GoogleChrome:     ".config/google-chrome",
	GoogleChromeBeta: ".config/google-chrome-beta",
	GoogleChromeDev:  ".config/google-chrome-unstable",
	Brave:            ".config/BraveSoftware/Brave-Browser",
	Chromium:         ".config/chromium",
	Yandex:           ".config/yandex-browser-beta",
	Opera:            ".config/opera",
	Vivaldi:          ".config/vivaldi",
}

// ChromeBrowserTypeToString maps browser types to their string representations
var ChromeBrowserTypeToString = map[ChromeBrowserType]string{
	GoogleChrome:       "chrome",
	GoogleChromeBeta:   "chrome_beta",
	GoogleChromeDev:    "chrome_dev",
	GoogleChromeCanary: "chrome_canary",
	Brave:              "brave",
	Chromium:           "chromium",
	Yandex:             "yandex",
	Opera:              "opera",
	Edge:               "edge",
	EdgeBeta:           "edge_beta",
	Vivaldi:            "vivaldi",
	Arc:                "arc",
}

// GetChromeBrowserName returns the string representation of a ChromeBrowserType
func GetChromeBrowserName(t ChromeBrowserType) string {
	name, ok := ChromeBrowserTypeToString[t]
	if !ok {
		return "" // Return an empty string if not found
	}
	return name
}

// GetChromePathSuffixMap returns the appropriate path list based on the operating system
func GetChromePathSuffixMap() map[ChromeBrowserType]string {
	switch runtime.GOOS {
	case "windows":
		return WindowsPathList
	case "darwin":
		return MacOSPathList
	default:
		return LinuxPathList
	}
}

// isValidChromeProfile returns true if the given path contains either the Preferences
// or the Secure Preferences file
func isValidChromeProfile(path string) bool {
	for _, configFileName := range possibleConfigFileNames {
		preferencesFilePath := filepath.Join(path, configFileName)

		if _, err := os.Stat(preferencesFilePath); err == nil {
			// If os.Stat returns no error, the file exists and is accessible
			return true
		}
	}

	return false
}

// listDirectoriesInDirectory returns a slice of subdirectories within the given path.
func listDirectoriesInDirectory(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, filepath.Join(path, e.Name()))
		}
	}
	return dirs, nil
}

// GetChromeProfilePathList attempts to discover valid Chrome profiles
// based on user information and known Chrome installation paths.
func GetChromeProfilePathList() ([]ChromeProfilePath, error) {
	// Get User directories
	userInfoList := map[string]string{}
	for _, possibleHome := range HomeDirLocations[runtime.GOOS] {
		userDirs, err := os.ReadDir(possibleHome)
		if err != nil {
			return nil, err
		}

		for _, userDir := range userDirs {
			if userDir.IsDir() {
				userInfoList[userDir.Name()] = filepath.Join(possibleHome, userDir.Name())
			}
		}
	}

	var output []ChromeProfilePath

	for userName, userPath := range userInfoList {
		// Prepare a ChromeProfilePath struct, fill the UID for each user.
		chromeProfile := ChromeProfilePath{
			UserName: userName,
		}

		pathSuffixMap := GetChromePathSuffixMap()

		for browserType, pathSuffix := range pathSuffixMap {
			// Set the browser type in the profile.
			chromeProfile.Type = browserType

			// Construct the path to the user's Chrome directory.
			path := filepath.Join(userPath, pathSuffix)

			// Attempt to resolve symlinks.
			absoluteChromePath, err := filepath.EvalSymlinks(path)
			if err != nil {
				// If an error occurs, just use the original path.
				absoluteChromePath = path
			}

			// Check if this directory itself is a valid Chrome profile.
			if isValidChromeProfile(absoluteChromePath) {
				chromeProfile.Value = absoluteChromePath
				output = append(output, chromeProfile)
				continue
			}

			// Otherwise, attempt to find subdirectories that may be valid profiles.
			subfolders, err := listDirectoriesInDirectory(absoluteChromePath)
			if err != nil {
				// If there's an error listing directories, skip this folder.
				continue
			}

			// Check each subfolder for a valid Chrome profile.
			for _, subfolder := range subfolders {
				absSubfolder, err := filepath.EvalSymlinks(subfolder)
				if err != nil {
					absSubfolder = subfolder
				}

				if isValidChromeProfile(absSubfolder) {
					chromeProfile.Value = absSubfolder
					output = append(output, chromeProfile)
					continue
				}
			}
		}
	}

	return output, nil
}
