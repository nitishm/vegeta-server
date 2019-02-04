package attacker

import (
	"vegeta-server/internal/app/server/models"
)

type attack struct {

}

type IScheduler interface {
	Schedule(attack models.Attack) models.AttackResponse
}

// TODO: scheduler must do more to figure out system capacity
type scheduler struct {
	ch chan *task
	dispatcher IDispatcher
	quit chan struct{}
}

func NewScheduler(dispatcher IDispatcher, quit chan struct{}) *scheduler {
	s := &scheduler{
		ch:      make(chan *task),
		dispatcher: dispatcher,
		quit:    quit,
	}

	go s.attackScheduler()

	return s
}

func (s *scheduler) Schedule(attack models.Attack) models.AttackResponse {
	task := newTask(attack)
	// Schedule the test
	s.ch <- task

	// Return the UUID and Status = scheduled
	return models.AttackResponse{
		task.id,
		task.status,
	}
}

func (s *scheduler) attackScheduler() {
	for {
		select {
		case task := <-s.ch:
			// TODO: Scheduler should check if it can schedule attack
			// 		 If not defer to later (maintain a separate queue)
			s.dispatcher.Dispatch(task)
		case <-s.quit:
			break
		}
	}
}