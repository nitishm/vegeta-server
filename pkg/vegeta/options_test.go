package vegeta

import (
	"net/http"
	"reflect"
	"testing"
	"time"
	"vegeta-server/models"

	vegeta "github.com/tsenart/vegeta/lib"
)

func TestNewAttackOptsFromAttackParams_WithHeaders(t *testing.T) {
	tests := []struct {
		name   string
		params models.AttackParams
		want   AttackOpts
	}{
		{"with-headers", models.AttackParams{
			Rate:     5,
			Duration: "10s",
			Headers: []models.AttackHeader{
				{
					Key:   "X-Test-Key",
					Value: "test-value",
				},
			},
		}, AttackOpts{
			Name:     "with-headers",
			Duration: 10 * time.Second,
			Rate: vegeta.Rate{
				Freq: 5,
				Per:  time.Second,
			},
			Target: vegeta.Target{
				Header: http.Header{
					"X-Test-Key": []string{"test-value"},
				},
			},
		}},
		{"without-headers", models.AttackParams{
			Rate:     5,
			Duration: "10s",
		}, AttackOpts{
			Name:     "without-headers",
			Duration: 10 * time.Second,
			Rate: vegeta.Rate{
				Freq: 5,
				Per:  time.Second,
			},
			Target: vegeta.Target{
				Header: make(http.Header),
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAttackOptsFromAttackParams(tt.name, tt.params)
			if !reflect.DeepEqual(got.Target.Header, tt.want.Target.Header) {
				t.Errorf("NewAttackOptsFromAttackParams() = %v, want %v", got.Target.Header, tt.want.Target.Header)
			} else if err != nil {
				t.Error(err)
			}
		})
	}
}
