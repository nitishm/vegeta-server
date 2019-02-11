package models

type AttackInfo struct {
	// id is a attack UUID generated for each attack submitted
	ID string `json:"id,omitempty"`
	// Status captures the attack status in the scheduler pipeline
	Status AttackStatus `json:"status,omitempty"`
	// Params captures the attack parameters
	Params AttackParams `json:"params,omitempty"`
}

type AttackDetails struct {
	AttackInfo
	Result []byte `json:"result,omitempty"`
}
