package models

import (
	"time"
)

// AttackInfo encapsulates the attack information for attacks
// submitted to the dispatcher
type AttackInfo struct {
	// id is a attack UUID generated for each attack submitted
	ID string `json:"id,omitempty"`
	// Status captures the attack status in the scheduler pipeline
	Status AttackStatus `json:"status,omitempty"`
	// Params captures the attack parameters
	Params    AttackParams `json:"params,omitempty"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}

// AttackDetails captures the AttackInfo for COMPLETED attacks,
// along with the result as a byte array
type AttackDetails struct {
	AttackInfo
	Result []byte `json:"result,omitempty"`
}

// FilterParams defines a map structure for the filter parameters received via
// query params in the request URL
type FilterParams map[string]interface{}

// Filter defines a type that must be implemented by
// an attack filter
type Filter func(AttackDetails) bool

// StatusFilter implements a attack status filter
// in the Filter function format
func StatusFilter(status string) Filter {
	return func(a AttackDetails) bool {
		if status == "" {
			return true
		}

		if a.Status == AttackStatus(status) {
			return true
		}

		return false
	}
}

// CreationBeforeFilter implements an attack created_before filter
// in the Filter function format
func CreationBeforeFilter(d string) Filter {
	return func(a AttackDetails) bool {
		if d == "" {
			return true
		}

		const layoutUser = "2006-01-02 15:04:05"
		t, err := time.ParseInLocation(layoutUser, d, time.Local)
		// If parsing failed, don't filter
		if err != nil {
			return true
		}

		attackTime, _ := time.Parse(time.RFC1123, a.CreatedAt)

		return attackTime.Before(t)
	}
}

// CreationAfterFilter implements an attack created_after filter
// in the Filter function format
func CreationAfterFilter(d string) Filter {
	return func(a AttackDetails) bool {
		if d == "" {
			return true
		}
		const layoutUser = "2006-01-02 15:04:05"
		t, err := time.Parse(layoutUser, d)

		// If parsing failed, don't filter
		if err != nil {
			return true
		}

		attackTime, _ := time.Parse(time.RFC1123, a.CreatedAt)

		return attackTime.After(t)
	}
}
