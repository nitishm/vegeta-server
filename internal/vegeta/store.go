package vegeta

import (
	"fmt"
)

var DefaultStorage = make(ReportMap)

type StorageInterface interface {
	Add(string, string) error
	Get(string) (string, error)
	List() map[string]string
}

type ReportMap map[string]string

func (r ReportMap) Add(uuid string, report string) error {
	if r == nil {
		return fmt.Errorf("Report Map has not been initialized")
	}
	r[uuid] = report

	return nil
}

func (r ReportMap) Get(uuid string) (string, error) {
	if v, ok := r[uuid]; !ok {
		return "", fmt.Errorf("No attack entry with UUID %v found in Report Map", uuid)
	} else {
		return v, nil
	}
}

func (r ReportMap) List() map[string]string {
	return r
}
