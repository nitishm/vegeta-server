package dispatcher

import (
	"fmt"
	"vegeta-server/pkg/vegeta"

	log "github.com/sirupsen/logrus"

	"sync"
	"vegeta-server/models"
)

var (
	defaultDB       = models.NewTaskMap()
	defaultAttackFn = vegeta.Attack
)

// IDispatcher provides an interface for attack dispatch operations.
type IDispatcher interface {
	// Run the dispatcher event loop
	Run(chan struct{})
	// Dispatch an attack. Used by the client/handler
	Dispatch(models.AttackParams) (*models.AttackResponse, error)
	// Cancel a scheduled/on-going attack
	Cancel(string, bool) error

	// Get the attack status, params and ID for a single attack
	Get(string) (*models.AttackResponse, error)
	// List the attack status, params and ID for all submitted attacks.
	List() []*models.AttackResponse
}

type dispatcher struct {
	mu       *sync.RWMutex
	tasks    map[string]ITask
	attackFn AttackFunc
	submitCh chan ITask
	updateCh chan UpdateMessage
	db       models.IAttackStore
}

// NewDispatcher constructs a new instance of the dispatcher object.
func NewDispatcher(db models.IAttackStore, fn AttackFunc) *dispatcher { // nolint: golint
	if db == nil {
		db = defaultDB
	}

	if fn == nil {
		fn = defaultAttackFn
	}

	d := &dispatcher{
		&sync.RWMutex{},
		make(map[string]ITask),
		fn,
		make(chan ITask, 10),
		make(chan UpdateMessage, 20),
		db,
	}
	d.log(nil).Info("creating new dispatcher")
	return d
}

// Dispatch implements the attack dispatcher method, used by the client to schedule new attacks
func (d *dispatcher) Dispatch(params models.AttackParams) (*models.AttackResponse, error) {
	task := NewTask(d.updateCh, params)
	id := task.ID()
	status := task.Status()
	fields := log.Fields{
		"ID":     id,
		"Status": status,
	}

	// Track the task
	d.mu.Lock()
	d.tasks[task.ID()] = task
	d.mu.Unlock()

	// Add to database
	_ = d.db.Add(attackDetailFromTask(task))

	d.log(fields).Info("dispatching new attack")
	d.submitCh <- task

	attackDetails, err := d.db.GetByID(id)
	if err != nil {
		return nil, err
	}
	resp := models.AttackResponse(attackDetails.AttackInfo)
	return &resp, nil
}

// Run the dispatcher event loop to dispatch new attacks,
// receive status updates for scheduled attacks and update the
// storage.
func (d *dispatcher) Run(quit chan struct{}) {
	defer close(d.submitCh)
	d.log(nil).Info("starting dispatcher")
	for {
		select {
		case task := <-d.submitCh:
			fields := log.Fields{
				"ID":     task.ID(),
				"Status": task.Status(),
			}

			d.log(fields).Debug("received task")

			if err := task.Run(d.attackFn); err != nil {
				d.log(fields).WithError(err).Errorf("failed to run %s", task.ID())
				continue
			}
		case update := <-d.updateCh:
			fields := log.Fields{
				"ID":     update.ID,
				"Status": update.Status,
			}

			d.mu.RLock()
			task := d.tasks[update.ID]
			d.mu.RUnlock()

			if err := d.db.Update(task.ID(), attackDetailFromTask(task)); err != nil {
				d.log(fields).WithError(err).Error("attack update error")
				continue
			}
			d.log(fields).Debug("received update for attack")
		case <-quit:
			d.mu.RLock()
			for _, task := range d.tasks {
				task.Cancel()
			}
			d.mu.RUnlock()
			d.log(nil).Warning("gracefully shutting down the dispatcher")
			return
		}
	}
}

// Cancel an attack by ID.
func (d *dispatcher) Cancel(id string, cancel bool) error {
	fields := log.Fields{
		"ID":       id,
		"ToCancel": cancel,
	}

	d.log(fields).Info("canceling attack")

	d.mu.RLock()
	t, ok := d.tasks[id]
	if !ok {
		d.log(fields).Error("task not found")
		d.mu.RUnlock()
		return fmt.Errorf("cannot find task with id %s", id)
	}
	d.mu.RUnlock()

	if cancel {
		err := t.Cancel()
		if err != nil {
			d.log(fields).WithError(err).Error("failed to cancel task")
			return err
		}
	}

	return nil
}

// Get an attack by ID
func (d *dispatcher) Get(id string) (*models.AttackResponse, error) {
	fields := log.Fields{
		"ID": id,
	}

	d.log(fields).Debug("getting attack status")

	attackDetails, err := d.db.GetByID(id)
	if err != nil {
		return nil, err
	}
	resp := models.AttackResponse(attackDetails.AttackInfo)
	return &resp, nil
}

// List all submitted attacks
func (d *dispatcher) List() []*models.AttackResponse {
	d.log(nil).Debug("getting attack list")

	responses := make([]*models.AttackResponse, 0)

	for _, attackDetails := range d.db.GetAll() {
		resp := models.AttackResponse(attackDetails.AttackInfo)
		responses = append(responses, &resp)
	}
	return responses
}

func (d *dispatcher) log(fields map[string]interface{}) *log.Entry {
	l := log.WithField("component", "dispatcher")

	if fields != nil {
		l = l.WithFields(fields)
	}

	return l
}
