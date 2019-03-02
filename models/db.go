package models

import (
	"fmt"
	"sync"
)

// IAttackStore captures all methods related to storing and retrieving attack details
type IAttackStore interface {
	// Add item by its ID string
	Add(AttackDetails) error

	// GetAll items
	GetAll(filters FilterParams) []AttackDetails
	// GetByID gets an item by its ID
	GetByID(string) (AttackDetails, error)

	// Update multiple fields in an item
	Update(string, AttackDetails) error
	// Set a member field
	//Set(string, string, interface{}) error

	// Delete an item by ID
	Delete(string) error
}

var mu sync.RWMutex

// TaskMap is a map of attack ID's to their AttackDetails
type TaskMap map[string]AttackDetails

// NewTaskMap constructs a new instance of TaskMap
func NewTaskMap() TaskMap {
	return make(TaskMap)
}

// Add attack details by ID to store
func (tm TaskMap) Add(attack AttackDetails) error {
	mu.Lock()
	defer mu.Unlock()

	tm[attack.ID] = attack

	return nil
}

// GetAll attacks and details from store
func (tm TaskMap) GetAll(filterParams FilterParams) []AttackDetails {
	mu.RLock()
	defer mu.RUnlock()

	filters := createFilterChain(filterParams)
	attacks := make([]AttackDetails, 0)
	for _, attack := range tm {
		for _, filter := range filters {
			if !filter(attack) {
				goto skip
			}
		}
		attacks = append(attacks, attack)
	skip:
	}

	return attacks
}

// GetByID returns an attack detail by ID
func (tm TaskMap) GetByID(id string) (AttackDetails, error) {
	mu.RLock()
	defer mu.RUnlock()

	attack, ok := tm[id]
	if !ok {
		return AttackDetails{}, fmt.Errorf("attack with id %s not found", id)
	}

	return attack, nil
}

// Update an attack detail in the store
func (tm TaskMap) Update(id string, attack AttackDetails) error {
	if attack.ID != id {
		return fmt.Errorf("update ID %s and attack ID %s do not match", id, attack.ID)
	}

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

// Delete an attack by ID from the store
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

func createFilterChain(params FilterParams) []Filter {
	filters := make([]Filter, 0)
	if status, ok := params["status"]; ok {
		filters = append(
			filters,
			StatusFilter(status.(string)),
		)
	}
	return filters
}
