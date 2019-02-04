package attacker

import (
	"vegeta-server/internal/app/server/models"
)

// IAttacker is the attacker interface
type IAttacker interface {
	Submit(models.Attack) *models.AttackResponse
	Get(string) (*models.AttackResponse, error)
	List() []*models.AttackResponse
	Cancel(string, bool) (*models.AttackResponse, error)
}

type attacker struct {
	// scheduler schedules an attack to run
	scheduler IScheduler
}

// NewAttacker returns an implementation of the
// IAttacker interface
func NewAttacker(scheduler IScheduler) *attacker {
	return &attacker{
		scheduler,
	}
}

func (a *attacker) Submit(attack models.Attack) *models.AttackResponse {
	return a.scheduler.Schedule(attack)
}

func (a *attacker) Get(id string) (*models.AttackResponse, error) {
	return a.scheduler.Get(id)
}

func (a *attacker) List() []*models.AttackResponse {
	return a.scheduler.List()
}

func (a *attacker) Cancel(id string, cancel bool) (*models.AttackResponse, error) {
	return a.scheduler.Cancel(id, cancel)
}
