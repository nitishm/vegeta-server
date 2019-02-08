package internal

import (
	"bytes"
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tsenart/vegeta/lib"
	"net/http"
	"strings"
	"sync"
	"time"
	"vegeta-server/internal/models"
)

type AttackFunc func(*AttackOpts) <-chan *vegeta.Result

func DefaultAttackFn(opts *AttackOpts) <-chan *vegeta.Result {
	atk := vegeta.NewAttacker(
		vegeta.Redirects(opts.Redirects),
		vegeta.Timeout(opts.Timeout),
		vegeta.Workers(opts.Workers),
		vegeta.KeepAlive(opts.Keepalive),
		vegeta.Connections(opts.Connections),
		vegeta.HTTP2(opts.HTTP2),
		vegeta.H2C(opts.H2c),
		vegeta.MaxBody(opts.MaxBody),
	)
	tr := vegeta.NewStaticTargeter(opts.Target)
	return atk.Attack(tr, opts.Rate, opts.Duration, opts.Name)
}

type ITask interface {
	ID() string
	Status() models.AttackStatus
	Params() models.AttackParams

	Run(AttackFunc) error
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
	return &task{
		&sync.RWMutex{},
		newAttackContext(),
		id,
		params,
		models.AttackResponseStatusScheduled,
	}
}

func (t *task) Run(fn AttackFunc) error {
	if t.status != models.AttackResponseStatusScheduled {
		return fmt.Errorf("cannot run task %s with status %s", t.id, t.status)
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.status = models.AttackResponseStatusRunning

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

	return nil
}

func (t *task) Fail() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.status = models.AttackResponseStatusFailed
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

func run(t *task, fn AttackFunc) error {

	opts, err := attackOptsFromModel(t.ID(), t.Params())
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
				err := t.Fail()
				if err != nil {
					log.Fatal(err)
				}
			}
		case <-t.ctx.Done():
			log.Warnf("AttackParams %s was canceled", t.id)
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

func attackOptsFromModel(id string, params models.AttackParams) (*AttackOpts, error) {
	rate := vegeta.Rate{Freq: int(params.Rate), Per: time.Second}

	// Set Duration
	dur, err := time.ParseDuration(params.Duration)
	if err != nil {
		return nil, err
	}

	// Set timeout
	timeout, _ := time.ParseDuration(params.Timeout)

	// Set target headers
	var hdr http.Header
	for _, h := range params.Headers {
		hdr.Add(h.Key, h.Value)
	}

	// Set resolvers
	resolvers := strings.Split(params.Resolvers, ",")

	// TODO: Set Local Address

	// TODO: Set TLS configuration

	// TODO: Set target body

	// Set Target
	tgt := vegeta.Target{
		Method: params.Target.Method,
		URL:    params.Target.URL,
		Header: hdr,
	}

	opts := &AttackOpts{
		Name:      id,
		Target:    tgt,
		Insecure:  params.Insecure,
		Duration:  dur,
		Timeout:   timeout,
		Rate:      rate,
		Redirects: int(params.Redirects),
		MaxBody:   params.MaxBody,
		Keepalive: params.Keepalive,
		Resolvers: resolvers,
	}
	opts.HTTP2 = params.Http2
	opts.H2c = params.H2c
	opts.Workers = uint64(params.Workers)

	return opts, nil
}
