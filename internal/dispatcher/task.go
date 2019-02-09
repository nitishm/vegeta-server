package dispatcher

import (
	"bytes"
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tsenart/vegeta/lib"
	"sync"
	"vegeta-server/internal/models"
)

type ITask interface {
	ID() string
	Status() models.AttackStatus
	Params() models.AttackParams

	Run(attackFunc) error
	Complete() error
	Cancel() error
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
	mu     *sync.RWMutex
	ctx    attackContext
	id     string
	params models.AttackParams
	status models.AttackStatus
}

func NewTask(params models.AttackParams) *task {
	id := uuid.NewV4().String()
	t := &task{
		&sync.RWMutex{},
		newAttackContext(),
		id,
		params,
		models.AttackResponseStatusScheduled,
	}
	t.log(nil).Debug("creating new task")
	return t
}

func (t *task) Run(fn attackFunc) error {
	if t.status != models.AttackResponseStatusScheduled {
		return fmt.Errorf("cannot run task %s with status %s", t.id, t.status)
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.status = models.AttackResponseStatusRunning
	t.log(nil).Debug("running")

	go run(t, fn)

	return nil
}

func (t *task) Complete() error {
	if t.status != models.AttackResponseStatusRunning {
		return fmt.Errorf("cannot mark completed for task %s with status %s", t.id, t.status)
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.status = models.AttackResponseStatusCompleted
	t.log(nil).Debug("completed")

	return nil
}

func (t *task) Cancel() error {
	if t.status == models.AttackResponseStatusCompleted || t.status == models.AttackResponseStatusFailed {
		return fmt.Errorf("cannot cancel task %s with status %s", t.id, t.status)
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.ctx.cancelFn()
	t.status = models.AttackResponseStatusCanceled
	t.log(nil).Debug("canceled")

	return nil
}

func (t *task) Fail() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.status = models.AttackResponseStatusFailed
	t.log(nil).Error("failed")
	return nil
}

func (t *task) ID() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.id
}

func (t *task) Status() models.AttackStatus {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.status
}

func (t *task) Params() models.AttackParams {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.params
}

func run(t *task, fn attackFunc) error {
	opts, err := models.NewAttackOptsFromAttackParams(t.id, t.params)
	if err != nil {
		return err
	}

	result := fn(opts)
	if result == nil {
		return fmt.Errorf("empty channel returned")
	}

	buf := bytes.NewBuffer(nil)
	enc := vegeta.NewEncoder(buf)
loop:
	for {
		select {
		case r, ok := <-result:
			if !ok {
				break loop
			}
			if err := enc.Encode(r); err != nil {
				_ = t.Fail()
			}
		case <-t.ctx.Done():
			t.log(nil).Warn("task was canceled", t.id)
			return nil
		}
	}

	// Write results to reporter channel
	//t.resCh <- &Result{
	//	entry.uuid,
	//	buf,
	//}

	// Mark attack as completed
	err = t.Complete()
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
