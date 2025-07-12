package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TimeFormat represents a time parsing format
type TimeFormat struct {
	Name    string
	Regex   *regexp.Regexp
	Handler func([]string, string) (uint, error)
}

// Public format variables for individual access
var (
	// Format12Hour represents 12-hour format with am/pm: 6 am, 6:30 pm, 12:45 am, etc.
	Format12Hour = TimeFormat{
		Name:    "12-hour with am/pm",
		Regex:   regexp.MustCompile(`\b(\d{1,2})(?::(\d{2}))?\s*(am|pm)\b`),
		Handler: parse12HourFormat,
	}

	// Format24Hour represents 24-hour format: 18:00, 18:30, 09:15, etc.
	Format24Hour = TimeFormat{
		Name:    "24-hour format",
		Regex:   regexp.MustCompile(`\b(\d{1,2}):(\d{2})\b`),
		Handler: parse24HourFormat,
	}

	// FormatMilitary represents military time format: 1542, 0900, 2359, etc.
	FormatMilitary = TimeFormat{
		Name:    "military time",
		Regex:   regexp.MustCompile(`\b(\d{4})\b`),
		Handler: parseMilitaryTime,
	}

	// FormatSimpleHour represents simple hour format: 6, 18, etc.
	FormatSimpleHour = TimeFormat{
		Name:    "simple hour",
		Regex:   regexp.MustCompile(`\b(\d{1,2})\s*(?:o'clock|oclock|oc|o'c)\b`),
		Handler: parseSimpleHour,
	}
)

// getDefaultFormats returns the default time formats
func getDefaultFormats() []TimeFormat {
	return []TimeFormat{
		Format12Hour,
		Format24Hour,
		FormatMilitary,
		FormatSimpleHour,
	}
}

// parse12HourFormat parses 12-hour format with am/pm: 6 am, 6:30 pm, 12:45 am, etc.
func parse12HourFormat(matches []string, _ string) (uint, error) {
	hour, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid hour: %s", matches[1])
	}
	minute := 0
	if matches[2] != "" {
		minute, err = strconv.Atoi(matches[2])
		if err != nil {
			return 0, fmt.Errorf("invalid minute: %s", matches[2])
		}
	}
	if hour < 1 || hour > 12 || minute < 0 || minute > 59 {
		return 0, fmt.Errorf("invalid time: %s:%s %s", matches[1], matches[2], matches[3])
	}
	if matches[3] == "am" {
		if hour == 12 {
			hour = 0
		}
	} else if matches[3] == "pm" {
		if hour != 12 {
			hour += 12
		}
	}
	return uint(hour*3600 + minute*60), nil
}

// parse24HourFormat parses 24-hour format: 18:00, 18:30, 09:15, etc.
func parse24HourFormat(matches []string, msg string) (uint, error) {
	hour, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid hour: %s", matches[1])
	}
	minute, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, fmt.Errorf("invalid minute: %s", matches[2])
	}
	// Check if 'am' or 'pm' follows the match in the message
	pattern := matches[0]
	idx := strings.Index(msg, pattern)
	if idx != -1 {
		endIdx := idx + len(pattern)
		if endIdx+2 <= len(msg) {
			suffix := msg[endIdx:]
			suffix = strings.TrimSpace(suffix)
			if strings.HasPrefix(suffix, "am") || strings.HasPrefix(suffix, "pm") {
				return 0, fmt.Errorf("24-hour match followed by am/pm, skip")
			}
		}
	}
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return 0, fmt.Errorf("invalid time: %s:%s", matches[1], matches[2])
	}
	return uint(hour*3600 + minute*60), nil
}

// parseMilitaryTime parses military time format: 1542, 0900, 2359, etc.
func parseMilitaryTime(matches []string, _ string) (uint, error) {
	timeStr := matches[1]
	if len(timeStr) != 4 {
		return 0, fmt.Errorf("invalid military time format: %s", timeStr)
	}
	hour, err := strconv.Atoi(timeStr[:2])
	if err != nil {
		return 0, fmt.Errorf("invalid hour: %s", timeStr[:2])
	}
	minute, err := strconv.Atoi(timeStr[2:])
	if err != nil {
		return 0, fmt.Errorf("invalid minute: %s", timeStr[2:])
	}
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return 0, fmt.Errorf("invalid military time: %s", timeStr)
	}
	return uint(hour*3600 + minute*60), nil
}

// parseSimpleHour parses simple hour format: 6, 18, etc.
func parseSimpleHour(matches []string, _ string) (uint, error) {
	hour, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid hour: %s", matches[1])
	}
	if hour < 0 || hour > 23 {
		return 0, fmt.Errorf("invalid hour: %d", hour)
	}
	return uint(hour * 3600), nil
}
