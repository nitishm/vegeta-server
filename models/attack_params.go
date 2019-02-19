package models

// AttackHeader provides a key/value object for headers
type AttackHeader struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// AttackParams request parameters
type AttackParams struct {
	Rate int `json:"rate,omitempty" binding:"required"`

	Connections int64 `json:"connections,omitempty"`
	Workers     int64 `json:"workers,omitempty"`
	MaxBody     int64 `json:"max-body,omitempty"`
	Redirects   int64 `json:"redirects,omitempty"`

	Key       string `json:"key,omitempty"`
	Laddr     string `json:"laddr,omitempty"`
	Duration  string `json:"duration,omitempty" binding:"required"`
	Body      string `json:"body,omitempty"`
	Cert      string `json:"cert,omitempty"`
	Resolvers string `json:"resolvers,omitempty"`
	RootCerts string `json:"root-certs,omitempty"`
	Timeout   string `json:"timeout,omitempty"`

	H2c       bool `json:"h2c,omitempty"`
	HTTP2     bool `json:"http2,omitempty"`
	Insecure  bool `json:"insecure,omitempty"`
	Keepalive bool `json:"keepalive,omitempty"`

	Target  Target         `json:"target,omitempty" binding:"required"`
	Headers []AttackHeader `json:"headers,omitempty"`
}

// Target request target parameters
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
	AttackResponseStatusCompleted AttackStatus = "completed"

	// AttackResponseStatusFailed captures enum value "failed"
	AttackResponseStatusFailed AttackStatus = "failed"
)

// AttackCancel request body
type AttackCancel struct {
	Cancel bool `json:"cancel" binding:"required"`
}
