package attacker

import (
	"fmt"
	"sync"
)

type IDispatcher interface {
	Dispatch(*task)
	Get(string) (*task, error)
	List() []*task
	Cancel(string, bool) (*task, error)
}

type dispatcher struct {
	mu       *sync.RWMutex
	tasks    map[string]*task
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

	// Run the task
	task.run(d.attackFn)

	return
}

func (d *dispatcher) Get(id string) (*task, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	t, ok := d.tasks[id]
	if !ok {
		err := fmt.Errorf("cannot find task with id %s", id)
		return nil, err
	}
	return t, nil
}

func (d *dispatcher) List() []*task {
	d.mu.RLock()
	defer d.mu.RUnlock()
	tasks := make([]*task, 0)
	for _, task := range d.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (d *dispatcher) Cancel(id string, cancel bool) (*task, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	t, ok := d.tasks[id]
	if !ok {
		return nil, fmt.Errorf("cannot find task with id %s", id)
	}

	if cancel {
		err := t.cancel()
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}
