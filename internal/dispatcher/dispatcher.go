package dispatcher

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"sync"
	"vegeta-server/models"
	"vegeta-server/pkg/vegeta"
)

type IDispatcher interface {
	Run(chan struct{})
	Dispatch(models.AttackParams) *models.AttackResponse
	Cancel(string, bool) (*models.AttackResponse, error)

	Get(string) (*models.AttackResponse, error)
	List() []*models.AttackResponse
}

type dispatcher struct {
	mu             *sync.RWMutex
	tasks          map[string]ITask
	attackFn       vegeta.AttackFunc
	newTaskCh      chan ITask
	attackUpdateCh chan models.AttackDetails
	db             models.IAttackStore
}

func NewDispatcher(db models.IAttackStore, fn vegeta.AttackFunc) *dispatcher { //nolint:golint
	d := &dispatcher{
		&sync.RWMutex{},
		make(map[string]ITask),
		fn,
		make(chan ITask, 10),
		make(chan models.AttackDetails, 20),
		db,
	}
	d.log(nil).Info("creating new dispatcher")
	return d
}

func (d *dispatcher) Dispatch(params models.AttackParams) *models.AttackResponse {
	d.mu.Lock()
	defer d.mu.Unlock()
	task := NewTask(d.attackUpdateCh, params)
	id := task.ID()
	status := task.Status()
	fields := log.Fields{
		"ID":     id,
		"Status": status,
	}

	// Track the task
	d.tasks[task.ID()] = task

	// Add to database
	_ = d.db.Add(models.AttackDetails{
		AttackInfo: models.AttackInfo{
			ID:     id,
			Status: status,
			Params: params,
		},
		Result: nil,
	})

	d.log(fields).Info("dispatching new attack")
	d.newTaskCh <- task

	return &models.AttackResponse{
		ID:     task.ID(),
		Status: task.Status(),
		Params: task.Params(),
	}
}

func (d *dispatcher) Get(id string) (*models.AttackResponse, error) {
	fields := log.Fields{
		"ID": id,
	}

	d.log(fields).Debug("getting attack status")

	d.mu.RLock()
	defer d.mu.RUnlock()
	t, ok := d.tasks[id]
	if !ok {
		err := fmt.Errorf("cannot find task with id %s", id)
		d.log(fields).Error("failed to find attack")
		return nil, err
	}
	response := &models.AttackResponse{
		ID:     t.ID(),
		Status: t.Status(),
		Params: t.Params(),
	}
	return response, nil
}

func (d *dispatcher) List() []*models.AttackResponse {
	d.log(nil).Debug("getting attack list")

	d.mu.RLock()
	defer d.mu.RUnlock()
	responses := make([]*models.AttackResponse, 0)
	for _, task := range d.tasks {
		responses = append(responses, &models.AttackResponse{
			ID:     task.ID(),
			Status: task.Status(),
			Params: task.Params(),
		})
	}
	return responses
}

func (d *dispatcher) Cancel(id string, cancel bool) (*models.AttackResponse, error) {
	fields := log.Fields{
		"ID":       id,
		"ToCancel": cancel,
	}

	d.log(fields).Info("canceling attack")

	d.mu.Lock()
	defer d.mu.Unlock()
	t, ok := d.tasks[id]
	if !ok {
		d.log(fields).Error("task not found")
		return nil, fmt.Errorf("cannot find task with id %s", id)
	}

	if cancel {
		err := t.Cancel()
		if err != nil {
			d.log(fields).WithError(err).Error("failed to cancel task")
			return nil, err
		}
	}

	return &models.AttackResponse{
		ID:     t.ID(),
		Status: t.Status(),
		Params: t.Params(),
	}, nil
}

func (d *dispatcher) log(fields map[string]interface{}) *log.Entry {
	l := log.WithField("component", "dispatcher")

	if fields != nil {
		l = l.WithFields(fields)
	}

	return l
}

func (d *dispatcher) Run(quit chan struct{}) {
	defer close(d.newTaskCh)
	d.log(nil).Info("starting dispatcher")
	for {
		select {
		case task := <-d.newTaskCh:
			fields := log.Fields{
				"ID":     task.ID(),
				"Status": task.Status(),
			}

			d.log(fields).Debug("received task")

			if err := task.Run(d.attackFn); err != nil {
				d.log(fields).WithError(err).Errorf("failed to run %s", task.ID())
				continue
			}
		case updatedAttack := <-d.attackUpdateCh:
			fields := log.Fields{
				"ID":     updatedAttack.ID,
				"Status": updatedAttack.Status,
			}

			if err := d.db.Update(updatedAttack.ID, updatedAttack); err != nil {
				d.log(fields).WithError(err).Error("attack update error")
				continue
			}
			d.log(fields).Debug("received update for attack")
		case <-quit:
			d.log(nil).Warning("gracefully shutting down the dispatcher")
			return
		}
	}
}
