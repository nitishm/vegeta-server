package attacker

import (
	"vegeta-server/internal/app/server/models"
)

type attack struct {
}

type IScheduler interface {
	Schedule(models.Attack) *models.AttackResponse
	Get(string) (*models.AttackResponse, error)
	List() []*models.AttackResponse
	Cancel(string, bool) (*models.AttackResponse, error)
}

// TODO: scheduler must do more to figure out system capacity
type scheduler struct {
	ch         chan *task
	dispatcher IDispatcher
	quit       chan struct{}
}

func NewScheduler(dispatcher IDispatcher, quit chan struct{}) *scheduler {
	s := &scheduler{
		ch:         make(chan *task),
		dispatcher: dispatcher,
		quit:       quit,
	}

	go s.attackScheduler()

	return s
}

func (s *scheduler) Schedule(attack models.Attack) *models.AttackResponse {
	task := newTask(attack)

	// Schedule the test
	ok := s.schedule()
	if !ok {
		// TODO: Do something here
	}

	s.ch <- task

	// Return the UUID and Status = scheduled
	return &models.AttackResponse{
		ID:     task.id,
		Status: task.status,
	}
}

func (s *scheduler) Get(id string) (*models.AttackResponse, error) {
	t, err := s.dispatcher.Get(id)
	if err != nil {
		return nil, err
	}

	resp := &models.AttackResponse{
		ID:     t.id,
		Status: t.status,
	}

	return resp, err
}

func (s *scheduler) List() []*models.AttackResponse {
	responses := make([]*models.AttackResponse, 0)
	for _, task := range s.dispatcher.List() {
		responses = append(responses, &models.AttackResponse{
			ID:     task.getID(),
			Status: task.getStatus(),
		})
	}

	return responses
}

func (s *scheduler) Cancel(id string, cancel bool) (*models.AttackResponse, error) {
	t, err := s.dispatcher.Cancel(id, cancel)
	if err != nil {
		return nil, err
	}

	return &models.AttackResponse{
		t.getID(),
		t.getStatus(),
	}, nil
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

func (s *scheduler) schedule() bool {
	// TODO: check system resources to schedule an attacks
	return true
}
