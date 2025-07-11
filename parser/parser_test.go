package parser

import (
	"regexp"
	"testing"
)

func TestTimeParser_ParseTimeFromMessage(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected uint
		hasError bool
	}{
		// Valid 24-hour format tests
		{"24-hour format", "18:00", 64800, false},
		{"24-hour format with leading zero", "09:15", 33300, false},
		{"24-hour format midnight", "00:00", 0, false},
		{"24-hour format end of day", "23:59", 86340, false},

		// Valid 12-hour format tests
		{"12-hour am", "6 am", 21600, false},
		{"12-hour pm", "6 pm", 64800, false},
		{"12-hour with minutes am", "6:30 am", 23400, false},
		{"12-hour with minutes pm", "6:30 pm", 66600, false},
		{"12-hour 12 am", "12 am", 0, false},
		{"12-hour 12 pm", "12 pm", 43200, false},
		{"12-hour 12:30 am", "12:30 am", 1800, false},
		{"12-hour 12:30 pm", "12:30 pm", 45000, false},
		{"12-hour 12:30pm", "12:30pm", 45000, false},
		{"12-hour 12:30am", "12:30am", 1800, false},

		// Valid military time tests
		{"military time", "1542", 56520, false},
		{"military time with leading zero", "0900", 32400, false},
		{"military time midnight", "0000", 0, false},
		{"military time end of day", "2359", 86340, false},

		// Valid simple hour tests
		{"simple hour", "9 o'clock", 32400, false},
		{"simple hour oclock", "9 oclock", 32400, false},
		{"simple hour oc", "9 oc", 32400, false},
		{"simple hour o'c", "9 o'c", 32400, false},

		// Tests with context
		{"time in sentence", "Let's meet at 6:30 pm tomorrow", 66600, false},
		{"time with military context", "The meeting is at 1542 hours", 56520, false},
		{"time with am context", "I'll be there at 6 am sharp", 21600, false},
		{"time with 24-hour context", "Flight departs at 18:00", 64800, false},

		// Error cases
		{"invalid hour 24-hour", "24:00", 0, true},
		{"invalid minute 24-hour", "12:60", 0, true},
		{"invalid hour 12-hour", "13 am", 0, true},
		{"invalid minute 12-hour", "12:60 pm", 0, true},
		{"invalid military time", "2500", 0, true},
		{"invalid military time format", "123", 0, true},
		{"no time in message", "Hello world", 0, true},
		{"empty message", "", 0, true},
	}

	tp := NewTimeParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tp.ParseTimeFromMessage(tt.message)

			if tt.hasError {
				if err == nil {
					t.Errorf("ParseTimeFromMessage() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ParseTimeFromMessage() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("ParseTimeFromMessage() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestTimeParser_ParseTimeFromMessage_EdgeCases(t *testing.T) {
	tp := NewTimeParser()

	// Test case sensitivity
	message := "6 AM"
	expected := uint(21600)

	result, err := tp.ParseTimeFromMessage(message)
	if err != nil {
		t.Errorf("ParseTimeFromMessage() unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("ParseTimeFromMessage() = %v, want %v", result, expected)
	}

	// Test multiple time formats in same message (should return first match)
	message = "Let's meet at 6:30 pm and also at 18:00"
	result, err = tp.ParseTimeFromMessage(message)
	if err != nil {
		t.Errorf("ParseTimeFromMessage() unexpected error: %v", err)
	}
	// Should return the first match (6:30 pm = 66600 seconds)
	if result != 66600 {
		t.Errorf("ParseTimeFromMessage() = %v, want %v", result, 66600)
	}
}

func TestNewTimeParserWithFormats(t *testing.T) {
	// Test custom parser with only 24-hour format
	customFormats := []TimeFormat{
		{
			Name:    "24-hour format",
			Regex:   regexp.MustCompile(`\b(\d{1,2}):(\d{2})\b`),
			Handler: parse24HourFormat,
		},
	}

	tp := NewTimeParserWithFormats(customFormats)

	// Should parse 24-hour format
	seconds, err := tp.ParseTimeFromMessage("18:00")
	if err != nil {
		t.Errorf("Expected to parse 24-hour format, got error: %v", err)
	}
	if seconds != 64800 {
		t.Errorf("Expected 64800 seconds, got %d", seconds)
	}

	// Should not parse 12-hour format
	_, err = tp.ParseTimeFromMessage("6 pm")
	if err == nil {
		t.Errorf("Expected error for 12-hour format with custom parser")
	}
}

func TestParse12HourFormat(t *testing.T) {
	tests := []struct {
		name     string
		matches  []string
		expected uint
		hasError bool
	}{
		{"6 am", []string{"6 am", "6", "", "am"}, 21600, false},
		{"6:30 pm", []string{"6:30 pm", "6", "30", "pm"}, 66600, false},
		{"12 am", []string{"12 am", "12", "", "am"}, 0, false},
		{"12 pm", []string{"12 pm", "12", "", "pm"}, 43200, false},
		{"12:30 am", []string{"12:30 am", "12", "30", "am"}, 1800, false},
		{"12:30 pm", []string{"12:30 pm", "12", "30", "pm"}, 45000, false},
		{"invalid hour", []string{"13 am", "13", "", "am"}, 0, true},
		{"invalid minute", []string{"12:60 pm", "12", "60", "pm"}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parse12HourFormat(tt.matches, "")

			if tt.hasError {
				if err == nil {
					t.Errorf("parse12HourFormat() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("parse12HourFormat() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("parse12HourFormat() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestParse24HourFormat(t *testing.T) {
	tests := []struct {
		name     string
		matches  []string
		message  string
		expected uint
		hasError bool
	}{
		{"18:00", []string{"18:00", "18", "00"}, "18:00", 64800, false},
		{"09:15", []string{"09:15", "09", "15"}, "09:15", 33300, false},
		{"00:00", []string{"00:00", "00", "00"}, "00:00", 0, false},
		{"23:59", []string{"23:59", "23", "59"}, "23:59", 86340, false},
		{"invalid hour", []string{"24:00", "24", "00"}, "24:00", 0, true},
		{"invalid minute", []string{"12:60", "12", "60"}, "12:60", 0, true},
		{"with am/pm", []string{"18:00", "18", "00"}, "18:00 pm", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parse24HourFormat(tt.matches, tt.message)

			if tt.hasError {
				if err == nil {
					t.Errorf("parse24HourFormat() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("parse24HourFormat() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("parse24HourFormat() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestParseMilitaryTime(t *testing.T) {
	tests := []struct {
		name     string
		matches  []string
		expected uint
		hasError bool
	}{
		{"1542", []string{"1542", "1542"}, 56520, false},
		{"0900", []string{"0900", "0900"}, 32400, false},
		{"0000", []string{"0000", "0000"}, 0, false},
		{"2359", []string{"2359", "2359"}, 86340, false},
		{"invalid time", []string{"2500", "2500"}, 0, true},
		{"invalid format", []string{"123", "123"}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseMilitaryTime(tt.matches, "")

			if tt.hasError {
				if err == nil {
					t.Errorf("parseMilitaryTime() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("parseMilitaryTime() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("parseMilitaryTime() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestParseSimpleHour(t *testing.T) {
	tests := []struct {
		name     string
		matches  []string
		expected uint
		hasError bool
	}{
		{"9", []string{"9", "9"}, 32400, false},
		{"18", []string{"18", "18"}, 64800, false},
		{"0", []string{"0", "0"}, 0, false},
		{"23", []string{"23", "23"}, 82800, false},
		{"invalid hour", []string{"24", "24"}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseSimpleHour(tt.matches, "")

			if tt.hasError {
				if err == nil {
					t.Errorf("parseSimpleHour() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("parseSimpleHour() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("parseSimpleHour() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}
