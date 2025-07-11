package timezones

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

var TimezoneLocations []string

func init() {
	tzs, err := getValidTimezones("/usr/share/zoneinfo")
	if err != nil {
		panic(err)
	}
	TimezoneLocations = tzs
}

// Filters out known system files and paths
func isValidZonePath(path, base string) bool {
	rel := strings.TrimPrefix(path, base+"/")

	// Skip top-level files and known invalids
	switch rel {
	case "localtime", "posixrules", "Factory", "leap-seconds.list", "tzdata.zi":
		return false
	}

	// Skip some system directories
	if strings.HasPrefix(rel, "posix/") ||
		strings.HasPrefix(rel, "right/") ||
		strings.HasPrefix(rel, "SystemV") ||
		strings.HasPrefix(rel, "Etc/") ||
		strings.Contains(rel, "/.") ||
		strings.HasSuffix(rel, ".tab") ||
		strings.Contains(rel, "leap") {
		return false
	}

	return true
}

func getValidTimezones(base string) ([]string, error) {
	var zones []string

	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if isValidZonePath(path, base) {
			zone := strings.TrimPrefix(path, base+"/")

			// Final validation: test with time.LoadLocation
			if _, err := time.LoadLocation(zone); err == nil {
				zones = append(zones, zone)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return zones, nil
}
