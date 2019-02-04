package attacker

import (
	"bytes"
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tsenart/vegeta/lib"
	"net/http"
	"strings"
	"time"
	"vegeta-server/internal/app/server/models"
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

	go run(t, fn)

	return nil
}

func (t *task) complete() error {
	if t.status != models.AttackResponseStatusRunning {
		return fmt.Errorf("cannot mark completed for task %s with status %s", t.id, t.status)
	}

	t.status = models.AttackResponseStatusCompleted

	return nil
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

func (t *task) getID() string {
	return t.id
}

func (t *task) getStatus() models.AttackStatus {
	return t.status
}

func run(t *task, fn AttackFunc) error {
	opts, err := attackOptsFromModel(t.id, t.params)
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
				err := t.fail()
				if err != nil {
					log.Fatal(err)
				}
			}
		case <-t.ctx.Done():
			log.Warnf("Attack %s was canceled", t.id)
			return nil
		}
	}

	// Write results to reporter channel
	//t.resCh <- &Result{
	//	entry.uuid,
	//	buf,
	//}

	// Mark attack as completed
	err = t.complete()
	if err != nil {
		log.WithError(err).Error("Failed to Complete")
		_ = t.fail()
	}

	return nil
}

func attackOptsFromModel(id string, params models.Attack) (*AttackOpts, error) {
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
