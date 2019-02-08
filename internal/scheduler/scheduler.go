package scheduler

import (
	"vegeta-server/internal"
	"vegeta-server/internal/dispatcher"
	"vegeta-server/internal/models"
)

type SchedulerFn func(dispatcher.IDispatcher, chan internal.ITask, chan struct{})

type attack struct {
}

type IScheduler interface {
	Run(chan struct{})
	Schedule(models.AttackParams) *models.AttackResponse
	Get(string) (*models.AttackResponse, error)
	List() []*models.AttackResponse
	Cancel(string, bool) (*models.AttackResponse, error)
}

// TODO: scheduler must do more to figure out system capacity
type scheduler struct {
	schedulerFn SchedulerFn
	taskCh      chan internal.ITask
	dispatcher  dispatcher.IDispatcher
	quit        chan struct{}
}

func NewScheduler(dispatcher dispatcher.IDispatcher, schedulerFn SchedulerFn) *scheduler {
	s := &scheduler{
		schedulerFn: schedulerFn,
		taskCh:      make(chan internal.ITask),
		dispatcher:  dispatcher,
	}

	return s
}

func (s *scheduler) Run(quit chan struct{}) {
	s.schedulerFn(s.dispatcher, s.taskCh, quit)
}

func (s *scheduler) Schedule(attack models.AttackParams) *models.AttackResponse {
	task := internal.NewTask(attack)

	// Schedule the test
	ok := s.schedule()
	if !ok {
		// TODO: Do something here
	}

	s.taskCh <- task

	// Return the UUID and Status = scheduled
	return &models.AttackResponse{
		ID:     task.ID(),
		Status: task.Status(),
	}
}

func (s *scheduler) Get(id string) (*models.AttackResponse, error) {
	t, err := s.dispatcher.Get(id)
	if err != nil {
		return nil, err
	}

	resp := &models.AttackResponse{
		ID:     t.ID(),
		Status: t.Status(),
	}

	return resp, err
}

func (s *scheduler) List() []*models.AttackResponse {
	responses := make([]*models.AttackResponse, 0)
	for _, task := range s.dispatcher.List() {
		responses = append(responses, &models.AttackResponse{
			ID:     task.ID(),
			Status: task.Status(),
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
		t.ID(),
		t.Status(),
	}, nil
}

func (s *scheduler) schedule() bool {
	// TODO: check system resources to schedule an attacks
	return true
}

func DefaultSchedulerFn(dispatcher dispatcher.IDispatcher, taskCh chan internal.ITask, quit chan struct{}) {
	for {
		select {
		case task := <-taskCh:
			// TODO: Scheduler should check if it can schedule attack
			// 		 If not defer to later (maintain a separate queue)
			dispatcher.Dispatch(task)
		case <-quit:
			break
		}
	}
}
