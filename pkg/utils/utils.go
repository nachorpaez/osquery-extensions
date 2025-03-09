package utils

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type findFile struct {
	username string
}

type FindFileOpt func(*findFile)

func WithUsername(username string) FindFileOpt {
	return func(ff *findFile) {
		ff.username = username
	}
}

var HomeDirLocations = map[string][]string{
	"windows": {"/Users"}, // windows10 uses /Users
	"darwin":  {"/Users"},
	"linux":   {"/home/"},
}
var homeDirDefaultLocation = []string{"/home"}

type UserFileInfo struct {
	User string
	Path string
}

// findFileInUserDirs looks for the existence of a specified path as a
// subdirectory of users' home directories. It does this by searching
// likely paths
func FindFileInUserDirs(pattern string, opts ...FindFileOpt) ([]UserFileInfo, error) {
	ff := &findFile{}

	for _, opt := range opts {
		opt(ff)
	}

	homedirRoots, ok := HomeDirLocations[runtime.GOOS]
	if !ok {
		homedirRoots = homeDirDefaultLocation
		log.Printf("Platform not found using default home_dir_root: %s", homedirRoots)
	}

	foundPaths := []UserFileInfo{}

	if ff.username == "" {
		for _, possibleHome := range homedirRoots {

			userDirs, err := os.ReadDir(possibleHome)
			if err != nil {
				// This possibleHome doesn't exist. Move on
				continue
			}

			// For each user's dir, in this possibleHome, check
			for _, ud := range userDirs {
				userPathPattern := filepath.Join(possibleHome, ud.Name(), pattern)
				fullPaths, err := filepath.Glob(userPathPattern)
				if err != nil {
					// skipping ErrBadPattern
					log.Printf("Bad file pattern %s", userPathPattern)
					continue
				}
				// If the found path is a file, add it to the list
				for _, fullPath := range fullPaths {
					if stat, err := os.Stat(fullPath); err == nil && stat.Mode().IsRegular() {
						foundPaths = append(foundPaths, UserFileInfo{
							User: ud.Name(),
							Path: fullPath,
						})
					}
				}
			}
		}
		return foundPaths, nil
	}

	// We have a username. Future normal path here
	for _, possibleHome := range homedirRoots {
		userPathPattern := filepath.Join(possibleHome, ff.username, pattern)
		fullPaths, err := filepath.Glob(userPathPattern)
		if err != nil {
			continue
		}
		for _, fullPath := range fullPaths {
			if stat, err := os.Stat(fullPath); err == nil && stat.Mode().IsRegular() {
				foundPaths = append(foundPaths, UserFileInfo{
					User: ff.username,
					Path: fullPath,
				})
			}
		}
	}
	return foundPaths, nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Btoi(value bool) int {
	if value {
		return 1
	}
	return 0
}
