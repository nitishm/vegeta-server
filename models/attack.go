package models

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
