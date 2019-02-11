package reporter

import (
	"bytes"
	"vegeta-server/models"
	"vegeta-server/pkg/vegeta"
)

type IReporter interface {
	// Get report in (default) JSON format
	Get(string) (string, error)
	// GetAll gets all reports in (default) JSON format
	GetAll() []string

	// Get report in specified format (supported: JSON/Histogram/Text
	GetInFormat(string, vegeta.Format) (string, error)
	// GetAll gets all reports in specified format (supported: JSON/Histogram/Text
	GetAllInFormat(vegeta.Format) []string

	// Delete report from store
	Delete(string) error
}

type reporter struct {
	db models.IAttackStore
}

func NewReporter(db models.IAttackStore) *reporter {
	return &reporter{
		db,
	}
}

func (r *reporter) Get(id string) (string, error) {
	attack, err := r.db.GetById(id)
	if err != nil {
		return "", err
	}

	result := attack.Result
	report, err := vegeta.CreateReportFromReader(bytes.NewBuffer(result), vegeta.JSONFormat)
	if err != nil {
		return "", err
	}
	return report, nil
}

func (r *reporter) GetAll() []string {
	attacks := r.db.GetAll()
	reports := make([]string, 0)
	for _, attack := range attacks {
		report, err := vegeta.CreateReportFromReader(bytes.NewBuffer(attack.Result), vegeta.JSONFormat)
		if err != nil {
			continue
		}
		reports = append(reports, report)
	}
	return reports
}

func (r *reporter) GetInFormat(id string, formats vegeta.Format) (string, error) {
	attack, err := r.db.GetById(id)
	if err != nil {
		return "", err
	}

	result := attack.Result
	report, err := vegeta.CreateReportFromReader(bytes.NewBuffer(result), vegeta.JSONFormat)
	if err != nil {
		return "", err
	}
	return report, nil
}

func (r *reporter) GetAllInFormat(format vegeta.Format) []string {
	attacks := r.db.GetAll()
	reports := make([]string, 0)
	for _, attack := range attacks {
		report, err := vegeta.CreateReportFromReader(bytes.NewBuffer(attack.Result), format)
		if err != nil {
			continue
		}
		reports = append(reports, report)
	}
	return reports
}

func (r *reporter) Delete(id string) error {
	panic("implement me")
}
