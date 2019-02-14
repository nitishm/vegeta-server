package dispatcher

import (
	"bytes"
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	log "github.com/sirupsen/logrus"
	vlib "github.com/tsenart/vegeta/lib"

	"io"
	"io/ioutil"
	"sync"
	"vegeta-server/models"
	"vegeta-server/pkg/vegeta"
)

// ITask defines an interface for attack tasks
type ITask interface {
	// ID returns the attack task ID
	ID() string
	// Status returns the attack task status
	Status() models.AttackStatus
	// Params returns the attack task params
	Params() models.AttackParams

	// Run the attack using the configured attack function.
	Run(vegeta.AttackFunc) error
	// Complete changes task status to completed
	Complete(io.Reader) error
	// Cancel changes task status to canceled
	Cancel() error
	// Fail changes task status to failed
	Fail() error
}

type attackContext struct {
	context.Context
	cancelFn context.CancelFunc
}

func newAttackContext() attackContext {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return attackContext{
		ctx,
		cancel,
	}
}

type task struct {
	mu  *sync.RWMutex
	ctx attackContext

	id     string
	params models.AttackParams
	status models.AttackStatus

	updateCh chan models.AttackDetails
}

// NewTask returns a new instance of a task object
func NewTask(updateCh chan models.AttackDetails, params models.AttackParams) *task { //nolint: golint
	id := uuid.NewV4().String()
	t := &task{
		&sync.RWMutex{},
		newAttackContext(),
		id,
		params,
		models.AttackResponseStatusScheduled,
		updateCh,
	}
	t.log(nil).Debug("creating new task")
	return t
}

func (t *task) update(status models.AttackStatus, result io.Reader) {
	t.status = status
	details := models.AttackDetails{
		AttackInfo: models.AttackInfo{
			ID:     t.id,
			Status: t.status,
			Params: t.params,
		},
	}

	if result != nil {
		buf, _ := ioutil.ReadAll(result)
		details.Result = buf
	}

	t.updateCh <- details
}

// Run an attack task using the passed in attack function
func (t *task) Run(fn vegeta.AttackFunc) error {
	if t.status != models.AttackResponseStatusScheduled {
		return fmt.Errorf("cannot run task %s with status %s", t.id, t.status)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.update(models.AttackResponseStatusRunning, nil)

	t.log(nil).Debug("running")

	go run(t, fn) //nolint: errcheck

	return nil
}

// Complete marks a task as completed
func (t *task) Complete(result io.Reader) error {
	if t.status != models.AttackResponseStatusRunning {
		return fmt.Errorf("cannot mark completed for task %s with status %s", t.id, t.status)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.update(models.AttackResponseStatusCompleted, result)

	t.log(nil).Debug("completed")

	return nil
}

// Cancel invokes the context cancel and marks a task as canceled
func (t *task) Cancel() error {
	if t.status == models.AttackResponseStatusCompleted || t.status == models.AttackResponseStatusFailed {
		return fmt.Errorf("cannot cancel task %s with status %s", t.id, t.status)
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.ctx.cancelFn()

	t.update(models.AttackResponseStatusCanceled, nil)

	t.log(nil).Debug("canceled")

	return nil
}

// Fail marks a task as failed
func (t *task) Fail() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.update(models.AttackResponseStatusFailed, nil)

	t.log(nil).Error("failed")
	return nil
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

// Params returns a the confgured attack params
func (t *task) Params() models.AttackParams {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.params
}

// TODO: Remove dependency on vegeta lib. Move functionality to pkg/vegeta package.
func run(t *task, fn vegeta.AttackFunc) error {
	opts, err := vegeta.NewAttackOptsFromAttackParams(t.id, t.params)
	if err != nil {
		return err
	}

	result := fn(opts)
	if result == nil {
		return fmt.Errorf("empty channel returned")
	}

	buf := bytes.NewBuffer(nil)
	enc := vlib.NewEncoder(buf)
loop:
	for {
		select {
		case r, ok := <-result:
			if !ok {
				break loop
			}
			if err = enc.Encode(r); err != nil {
				_ = t.Fail()
			}
		case <-t.ctx.Done():
			t.log(nil).Warn("task was canceled", t.id)
			return nil
		}
	}

	// Mark attack as completed
	err = t.Complete(buf)
	if err != nil {
		log.WithError(err).Error("Failed to Complete")
		_ = t.Fail()
	}

	return nil
}

func (t *task) log(fields map[string]interface{}) *log.Entry {
	l := log.WithField("component", "task")

	l = l.WithFields(log.Fields{
		"ID":     t.id,
		"Status": t.status,
	})

	if fields != nil {
		l = l.WithFields(fields)
	}
	return l
}
