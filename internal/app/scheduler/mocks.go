package scheduler

import (
	"fmt"
	"vegeta-server/internal/app/server/models"
)

type MockScheduler struct {
	Responses map[string]*models.AttackResponse
}

func (ms *MockScheduler) Schedule(models.AttackParams) *models.AttackResponse {
	return &models.AttackResponse{
		Status: models.AttackResponseStatusScheduled,
	}
}

func (ms *MockScheduler) Get(id string) (*models.AttackResponse, error) {
	r, ok := ms.Responses[id]
	if !ok {
		return nil, fmt.Errorf("[MOCK] Not found")
	}
	return r, nil
}

func (ms *MockScheduler) List() []*models.AttackResponse {
	responses := make([]*models.AttackResponse, 0)
	for _, resp := range ms.Responses {
		responses = append(responses, resp)
	}

	return responses
}

func (ms *MockScheduler) Cancel(id string, cancel bool) (*models.AttackResponse, error) {
	resp := &models.AttackResponse{
		ID:     id,
		Status: models.AttackResponseStatusRunning,
	}

	if cancel {
		resp.Status = models.AttackResponseStatusCanceled
	}

	return resp, nil
}

func (ms *MockScheduler) Run(quit chan struct{}) {
	for {
		select {
		case <-quit:
			return
		}
	}
}
