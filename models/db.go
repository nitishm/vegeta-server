package models

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gomodule/redigo/redis"
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

// Redis stores all Attack/Report information in a redis database
type Redis struct {
	connFn func() redis.Conn
}

func NewRedis(f func() redis.Conn) Redis {
	return Redis{
		f,
	}
}

func (r Redis) Add(attack AttackDetails) error {
	conn := r.connFn()
	v, err := json.Marshal(attack)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", attack.ID, v)
	if err != nil {
		return err
	}
	return nil
}

func (r Redis) GetAll(filterParams FilterParams) []AttackDetails {
	var attacks []AttackDetails
	filters := createFilterChain(filterParams)
	conn := r.connFn()
	res, err := conn.Do("KEYS", "*")
	if err != nil {
		return nil
	}
	if attackIDs, ok := res.([]interface{}); !ok {
		return nil
	} else {
		for _, attackID := range attackIDs {
			var attack AttackDetails

			res, err := conn.Do("GET", attackID)
			if err != nil {
				return nil
			}

			if err := json.Unmarshal(res.([]byte), &attack); err != nil {
				return nil
			} else {
				for _, filter := range filters {
					if !filter(attack) {
						goto skip
					}
				}
				attacks = append(attacks, attack)
			skip:
			}
		}
	}

	return attacks
}

func (r Redis) GetByID(id string) (AttackDetails, error) {
	var attack AttackDetails
	conn := r.connFn()
	res, err := conn.Do("GET", id)
	if err != nil {
		return attack, err
	}

	err = json.Unmarshal(res.([]byte), &attack)
	if err != nil {
		return attack, err
	}

	return attack, err
}

func (r Redis) Update(id string, attack AttackDetails) error {
	if attack.ID != id {
		return fmt.Errorf("update ID %s and attack ID %s do not match", id, attack.ID)
	}
	return r.Add(attack)
}

func (r Redis) Delete(id string) error {
	conn := r.connFn()
	_, err := conn.Do("DEL", id)
	if err != nil {
		return err
	}
	return nil
}

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
	if createdBefore, ok := params["created_before"]; ok {
		filters = append(
			filters,
			CreationBeforeFilter(createdBefore.(string)),
		)
	}
	if createdAfter, ok := params["created_after"]; ok {
		filters = append(
			filters,
			CreationAfterFilter(createdAfter.(string)),
		)
	}
	return filters
}
