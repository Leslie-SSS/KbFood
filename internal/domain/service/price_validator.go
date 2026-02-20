package service

import (
	"math"
	"time"

	"kbfood/internal/pkg/errors"
)

// PriceValidator implements Dutch auction model price validation
// The model assumes prices only decrease (one-way downward)
type PriceValidator struct {
	maxDropRatio float64 // Maximum drop ratio (default: 0.5 = 50%)
	maxRiseRatio float64 // Maximum rise ratio (default: 5.0 = 5x)
	minPrice     float64 // Minimum valid price (default: 1.0)
}

// NewPriceValidator creates a new price validator with default settings
func NewPriceValidator() *PriceValidator {
	return &PriceValidator{
		maxDropRatio: 0.5, // Single drop > 50% is considered noise
		maxRiseRatio: 5.0, // Single rise > 5x is considered error
		minPrice:     1.0, // Minimum 1 yuan
	}
}

// NewPriceValidatorWithConfig creates a new price validator with custom settings
func NewPriceValidatorWithConfig(maxDrop, maxRise, min float64) *PriceValidator {
	return &PriceValidator{
		maxDropRatio: maxDrop,
		maxRiseRatio: maxRise,
		minPrice:     min,
	}
}

// ValidateUpdate validates a price update based on the Dutch auction model
//
// Rules:
// 1. Cross-day: Always trust (price resets)
// 2. Same-day drop: Validate against noise threshold
// 3. Same-day rise: Allow for error correction, but limit extreme rises
//
// Returns the final price to use and any error
func (v *PriceValidator) ValidateUpdate(
	oldPrice, newPrice float64,
	lastUpdateTime time.Time,
) (float64, error) {
	// Validate inputs
	if math.IsNaN(oldPrice) || math.IsNaN(newPrice) {
		return oldPrice, errors.New(errors.InvalidInput, "price is NaN")
	}
	if math.IsInf(oldPrice, 0) || math.IsInf(newPrice, 0) {
		return oldPrice, errors.New(errors.InvalidInput, "price is infinite")
	}
	if newPrice < 0 {
		return oldPrice, errors.New(errors.ErrPriceBelowMin, "price is negative")
	}

	now := time.Now()
	oldDate := lastUpdateTime.UTC().Truncate(24 * time.Hour)
	today := now.UTC().Truncate(24 * time.Hour)

	// Rule 0: Cross-day - unconditional trust (price reset)
	// But still validate the new price is valid
	if !oldDate.Equal(today) {
		// Apply same basic validation as same-day
		if newPrice < v.minPrice {
			return oldPrice, errors.New(errors.ErrPriceBelowMin, "price below minimum threshold")
		}
		// Cross-day allows any positive price (price reset)
		return newPrice, nil
	}

	// Rule 1: Minimum price protection (applies to both cross-day and same-day)
	if newPrice < v.minPrice {
		return oldPrice, errors.New(errors.ErrPriceBelowMin, "price below minimum threshold")
	}

	// Handle zero or negative old price - same day
	if oldPrice <= 0 {
		// If old price is invalid, accept new price if it passes basic validation
		// This means we're setting the initial price
		return newPrice, nil
	}

	// Rule 2: Block extreme drops (prevent 110 -> 48 scenario)
	// Only check if price is dropping
	if newPrice < oldPrice {
		dropRatio := newPrice / oldPrice
		if dropRatio < v.maxDropRatio {
			return oldPrice, errors.New(errors.ErrPriceDropExceeded,
				"price drop exceeded threshold")
		}
	}

	// Rule 3: Allow corrective rises (prevent OCR errors from persisting)
	// In Dutch auction model, same-day rises are impossible in business logic
	// So they must be error corrections
	if newPrice > oldPrice {
		// But prevent sky-high errors (e.g., 48 -> 4800)
		// Only check if old price is significant enough
		if oldPrice >= 10 {
			riseRatio := newPrice / oldPrice
			if riseRatio > v.maxRiseRatio {
				return oldPrice, errors.New(errors.ErrPriceRiseExceeded,
					"price rise exceeded threshold")
			}
		}
		return newPrice, nil
	}

	// Rule 4: Normal price drop - allow
	return newPrice, nil
}

// IsSameDay checks if two times are on the same day
func (v *PriceValidator) IsSameDay(t1, t2 time.Time) bool {
	return t1.UTC().Truncate(24 * time.Hour).Equal(t2.UTC().Truncate(24 * time.Hour))
}

// IsCrossDay checks if two times are on different days
func (v *PriceValidator) IsCrossDay(t1, t2 time.Time) bool {
	return !v.IsSameDay(t1, t2)
}

// ValidatePriceChange validates a price change and returns if it's valid
func (v *PriceValidator) ValidatePriceChange(oldPrice, newPrice float64) error {
	// Validate inputs
	if math.IsNaN(oldPrice) || math.IsNaN(newPrice) {
		return errors.New(errors.InvalidInput, "price is NaN")
	}
	if math.IsInf(oldPrice, 0) || math.IsInf(newPrice, 0) {
		return errors.New(errors.InvalidInput, "price is infinite")
	}
	if newPrice < 0 || oldPrice < 0 {
		return errors.New(errors.ErrPriceBelowMin, "price is negative")
	}

	if newPrice < v.minPrice {
		return errors.New(errors.ErrPriceBelowMin, "price below minimum threshold")
	}

	// Handle zero old price
	if oldPrice == 0 {
		return nil // Accept any new price if old price was zero
	}

	// Check for extreme drop
	if newPrice < oldPrice {
		dropRatio := newPrice / oldPrice
		if dropRatio < v.maxDropRatio {
			return errors.New(errors.ErrPriceDropExceeded, "price drop exceeded threshold")
		}
	}

	// Check for extreme rise
	if newPrice > oldPrice && oldPrice >= 10 {
		riseRatio := newPrice / oldPrice
		if riseRatio > v.maxRiseRatio {
			return errors.New(errors.ErrPriceRiseExceeded, "price rise exceeded threshold")
		}
	}

	return nil
}

// GetEffectivePrice returns the effective price after validation
// If validation fails, returns the old price
func (v *PriceValidator) GetEffectivePrice(
	oldPrice, newPrice float64,
	lastUpdateTime time.Time,
) float64 {
	price, _ := v.ValidateUpdate(oldPrice, newPrice, lastUpdateTime)
	return price
}
