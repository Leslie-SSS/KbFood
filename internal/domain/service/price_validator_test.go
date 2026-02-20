package service

import (
	"math"
	"testing"
	"time"

	"kbfood/internal/pkg/errors"
)

func TestPriceValidator_ValidateUpdate_CrossDay(t *testing.T) {
	v := NewPriceValidator()

	oldPrice := 110.0
	newPrice := 48.0
	lastUpdate := time.Now().Add(-25 * time.Hour) // Cross day

	finalPrice, err := v.ValidateUpdate(oldPrice, newPrice, lastUpdate)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if finalPrice != newPrice {
		t.Errorf("expected new price %f, got %f", newPrice, finalPrice)
	}
}

func TestPriceValidator_ValidateUpdate_SameDayDrop(t *testing.T) {
	v := NewPriceValidator()

	oldPrice := 110.0
	newPrice := 48.0 // 56% drop, should be blocked
	lastUpdate := time.Now()

	finalPrice, err := v.ValidateUpdate(oldPrice, newPrice, lastUpdate)

	if err == nil {
		t.Error("expected error for excessive drop")
	}
	if err != nil && err.(*errors.AppError).Code != errors.ErrPriceDropExceeded {
		t.Errorf("expected ErrPriceDropExceeded, got %v", err)
	}
	if finalPrice != oldPrice {
		t.Errorf("expected old price %f, got %f", oldPrice, finalPrice)
	}
}

func TestPriceValidator_ValidateUpdate_SameDayValidDrop(t *testing.T) {
	v := NewPriceValidator()

	oldPrice := 100.0
	newPrice := 70.0 // 30% drop, should be allowed
	lastUpdate := time.Now()

	finalPrice, err := v.ValidateUpdate(oldPrice, newPrice, lastUpdate)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if finalPrice != newPrice {
		t.Errorf("expected new price %f, got %f", newPrice, finalPrice)
	}
}

func TestPriceValidator_ValidateUpdate_SameDayRise(t *testing.T) {
	v := NewPriceValidator()

	oldPrice := 48.0
	newPrice := 110.0 // Rise for error correction
	lastUpdate := time.Now()

	finalPrice, err := v.ValidateUpdate(oldPrice, newPrice, lastUpdate)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if finalPrice != newPrice {
		t.Errorf("expected new price %f, got %f", newPrice, finalPrice)
	}
}

func TestPriceValidator_ValidateUpdate_ExtremeRise(t *testing.T) {
	v := NewPriceValidator()

	oldPrice := 48.0
	newPrice := 480.0 // 10x rise, should be blocked
	lastUpdate := time.Now()

	finalPrice, err := v.ValidateUpdate(oldPrice, newPrice, lastUpdate)

	if err == nil {
		t.Error("expected error for excessive rise")
	}
	if err != nil && err.(*errors.AppError).Code != errors.ErrPriceRiseExceeded {
		t.Errorf("expected ErrPriceRiseExceeded, got %v", err)
	}
	if finalPrice != oldPrice {
		t.Errorf("expected old price %f, got %f", oldPrice, finalPrice)
	}
}

func TestPriceValidator_ValidateUpdate_MinPrice(t *testing.T) {
	v := NewPriceValidator()

	oldPrice := 100.0
	newPrice := 0.5 // Below minimum
	lastUpdate := time.Now()

	finalPrice, err := v.ValidateUpdate(oldPrice, newPrice, lastUpdate)

	if err == nil {
		t.Error("expected error for below minimum price")
	}
	if err != nil && err.(*errors.AppError).Code != errors.ErrPriceBelowMin {
		t.Errorf("expected ErrPriceBelowMin, got %v", err)
	}
	if finalPrice != oldPrice {
		t.Errorf("expected old price %f, got %f", oldPrice, finalPrice)
	}
}

func TestPriceValidator_ValidateUpdate_ZeroOldPrice(t *testing.T) {
	v := NewPriceValidator()

	oldPrice := 0.0
	newPrice := 10.0
	lastUpdate := time.Now()

	finalPrice, err := v.ValidateUpdate(oldPrice, newPrice, lastUpdate)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if finalPrice != newPrice {
		t.Errorf("expected new price %f, got %f", newPrice, finalPrice)
	}
}

func TestPriceValidator_ValidateUpdate_NegativePrice(t *testing.T) {
	v := NewPriceValidator()

	oldPrice := 100.0
	newPrice := -10.0
	lastUpdate := time.Now()

	finalPrice, err := v.ValidateUpdate(oldPrice, newPrice, lastUpdate)

	if err == nil {
		t.Error("expected error for negative price")
	}
	if finalPrice != oldPrice {
		t.Errorf("expected old price %f, got %f", oldPrice, finalPrice)
	}
}

func TestPriceValidator_ValidateUpdate_NaNPrice(t *testing.T) {
	v := NewPriceValidator()

	oldPrice := 100.0
	newPrice := math.NaN()
	lastUpdate := time.Now()

	finalPrice, err := v.ValidateUpdate(oldPrice, newPrice, lastUpdate)

	if err == nil {
		t.Error("expected error for NaN price")
	}
	if finalPrice != oldPrice {
		t.Errorf("expected old price %f, got %f", oldPrice, finalPrice)
	}
}

func TestPriceValidator_ValidateUpdate_InfPrice(t *testing.T) {
	v := NewPriceValidator()

	oldPrice := 100.0
	newPrice := math.Inf(1)
	lastUpdate := time.Now()

	finalPrice, err := v.ValidateUpdate(oldPrice, newPrice, lastUpdate)

	if err == nil {
		t.Error("expected error for infinite price")
	}
	if finalPrice != oldPrice {
		t.Errorf("expected old price %f, got %f", oldPrice, finalPrice)
	}
}

func TestPriceValidator_ValidateUpdate_SmallPriceRise(t *testing.T) {
	v := NewPriceValidator()

	oldPrice := 5.0  // Below 10, so no rise limit
	newPrice := 50.0 // 10x rise, should be allowed
	lastUpdate := time.Now()

	finalPrice, err := v.ValidateUpdate(oldPrice, newPrice, lastUpdate)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if finalPrice != newPrice {
		t.Errorf("expected new price %f, got %f", newPrice, finalPrice)
	}
}

func TestPriceValidator_IsSameDay(t *testing.T) {
	v := NewPriceValidator()

	t1 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC)
	t3 := time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC)

	if !v.IsSameDay(t1, t2) {
		t.Error("expected same day")
	}
	if v.IsSameDay(t1, t3) {
		t.Error("expected different day")
	}
}

func TestPriceValidator_ValidatePriceChange_ZeroOldPrice(t *testing.T) {
	v := NewPriceValidator()

	err := v.ValidatePriceChange(0, 10.0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPriceValidator_ValidatePriceChange_NegativePrice(t *testing.T) {
	v := NewPriceValidator()

	err := v.ValidatePriceChange(100.0, -10.0)
	if err == nil {
		t.Error("expected error for negative price")
	}
}

func TestPriceValidator_GetEffectivePrice(t *testing.T) {
	v := NewPriceValidator()

	oldPrice := 100.0
	newPrice := 80.0
	lastUpdate := time.Now()

	price := v.GetEffectivePrice(oldPrice, newPrice, lastUpdate)
	if price != newPrice {
		t.Errorf("expected %f, got %f", newPrice, price)
	}
}
