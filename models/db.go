package models

import (
	"fmt"
	"sync"
)

type IAttackStore interface {
	// Add item by its ID string
	Add(AttackDetails) error

	// GetAll items
	GetAll() []AttackDetails
	// GetById gets an item by its ID
	GetById(string) (AttackDetails, error)

	// Update multiple fields in an item
	Update(string, AttackDetails) error
	// Set a member field
	//Set(string, string, interface{}) error

	// Delete an item by ID
	Delete(string) error
}

var mu sync.RWMutex

type TaskMap map[string]AttackDetails

func NewTaskMap() TaskMap {
	return make(TaskMap)
}

func (tm TaskMap) Add(attack AttackDetails) error {
	mu.Lock()
	defer mu.Unlock()

	tm[attack.ID] = attack

	return nil
}

func (tm TaskMap) GetAll() []AttackDetails {
	mu.RLock()
	defer mu.RUnlock()

	attacks := make([]AttackDetails, 0)
	for _, attack := range tm {
		attacks = append(attacks, attack)
	}

	return attacks
}

func (tm TaskMap) GetById(id string) (AttackDetails, error) {
	mu.RLock()
	defer mu.RUnlock()

	attack, ok := tm[id]
	if !ok {
		return AttackDetails{}, fmt.Errorf("attack with id %s not found", id)
	}

	return attack, nil
}

func (tm TaskMap) Update(id string, attack AttackDetails) error {
	mu.RLock()
	_, ok := tm[id]
	mu.RUnlock()
	if !ok {
		return fmt.Errorf("attack with id %s not found", id)
	}

	mu.Lock()
	tm[id] = attack
	mu.Unlock()

	return nil
}

func (tm TaskMap) Delete(id string) error {
	mu.RLock()
	_, ok := tm[id]
	mu.RUnlock()
	if !ok {
		return fmt.Errorf("attack with id %s not found", id)
	}

	mu.Lock()
	defer mu.Unlock()

	delete(tm, id)

	return nil
}
