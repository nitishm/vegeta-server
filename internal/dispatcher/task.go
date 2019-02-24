package dispatcher

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"

	log "github.com/sirupsen/logrus"

	"io"
	"io/ioutil"
	"vegeta-server/models"
)

// AttackFunc provides type used by the attacker class
type AttackFunc func(string, models.AttackParams, chan struct{}) (io.Reader, error)

// ITask defines an interface for attack tasks
type ITask interface {
	ITaskGetter
	ITaskActions
}

// ITaskGetter defines an interface for the task getter methods
type ITaskGetter interface {
	// ID returns the attack task ID
	ID() string
	// Status returns the attack task status
	Status() models.AttackStatus
	// Params returns the attack task params
	Params() models.AttackParams
	// CreatedAt returns the created at timestamp
	CreatedAt() time.Time
	// UpdatedAt returns the updated at timestamp
	UpdatedAt() time.Time
	// Result returns the result as a byte array
	Result() io.Reader
}

// ITaskActions defines an interface for the task action methods
type ITaskActions interface {
	// Run the attack using the configured attack function.
	Run(AttackFunc) error
	// Complete changes task status to completed
	Complete(io.Reader) error
	// Cancel changes task status to canceled
	Cancel() error
	// Fail changes task status to failed
	Fail() error
	// SendUpdate sends an update on the update chan to the caller
	SendUpdate()
}

// UpdateMessage is a message type used to send updates to the dispatcher
// regarding any status changes.
type UpdateMessage struct {
	ID     string
	Status models.AttackStatus
}

type task struct {
	mu     sync.RWMutex
	id     string
	params models.AttackParams
	status models.AttackStatus
	result *bytes.Buffer

	createdAt time.Time
	updatedAt time.Time

	updateCh chan UpdateMessage
	quit     chan struct{}
}

// NewTask returns a new instance of a task object
func NewTask(updateCh chan UpdateMessage, params models.AttackParams) *task { //nolint: golint
	id := uuid.NewV4().String()
	t := &task{
		sync.RWMutex{},
		id,
		params,
		models.AttackResponseStatusScheduled,
		bytes.NewBuffer(make([]byte, 0)),

		time.Now(),
		time.Now(),

		updateCh,
		make(chan struct{}),
	}

	t.log(nil).Debug("creating new task")

	return t
}

// Run an attack task using the passed in attack function
func (t *task) Run(fn AttackFunc) error {
	status := t.Status()
	id := t.ID()

	if status != models.AttackResponseStatusScheduled {
		return fmt.Errorf("cannot run task %s with status %s", id, status)
	}

	t.log(nil).Debug("running")

	go run(t, fn) //nolint: errcheck

	t.mu.Lock()
	t.status = models.AttackResponseStatusRunning
	t.mu.Unlock()

	t.SendUpdate()

	return nil
}

// Complete marks a task as completed
func (t *task) Complete(result io.Reader) error {
	status := t.Status()
	id := t.ID()

	if status != models.AttackResponseStatusRunning {
		return fmt.Errorf("cannot mark completed for task %s with status %s", id, status)
	}

	buf, err := ioutil.ReadAll(result)
	if err != nil {
		return err
	}

	t.mu.Lock()
	t.status = models.AttackResponseStatusCompleted
	t.result = bytes.NewBuffer(buf)
	t.mu.Unlock()

	t.SendUpdate()

	t.log(nil).Debug("completed")

	return nil
}

// Cancel invokes the context cancel and marks a task as canceled
func (t *task) Cancel() error {
	status := t.Status()
	id := t.ID()

	if status == models.AttackResponseStatusCompleted || status == models.AttackResponseStatusFailed || status == models.AttackResponseStatusCanceled { // nolint: lll
		return fmt.Errorf("cannot cancel task %s with status %s", id, status)
	}

	t.quit <- struct{}{}

	t.mu.Lock()
	t.status = models.AttackResponseStatusCanceled
	t.mu.Unlock()

	t.SendUpdate()

	t.log(nil).Debug("canceled")

	return nil
}

// Fail marks a task as failed
func (t *task) Fail() error {
	t.mu.Lock()
	t.status = models.AttackResponseStatusFailed
	t.mu.Unlock()

	t.SendUpdate()

	t.log(nil).Error("failed")
	return nil
}

// SendUpdate to send a status update on the update channel
func (t *task) SendUpdate() {
	t.mu.Lock()
	t.updatedAt = time.Now()
	t.mu.Unlock()

	t.updateCh <- UpdateMessage{
		t.id,
		t.status,
	}
}

// ID returns the task identifier
func (t *task) ID() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.id

}

// Status returns the latest task status
func (t *task) Status() models.AttackStatus {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.status
}

// Params returns a the configured attack params
func (t *task) Params() models.AttackParams {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.params
}

// CreatedAt returns the created at timestamp
func (t *task) CreatedAt() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.createdAt
}

// UpdatedAt returns the created at timestamp
func (t *task) UpdatedAt() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.updatedAt
}

// Result returns the result as a io.Reader
func (t *task) Result() io.Reader {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.result
}

func run(t *task, fn AttackFunc) {
	buf, err := fn(t.id, t.params, t.quit)
	if err != nil {
		_ = t.Fail()
	}

	// Attack was canceled
	if buf == nil {
		return
	}

	// Mark attack as completed
	err = t.Complete(buf)
	if err != nil {
		log.WithError(err).Error("Failed to Complete")
		_ = t.Fail()
	}
}

func (t *task) log(fields map[string]interface{}) *log.Entry {
	l := log.WithField("component", "task")

	l = l.WithFields(log.Fields{
		"ID":     t.ID(),
		"Status": t.Status(),
	})

	if fields != nil {
		l = l.WithFields(fields)
	}
	return l
}

func attackDetailFromTask(t ITaskGetter) models.AttackDetails {
	details := models.AttackDetails{
		AttackInfo: models.AttackInfo{
			ID:        t.ID(),
			Status:    t.Status(),
			Params:    t.Params(),
			CreatedAt: t.CreatedAt().Format(time.RFC1123),
			UpdatedAt: t.UpdatedAt().Format(time.RFC1123),
		},
	}

	if t.Status() == models.AttackResponseStatusCompleted {
		result := t.Result()
		buf, _ := ioutil.ReadAll(result)
		details.Result = buf
	}

	return details
}
