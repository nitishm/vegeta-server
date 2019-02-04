package attacker

import "sync"

type IDispatcher interface {
	Dispatch(*task)
}

type dispatcher struct {
	mu *sync.RWMutex
	tasks map[string]*task
	attackFn AttackFunc
}

func NewDispatcher(fn AttackFunc) *dispatcher {
	return &dispatcher{
		&sync.RWMutex{},
		make(map[string]*task),
		fn,
	}

}

func (d *dispatcher) Dispatch(task *task) {
	d.mu.Lock()
	defer d.mu.Unlock()
	// Track the task
	d.tasks[task.id] = task

	// Schedule the task
	task.run(d.attackFn)

	return
}

