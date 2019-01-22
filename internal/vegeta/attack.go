package vegeta

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"vegeta-server/models"
	"vegeta-server/pkg/vegeta"

	uuid "github.com/satori/go.uuid"
	vegetalib "github.com/tsenart/vegeta/lib"
)

type attackContext struct {
	context.Context
	cancelFn context.CancelFunc
}
type attackEntry struct {
	ctx attackContext

	uuid   string
	status string

	resCh chan<- *Result
}

// Status returns the attack status
func (ae *attackEntry) Status() string {
	return ae.status
}

// Schedule an attack and once scheduled invoke the Run method
func (ae *attackEntry) Schedule(params interface{}) error {
	if ae.status == models.AttackResponseStatusCanceled || ae.status == models.AttackResponseStatusRunning || ae.status == models.AttackResponseStatusCompleted { //nolint:lll
		return fmt.Errorf("cannot schedule attack %s with status %v", ae.uuid, ae.status)
	}

	log.WithField("UUID", ae.uuid).Info("Scheduled")
	ae.status = models.AttackResponseStatusScheduled

	// TODO: Create a scheduler that is smart enough to schedule
	//attacks when there are enough resources available
	_ = ae.Run(params)

	return nil
}

// Run an attack against the target
func (ae *attackEntry) Run(params interface{}) error {
	if ae.status != models.AttackResponseStatusScheduled {
		return fmt.Errorf("cannot run attack %s with status %v", ae.uuid, ae.status)
	}

	log.WithField("UUID", ae.uuid).Info("Running")
	ae.status = models.AttackResponseStatusRunning

	attackOpts := convertParamsToAttackOpts(ae.uuid, params)

	go attackHandler(ae, attackOpts, vegeta.DefaultAttackFunc)

	return nil
}

// Complete marks the attack as completed
func (ae *attackEntry) Complete() error {
	if ae.status != models.AttackResponseStatusRunning {
		return fmt.Errorf("cannot mark attack %s completed when in status %v", ae.uuid, ae.status)
	}
	log.WithField("UUID", ae.uuid).Info("Completed")
	ae.status = models.AttackResponseStatusCompleted
	return nil
}

// Cancel an attack and update the status as canceled
func (ae *attackEntry) Cancel() error {
	if ae.status == models.AttackResponseStatusCompleted || ae.status == models.AttackResponseStatusFailed {
		return fmt.Errorf("cannot cancel attack %s  with status %v", ae.uuid, ae.status)
	}
	// Cancel the attack context
	ae.ctx.cancelFn()

	log.WithField("UUID", ae.uuid).Info("Canceled")
	ae.status = models.AttackResponseStatusCanceled
	return nil
}

// Fail marks the attack status as failed
func (ae *attackEntry) Fail() error {
	_ = ae.Cancel()
	log.WithField("UUID", ae.uuid).Info("Failed")
	ae.status = models.AttackResponseStatusFailed
	return nil
}

// UUID returns the ID of the attack
func (ae *attackEntry) UUID() string {
	return ae.uuid
}

type attackCmd struct {
	uuid   string
	params interface{}
}

// AttackIntf is the Attacker interface
type AttackIntf interface {
	Schedule(interface{}) string
	Status(string) (string, error)
	Cancel(string, bool) (string, error)

	Exists(string) bool
	List() map[string]string
}

// Attacker implements the AttackIntf
type Attacker struct {
	ch        chan attackCmd
	lock      sync.RWMutex
	scheduler map[string]*attackEntry
	quit      chan struct{}

	Results chan *Result
}

// NewAttacker returns an instance of a new attacker.
func NewAttacker() *Attacker {
	at := &Attacker{
		scheduler: make(map[string]*attackEntry),
		ch:        make(chan attackCmd),
		quit:      make(chan struct{}),
		Results:   make(chan *Result),
	}

	go at.startAttackHandler()

	return at
}

// Schedule adds the attack command to the scheduler
func (at *Attacker) Schedule(params interface{}) string {
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

// Status returns the status for an attack by its ID
func (at *Attacker) Status(uuid string) (string, error) {
	at.lock.RLock()
	defer at.lock.RUnlock()
	entry, ok := at.scheduler[uuid]
	if !ok {
		return "", fmt.Errorf("attack reference %s not found", uuid)
	}
	return entry.Status(), nil
}

// Cancel an attack by its ID
func (at *Attacker) Cancel(uuid string, cancel bool) (string, error) {
	at.lock.Lock()
	defer at.lock.Unlock()
	entry, ok := at.scheduler[uuid]
	if !ok {
		return "", fmt.Errorf("attack reference %s not found", uuid)
	}
	err := entry.Cancel()
	return entry.Status(), err
}

func (at *Attacker) Exists(uuid string) bool {
	_, ok := at.scheduler[uuid]
	return ok
}

func (at *Attacker) List() map[string]string {
	at.lock.RLock()
	defer at.lock.RUnlock()

	m := make(map[string]string)
	for uuid, entry := range at.scheduler {
		m[uuid] = entry.Status()
	}
	return m
}

func (at *Attacker) startAttackHandler() {
	fmt.Println("Starting Attack Handlers")
	for {
		select {
		case cmd := <-at.ch:
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)

			entry := &attackEntry{
				ctx:   attackContext{ctx, cancel},
				uuid:  cmd.uuid,
				resCh: at.Results,
			}

			// Create an entry in the attack database
			at.lock.Lock()
			at.scheduler[cmd.uuid] = entry
			// Mark attack as Scheduled
			err := entry.Schedule(cmd.params)
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

func convertParamsToAttackOpts(name string, params interface{}) *vegeta.AttackOpts {
	switch p := params.(type) {
	case *models.Attack:
		// Set Rate
		rate := vegetalib.Rate{Freq: int(*p.Rate), Per: time.Second}

		// Set Duration
		dur := time.Duration(*p.Duration)

		// Set timeout
		timeout, _ := time.ParseDuration(p.Timeout)

		// Set target headers
		var hdr http.Header
		for _, h := range p.Headers {
			hdr.Add(h.Key, h.Value)
		}

		// Set resolvers
		resolvers := strings.Split(p.Resolvers, ",")

		// TODO: Set Local Address

		// TODO: Set TLS configuration

		// TODO: Set target body

		// Set Target
		tgt := vegetalib.Target{
			Method: *p.Target.Method,
			URL:    p.Target.URL,
			Header: hdr,
		}

		opts := &vegeta.AttackOpts{
			Name:      name,
			Target:    tgt,
			Insecure:  p.Insecure,
			Duration:  dur,
			Timeout:   timeout,
			Rate:      rate,
			Redirects: int(p.Redirects),
			MaxBody:   p.MaxBody,
			Keepalive: p.Keepalive,
			Resolvers: resolvers,
		}
		if p.Http2 != nil {
			opts.HTTP2 = *p.Http2
		}
		if p.H2c != nil {
			opts.H2c = *p.H2c
		}
		if p.Workers != nil {
			opts.Workers = uint64(*p.Workers)
		}
		return opts
	default:
		return nil
	}
}

func attackHandler(entry *attackEntry, opts *vegeta.AttackOpts, fn vegeta.AttackFunc) {
	result := fn(opts)
	buf := bytes.NewBuffer(nil)
	enc := vegetalib.NewEncoder(buf)
loop:
	for {
		select {
		case r, ok := <-result:
			if !ok {
				break loop
			}
			if err := enc.Encode(r); err != nil {
				err := entry.Fail()
				if err != nil {
					log.Fatal(err)
				}
			}
		case <-entry.ctx.Done():
			log.Warnf("Attack %s was canceled", entry.uuid)
			break loop
		}
	}

	// Write results to reporter channel
	entry.resCh <- &Result{
		entry.uuid,
		buf,
	}

	// Mark attack as completed
	err := entry.Complete()
	if err != nil {
		log.WithError(err).Error("Failed to Complete")
		_ = entry.Fail()
	}
}
