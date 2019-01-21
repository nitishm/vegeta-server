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
		return fmt.Errorf("report map has not been initialized")
	}
	r[uuid] = report

	return nil
}

func (r ReportMap) Get(uuid string) (string, error) {
	v, ok := r[uuid]
	if !ok {
		return "", fmt.Errorf("no attack entry with UUID %v found in report map", uuid)
	}

	return v, nil
}

func (r ReportMap) List() map[string]string {
	return r
}
