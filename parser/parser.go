package parser

import (
	"fmt"
	"strings"
)

// MatchResult represents a parsed time match
type MatchResult struct {
	Start      int
	End        int
	PatternIdx int
	Matches    []string
}

// ParseResult represents the result of parsing a time from a message
type ParseResult struct {
	Seconds uint
	Error   error
}

// TimeParser handles parsing of time formats from messages
type TimeParser struct {
	formats []TimeFormat
}

// NewTimeParser creates a new TimeParser with default formats
func NewTimeParser() *TimeParser {
	return &TimeParser{
		formats: getDefaultFormats(),
	}
}

// NewTimeParserWithFormats creates a new TimeParser with custom formats
func NewTimeParserWithFormats(formats []TimeFormat) *TimeParser {
	return &TimeParser{
		formats: formats,
	}
}

// ParseTimeFromMessage parses number of seconds since midnight from a message
func (tp *TimeParser) ParseTimeFromMessage(message string) (uint, error) {
	lowerMessage := strings.ToLower(message)

	var allMatches []MatchResult
	for i, format := range tp.formats {
		locs := format.Regex.FindAllStringSubmatchIndex(lowerMessage, -1)
		for _, loc := range locs {
			start, end := loc[0], loc[1]
			matches := make([]string, len(loc)/2)
			for j := 0; j < len(loc)/2; j++ {
				subStart, subEnd := loc[2*j], loc[2*j+1]
				if subStart >= 0 && subEnd >= 0 {
					matches[j] = lowerMessage[subStart:subEnd]
				} else {
					matches[j] = ""
				}
			}
			if matches[0] != "" {
				allMatches = append(allMatches, MatchResult{start, end, i, matches})
			}
		}
	}

	if len(allMatches) > 0 {
		// Find the match with the lowest start index, then longest match, then lowest patternIdx
		best := allMatches[0]
		for _, m := range allMatches[1:] {
			if m.Start < best.Start ||
				(m.Start == best.Start && (m.End-m.Start) > (best.End-best.Start)) ||
				(m.Start == best.Start && (m.End-m.Start) == (best.End-best.Start) && m.PatternIdx < best.PatternIdx) {
				best = m
			}
		}
		seconds, err := tp.formats[best.PatternIdx].Handler(best.Matches, lowerMessage)
		if err == nil {
			return seconds, nil
		}
	}

	return 0, fmt.Errorf("no valid time format found in message: %s", message)
}
