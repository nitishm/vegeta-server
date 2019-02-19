package reporter

import (
	"bytes"
	"vegeta-server/models"
	"vegeta-server/pkg/vegeta"
)

// IReporter provides an interface for all report generation operations.
type IReporter interface {
	// Get report in (default) JSON format
	Get(string) ([]byte, error)
	// GetAll gets all reports in (default) JSON format
	GetAll() [][]byte

	// Get report in specified format (supported: JSON/Histogram/Text
	GetInFormat(string, vegeta.Format) ([]byte, error)

	// Delete report from store
	Delete(string) error
}

type reporter struct {
	db models.IAttackStore
}

// NewReporter returns an instance of the reporter object
func NewReporter(db models.IAttackStore) *reporter { //nolint: golint
	return &reporter{
		db,
	}
}

// Get returns an attack report by its ID as a byte array
func (r *reporter) Get(id string) ([]byte, error) {
	attack, err := r.db.GetByID(id)
	if err != nil {
		return nil, err
	}

	result := attack.Result
	report, err := vegeta.CreateReportFromReader(bytes.NewBuffer(result), attack.ID, vegeta.JSONFormat)
	if err != nil {
		return nil, err
	}
	return report, nil
}

// GetAll returns a list of attack reports in byte array format
// The default format, JSON is returned.
func (r *reporter) GetAll() [][]byte {
	attacks := r.db.GetAll()
	reports := make([][]byte, 0)
	for _, attack := range attacks {
		// Canceled attacks will have a nil result field
		if attack.Result == nil {
			continue
		}

		// Create report for all other attacks
		report, err := vegeta.CreateReportFromReader(bytes.NewBuffer(attack.Result), attack.ID, vegeta.JSONFormat)
		if err != nil {
			continue
		}
		reports = append(reports, report)
	}
	return reports
}

// GetInFormat returns a report in the specified format.
func (r *reporter) GetInFormat(id string, format vegeta.Format) ([]byte, error) {
	attack, err := r.db.GetByID(id)
	if err != nil {
		return nil, err
	}

	result := attack.Result
	if format == vegeta.BinaryFormat {
		return result, nil
	}

	report, err := vegeta.CreateReportFromReader(bytes.NewBuffer(result), attack.ID, format)
	if err != nil {
		return nil, err
	}
	return report, nil
}

// Delete removes a report from the storage
func (r *reporter) Delete(id string) error {
	return r.db.Delete(id)
}
