package vegeta

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"vegeta-server/models"

	"github.com/satori/go.uuid"
)

type AttackParams *models.Attack

type attackContext struct {
	context.Context
	cancelFn context.CancelFunc
}
type attackEntry struct {
	ctx    attackContext
	uuid   string
	status string
}

func (ae *attackEntry) Status() string {
	return ae.status
}

func (ae *attackEntry) Schedule() error {
	if ae.status == models.AttackResponseStatusCanceled || ae.status == models.AttackResponseStatusRunning || ae.status == models.AttackResponseStatusCompleted {
		return fmt.Errorf("Cannot schedule attack %s with status %v", ae.uuid, ae.status)
	}
	log.WithField("UUID", ae.uuid).Info("Scheduled")

	time.AfterFunc(time.Second*5, func() {
		_ = ae.Run()
	})

	ae.status = models.AttackResponseStatusScheduled
	return nil
}

func (ae *attackEntry) Run() error {
	if ae.status != models.AttackResponseStatusScheduled {
		return fmt.Errorf("Cannot run attack %s with status %v", ae.uuid, ae.status)
	}

	log.WithField("UUID", ae.uuid).Info("Running")
	ae.status = models.AttackResponseStatusRunning

	//TODO: Invoke the attack function
	attackFunc := func(entry *attackEntry) {
		select {
		case <-time.After(60 * time.Second):
			err := entry.Complete()
			if err != nil {
				log.WithError(err).Error("Failed to Complete")
				_ = entry.Fail()
			}
		case <-entry.ctx.Done():
			log.Warnf("Attack %s was canceled", entry.uuid)
			break
		}
	}

	go attackFunc(ae)
	return nil
}

func (ae *attackEntry) Complete() error {
	if ae.status != models.AttackResponseStatusRunning {
		return fmt.Errorf("Cannot mark attack %s completed when in status %v", ae.uuid, ae.status)
	}
	log.WithField("UUID", ae.uuid).Info("Completed")
	ae.status = models.AttackResponseStatusCompleted
	return nil
}

func (ae *attackEntry) Cancel() error {
	if ae.status == models.AttackResponseStatusCompleted || ae.status == models.AttackResponseStatusFailed || ae.status == models.AttackResponseStatusCanceled {
		return fmt.Errorf("Cannot cancel attack %s  with status %v", ae.uuid, ae.status)
	}
	// Cancel the attack context
	ae.ctx.cancelFn()

	log.WithField("UUID", ae.uuid).Info("Canceled")
	ae.status = models.AttackResponseStatusCanceled
	return nil
}

func (ae *attackEntry) Fail() error {
	_ = ae.Cancel()
	log.WithField("UUID", ae.uuid).Info("Failed")
	ae.status = models.AttackResponseStatusFailed
	return nil
}

func (ae *attackEntry) UUID() string {
	return ae.uuid
}

type attackCmd struct {
	uuid   string
	params interface{}
}

type AttackIntf interface {
	Schedule(AttackParams) string
	Status(string) (string, error)
}

type Attacker struct {
	ch        chan attackCmd
	lock      sync.RWMutex
	scheduler map[string]*attackEntry
	quit      chan struct{}
}

func NewAttacker() *Attacker {
	at := &Attacker{
		scheduler: make(map[string]*attackEntry),
		ch:        make(chan attackCmd),
		quit:      make(chan struct{}),
	}

	go at.startAttackHandler()

	return at
}

func (at *Attacker) Schedule(params AttackParams) string {
	// Generate a uuid for the attack
	id := uuid.NewV4().String()

	// Submit attack command params to the central attacker
	at.ch <- attackCmd{
		uuid:   id,
		params: params,
	}

	// Return the uuid to the user to check for status and report
	return id
}

func (at *Attacker) Status(uuid string) (string, error) {
	at.lock.RLock()
	defer at.lock.RUnlock()
	if entry, ok := at.scheduler[uuid]; !ok {
		return "", fmt.Errorf("attack reference %s not found", uuid)
	} else {
		return entry.Status(), nil
	}
}

func (at *Attacker) Cancel(uuid string, cancel bool) (string, error) {
	at.lock.Lock()
	defer at.lock.Unlock()
	if entry, ok := at.scheduler[uuid]; !ok {
		return "", fmt.Errorf("attack reference %s not found", uuid)
	} else {
		err := entry.Cancel()
		return entry.Status(), err
	}
}
func (at *Attacker) startAttackHandler() {
	fmt.Println("Starting Attack Handlers")
	for {
		select {
		case cmd := <-at.ch:
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)

			// Create an entry in the attack database
			at.lock.Lock()
			entry := &attackEntry{
				ctx:  attackContext{ctx, cancel},
				uuid: cmd.uuid,
			}

			at.scheduler[cmd.uuid] = entry
			// Mark attack as Scheduled
			err := entry.Schedule()
			if err != nil {
				log.WithError(err).Error("Failed to Schedule")
				_ = entry.Fail()
				continue
			}
			at.lock.Unlock()
		case <-at.quit:
			break
		}
	}
}
