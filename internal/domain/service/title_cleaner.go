package service

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"kbfood/internal/domain/entity"
)

const (
	// SimilarityThreshold for title matching
	SimilarityThreshold = 0.75
	// PromotionThreshold for candidate promotion
	PromotionThreshold = 3
	// MidSimilarityThreshold for title + price matching
	MidSimilarityThreshold = 0.5
	// PriceMatchThreshold in yuan
	PriceMatchThreshold = 1.0
)

var (
	// Pre-compiled regex for title cleaning
	titleCleanerRegex *regexp.Regexp
	regexOnce         sync.Once
)

func getTitleCleanerRegex() *regexp.Regexp {
	regexOnce.Do(func() {
		// Use \p{Han} for Chinese characters instead of \u escape
		titleCleanerRegex = regexp.MustCompile(`[^a-zA-Z0-9\p{Han}]`)
	})
	return titleCleanerRegex
}

// TitleCleaner handles title cleaning and candidate pool management
type TitleCleaner struct {
	similarityThreshold float64
	promotionThreshold  int
}

// NewTitleCleaner creates a new title cleaner
func NewTitleCleaner() *TitleCleaner {
	return &TitleCleaner{
		similarityThreshold: SimilarityThreshold,
		promotionThreshold:  PromotionThreshold,
	}
}

// CleanTitleForID cleans a title for ID generation
// Removes all non-alphanumeric characters (except Chinese)
func (tc *TitleCleaner) CleanTitleForID(title string) string {
	if title == "" {
		return ""
	}

	// Use pre-compiled regex for better performance
	reg := getTitleCleanerRegex()
	return reg.ReplaceAllString(title, "")
}

// GenerateID generates a unique ID from a title
func (tc *TitleCleaner) GenerateID(prefix, title string) string {
	cleaned := tc.CleanTitleForID(title)
	hash := md5.Sum([]byte(cleaned))
	return prefix + "_" + hex.EncodeToString(hash[:])
}

// CalculateSimilarity calculates the similarity between two strings
// Using Levenshtein distance
func (tc *TitleCleaner) CalculateSimilarity(s1, s2 string) float64 {
	// Both empty strings have no similarity (as per test expectation)
	if s1 == "" && s2 == "" {
		return 0.0
	}
	if s1 == s2 {
		return 1.0
	}
	if s1 == "" || s2 == "" {
		return 0.0
	}

	distance := levenshteinDistance(s1, s2)
	// Use rune count for proper Unicode (Chinese) character handling
	maxLen := max(len([]rune(s1)), len([]rune(s2)))
	if maxLen == 0 {
		return 0.0
	}

	return 1.0 - float64(distance)/float64(maxLen)
}

// IsHighSimilarity checks if two titles have high similarity
func (tc *TitleCleaner) IsHighSimilarity(title1, title2 string) bool {
	// Handle edge cases
	if title1 == "" || title2 == "" {
		return false
	}
	return tc.CalculateSimilarity(title1, title2) >= tc.similarityThreshold
}

// IsMidSimilarity checks if two titles have medium similarity
func (tc *TitleCleaner) IsMidSimilarity(title1, title2 string) bool {
	if title1 == "" || title2 == "" {
		return false
	}
	sim := tc.CalculateSimilarity(title1, title2)
	return sim >= MidSimilarityThreshold && sim < tc.similarityThreshold
}

// IsPriceMatch checks if two prices match within threshold
func (tc *TitleCleaner) IsPriceMatch(price1, price2 float64) bool {
	// Handle NaN and infinity
	if math.IsNaN(price1) || math.IsNaN(price2) {
		return false
	}
	if math.IsInf(price1, 0) || math.IsInf(price2, 0) {
		return false
	}

	diff := price1 - price2
	if diff < 0 {
		diff = -diff
	}
	return diff <= PriceMatchThreshold // Use <= to match exactly at threshold
}

// ElectStandardTitle selects the standard title from votes
func (tc *TitleCleaner) ElectStandardTitle(votes map[string]int) string {
	if len(votes) == 0 {
		return ""
	}

	var winner string
	maxVotes := 0

	for title, count := range votes {
		if count > maxVotes {
			maxVotes = count
			winner = title
		} else if count == maxVotes && len(title) > len(winner) {
			// Tie-breaker: prefer longer title (more descriptive)
			winner = title
		}
	}

	return winner
}

// ShouldPromote checks if a candidate should be promoted to master
func (tc *TitleCleaner) ShouldPromote(occurrences int) bool {
	return occurrences >= tc.promotionThreshold
}

// FindMatchingMaster finds a matching master product by title
func (tc *TitleCleaner) FindMatchingMaster(
	title string,
	masters []*entity.MasterProduct,
) *entity.MasterProduct {
	// Handle edge cases
	if title == "" || len(masters) == 0 {
		return nil
	}

	// First try high similarity match
	for _, m := range masters {
		if m == nil {
			continue
		}
		if tc.IsHighSimilarity(title, m.StandardTitle) {
			return m
		}
	}
	return nil
}

// FindMatchingMasterWithPrice finds a matching master by title + price
func (tc *TitleCleaner) FindMatchingMasterWithPrice(
	title string,
	price float64,
	masters []*entity.MasterProduct,
) *entity.MasterProduct {
	if title == "" || len(masters) == 0 {
		return nil
	}

	for _, m := range masters {
		if m == nil {
			continue
		}
		// Check for high similarity with price match, or mid similarity with price match
		sim := tc.CalculateSimilarity(title, m.StandardTitle)
		if (sim >= MidSimilarityThreshold) && tc.IsPriceMatch(price, m.Price) {
			return m
		}
	}
	return nil
}

// NormalizeTitle normalizes a title for comparison
// Removes extra spaces, converts to lowercase, etc.
func (tc *TitleCleaner) NormalizeTitle(title string) string {
	// Remove leading/trailing spaces
	title = strings.TrimSpace(title)

	// Remove consecutive spaces
	prev := ' '
	result := make([]rune, 0, len(title))

	for _, r := range title {
		if !unicode.IsSpace(r) || prev != ' ' {
			result = append(result, r)
			prev = r
		}
	}

	return string(result)
}

// levenshteinDistance calculates the Levenshtein distance between two strings
// Uses rune iteration for proper Unicode/Chinese character handling
func levenshteinDistance(s1, s2 string) int {
	// Convert to runes for proper Unicode handling (especially for Chinese)
	r1 := []rune(s1)
	r2 := []rune(s2)

	if len(r1) < len(r2) {
		// Swap to ensure r1 is the longer string
		r1, r2 = r2, r1
	}
	if len(r2) == 0 {
		return len(r1)
	}

	previousRow := make([]int, len(r2)+1)
	for i := range previousRow {
		previousRow[i] = i
	}

	for i, c1 := range r1 {
		currentRow := []int{i + 1}
		for j, c2 := range r2 {
			cost := 0
			if c1 != c2 {
				cost = 1
			}
			currentRow = append(currentRow, min(
				previousRow[j+1]+1,  // deletion
				currentRow[j]+1,     // insertion
				previousRow[j]+cost, // substitution
			))
		}
		previousRow = currentRow
	}

	return previousRow[len(previousRow)-1]
}

// Helper functions
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
