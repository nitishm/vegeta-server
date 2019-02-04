package attacker

import (
	"fmt"
	"github.com/satori/go.uuid"
	"vegeta-server/internal/app/server/models"
	"context"
)


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
	ctx    attackContext
	id     string
	params models.Attack
	status models.AttackStatus
}

func newTask(params models.Attack) *task {
	id := uuid.NewV4().String()
	return &task{
		newAttackContext(),
		id,
		params,
		models.AttackResponseStatusScheduled,
	}
}

func (t *task) run(fn AttackFunc) error {
	if t.status != models.AttackResponseStatusScheduled {
		return fmt.Errorf("cannot run task %s with status %s", t.id, t.status)
	}

	t.status = models.AttackResponseStatusRunning

	return run(fn, t.params)
}

func (t *task) cancel() error {
	if t.status == models.AttackResponseStatusCompleted || t.status == models.AttackResponseStatusFailed {
		return fmt.Errorf("cannot cancel task %s with status %s", t.id, t.status)
	}

	t.ctx.cancelFn()

	t.status = models.AttackResponseStatusCanceled

	return nil
}

func (t *task) fail() error {
	t.status = models.AttackResponseStatusFailed
	return nil
}

func run(fn AttackFunc, params models.Attack) error {
	opts := AttackOpts{}
	fn(&opts)
	return nil
}
