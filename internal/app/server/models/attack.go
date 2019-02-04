package models

// Attack request parameters
type Attack struct {
	Rate     int    `json:"rate,omitempty"`
	Duration string `json:"duration,omitempty"`
	Target   Target `json:"target,omitempty"`
}

// Attack request target parameters
type Target struct {
	Method string `json:"method,omitempty"`
	URL    string `json:"URL,omitempty"`
	Scheme string `json:"scheme,omitempty"`
}

// AttackStatus as a string enum
type AttackStatus string

const (

	// AttackResponseStatusScheduled captures enum value "scheduled"
	AttackResponseStatusScheduled AttackStatus = "scheduled"

	// AttackResponseStatusRunning captures enum value "running"
	AttackResponseStatusRunning AttackStatus = "running"

	// AttackResponseStatusCanceled captures enum value "canceled"
	AttackResponseStatusCanceled AttackStatus = "canceled"

	// AttackResponseStatusCompleted captures enum value "completed"
	AttackResponseStatusCompleted  AttackStatus= "completed"

	// AttackResponseStatusFailed captures enum value "failed"
	AttackResponseStatusFailed AttackStatus = "failed"
)

// AttackResponse with attacks UUID and AttackStatus
type AttackResponse struct {
	// id is a attack UUID generated for each attack submitted
	ID string `json:"id,omitempty"`
	// Status captures the attack status in the scheduler pipeline
	Status AttackStatus `json:"status,omitempty"`
}