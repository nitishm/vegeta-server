package dispatcher

import (
	"fmt"
	"sync"
	"vegeta-server/internal"
)

type IDispatcher interface {
	Dispatch(internal.ITask)
	Get(string) (internal.ITask, error)
	List() []internal.ITask
	Cancel(string, bool) (internal.ITask, error)
}

type dispatcher struct {
	mu       *sync.RWMutex
	tasks    map[string]internal.ITask
	attackFn internal.AttackFunc
}

func NewDispatcher(fn internal.AttackFunc) *dispatcher {
	return &dispatcher{
		&sync.RWMutex{},
		make(map[string]internal.ITask),
		fn,
	}

}

func (d *dispatcher) Dispatch(task internal.ITask) {
	d.mu.Lock()
	defer d.mu.Unlock()
	// Track the task
	d.tasks[task.ID()] = task

	// Run the task
	task.Run(d.attackFn)

	return
}

func (d *dispatcher) Get(id string) (internal.ITask, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	t, ok := d.tasks[id]
	if !ok {
		err := fmt.Errorf("cannot find task with id %s", id)
		return nil, err
	}
	return t, nil
}

func (d *dispatcher) List() []internal.ITask {
	d.mu.RLock()
	defer d.mu.RUnlock()
	tasks := make([]internal.ITask, 0)
	for _, task := range d.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (d *dispatcher) Cancel(id string, cancel bool) (internal.ITask, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	t, ok := d.tasks[id]
	if !ok {
		return nil, fmt.Errorf("cannot find task with id %s", id)
	}

	if cancel {
		err := t.Cancel()
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}
