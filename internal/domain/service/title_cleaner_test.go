package service

import (
	"testing"

	"kbfood/internal/domain/entity"
)

func TestTitleCleaner_CleanTitleForID(t *testing.T) {
	tc := NewTitleCleaner()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal Chinese title",
			input:    "巧克力草莓蛋糕",
			expected: "巧克力草莓蛋糕",
		},
		{
			name:     "Title with spaces and symbols",
			input:    "巧克力草莓蛋糕 (6寸)",
			expected: "巧克力草莓蛋糕6寸",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Mixed alphanumeric and Chinese",
			input:    "ABC123巧克力草莓",
			expected: "ABC123巧克力草莓",
		},
		{
			name:     "Only special characters",
			input:    "@#$%^&*()",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tc.CleanTitleForID(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestTitleCleaner_CalculateSimilarity(t *testing.T) {
	tc := NewTitleCleaner()

	tests := []struct {
		name     string
		s1       string
		s2       string
		expected float64
	}{
		{
			name:     "Identical strings",
			s1:       "巧克力草莓蛋糕",
			s2:       "巧克力草莓蛋糕",
			expected: 1.0,
		},
		{
			name:     "Similar strings",
			s1:       "巧克力草莓蛋糕",
			s2:       "巧克力草莓蛋糕6寸",
			expected: 0.8, // Approximate
		},
		{
			name:     "Different strings",
			s1:       "巧克力草莓蛋糕",
			s2:       "提拉米苏",
			expected: 0.0,
		},
		{
			name:     "Both empty",
			s1:       "",
			s2:       "",
			expected: 0.0,
		},
		{
			name:     "One empty",
			s1:       "巧克力",
			s2:       "",
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tc.CalculateSimilarity(tt.s1, tt.s2)
			if tt.name == "Identical strings" {
				if result != tt.expected {
					t.Errorf("expected %f, got %f", tt.expected, result)
				}
			} else if tt.expected == 0.0 {
				if result != 0.0 {
					t.Errorf("expected 0.0, got %f", result)
				}
			} else if result < tt.expected-0.2 || result > tt.expected+0.2 {
				t.Errorf("expected ~%f, got %f", tt.expected, result)
			}
		})
	}
}

func TestTitleCleaner_IsHighSimilarity(t *testing.T) {
	tc := NewTitleCleaner()

	if !tc.IsHighSimilarity("巧克力草莓蛋糕", "巧克力草莓蛋糕") {
		t.Error("identical strings should be high similarity")
	}

	if !tc.IsHighSimilarity("巧克力草莓蛋糕6寸", "巧克力草莓蛋糕") {
		t.Error("similar strings should be high similarity")
	}

	if tc.IsHighSimilarity("巧克力草莓蛋糕", "提拉米苏") {
		t.Error("different strings should not be high similarity")
	}

	// Edge cases
	if tc.IsHighSimilarity("", "巧克力") {
		t.Error("empty string should not be high similarity")
	}

	if tc.IsHighSimilarity("巧克力", "") {
		t.Error("empty string should not be high similarity")
	}
}

func TestTitleCleaner_IsMidSimilarity(t *testing.T) {
	tc := NewTitleCleaner()

	if tc.IsMidSimilarity("巧克力草莓蛋糕", "提拉米苏") {
		t.Error("very different strings should not be mid similarity")
	}

	// Edge cases
	if tc.IsMidSimilarity("", "巧克力") {
		t.Error("empty string should not be mid similarity")
	}
}

func TestTitleCleaner_ElectStandardTitle(t *testing.T) {
	tc := NewTitleCleaner()

	votes := map[string]int{
		"巧克力草莓蛋糕":      5,
		"巧克力草莓蛋糕6寸":    3,
		"巧克力草莓蛋糕 (6寸)": 2,
	}

	winner := tc.ElectStandardTitle(votes)
	if winner != "巧克力草莓蛋糕" {
		t.Errorf("expected '巧克力草莓蛋糕', got %s", winner)
	}

	// Empty votes
	emptyWinner := tc.ElectStandardTitle(map[string]int{})
	if emptyWinner != "" {
		t.Errorf("expected empty string, got %s", emptyWinner)
	}
}

func TestTitleCleaner_GenerateID(t *testing.T) {
	tc := NewTitleCleaner()

	id := tc.GenerateID("DT", "巧克力草莓蛋糕")
	if id == "" {
		t.Error("expected non-empty ID")
	}

	// Same title should generate same ID
	id2 := tc.GenerateID("DT", "巧克力草莓蛋糕")
	if id != id2 {
		t.Error("same title should generate same ID")
	}

	// Different title should generate different ID
	id3 := tc.GenerateID("DT", "提拉米苏")
	if id == id3 {
		t.Error("different titles should generate different IDs")
	}

	// Empty title should still generate ID
	id4 := tc.GenerateID("DT", "")
	if id4 == "" {
		t.Error("empty title should still generate ID")
	}
}

func TestTitleCleaner_IsPriceMatch(t *testing.T) {
	tc := NewTitleCleaner()

	if !tc.IsPriceMatch(10.0, 10.5) {
		t.Error("prices within threshold should match")
	}

	if !tc.IsPriceMatch(10.0, 10.0) {
		t.Error("identical prices should match")
	}

	// Exactly at threshold should match
	if !tc.IsPriceMatch(10.0, 11.0) {
		t.Error("price at threshold should match")
	}

	if tc.IsPriceMatch(10.0, 12.0) {
		t.Error("prices outside threshold should not match")
	}
}

func TestTitleCleaner_ShouldPromote(t *testing.T) {
	tc := NewTitleCleaner()

	if !tc.ShouldPromote(3) {
		t.Error("occurrences >= threshold should promote")
	}

	if tc.ShouldPromote(2) {
		t.Error("occurrences < threshold should not promote")
	}
}

func TestTitleCleaner_FindMatchingMaster_Empty(t *testing.T) {
	tc := NewTitleCleaner()

	// Empty masters
	result := tc.FindMatchingMaster("巧克力", []*entity.MasterProduct{})
	if result != nil {
		t.Error("should return nil for empty masters")
	}

	// Empty title
	result2 := tc.FindMatchingMaster("", []*entity.MasterProduct{})
	if result2 != nil {
		t.Error("should return nil for empty title")
	}
}

func TestTitleCleaner_FindMatchingMasterWithPrice_NilInSlice(t *testing.T) {
	tc := NewTitleCleaner()

	masters := []*entity.MasterProduct{
		{
			ID:            "1",
			StandardTitle: "巧克力",
			Price:         10.0,
		},
		nil, // Intentionally nil
		{
			ID:            "2",
			StandardTitle: "蛋糕",
			Price:         20.0,
		},
	}

	result := tc.FindMatchingMasterWithPrice("巧克力", 10.0, masters)
	if result == nil {
		t.Error("should find matching master")
	}
}

func TestTitleCleaner_NormalizeTitle(t *testing.T) {
	tc := NewTitleCleaner()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Leading/trailing spaces",
			input:    "  巧克力蛋糕  ",
			expected: "巧克力蛋糕",
		},
		{
			name:     "Multiple spaces",
			input:    "巧克力   蛋糕",
			expected: "巧克力 蛋糕",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tc.NormalizeTitle(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
