package dispatcher

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
	"vegeta-server/models"
	smocks "vegeta-server/models/mocks"

	"github.com/stretchr/testify/mock"
)

func TestNewDispatcher(t *testing.T) {
	type args struct {
		db models.IAttackStore
		fn AttackFunc
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "OK",
			args: args{
				db: &smocks.IAttackStore{},
				fn: func(s string, params models.AttackParams, i chan struct{}) (reader io.Reader, e error) {
					return strings.NewReader("hello world"), nil
				},
			},
			wantNil: false,
		},
		{
			name: "OK - defaults db",
			args: args{
				fn: func(s string, params models.AttackParams, i chan struct{}) (reader io.Reader, e error) {
					return strings.NewReader("hello world"), nil
				},
			},
			wantNil: false,
		},
		{
			name: "OK - defaults attack fn",
			args: args{
				db: &smocks.IAttackStore{},
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDispatcher(tt.args.db, tt.args.fn)
			if tt.wantNil && got != nil {
				t.Errorf("NewDispatcher() = %v, wantNit %v", got, tt.wantNil)
			}
		})
	}
}

func setupDispatcher(db models.IAttackStore) *dispatcher {
	d := &dispatcher{
		mu:    new(sync.RWMutex),
		tasks: make(map[string]ITask),
		attackFn: func(s string, params models.AttackParams, i chan struct{}) (reader io.Reader, e error) {
			return strings.NewReader("hello world"), nil
		},
		submitCh: make(chan ITask),
		updateCh: make(chan UpdateMessage),
		db:       db,
	}

	go func() {
		for {
			select {
			case <-d.submitCh:
				continue
			case <-d.updateCh:
				continue
			}
		}
	}()

	return d
}

func Test_dispatcher_Dispatch(t *testing.T) {
	type args struct {
		params models.AttackParams
	}
	tests := []struct {
		name    string
		db      func() models.IAttackStore
		args    args
		want    *models.AttackResponse
		wantErr bool
	}{
		{
			name: "OK",
			db: func() models.IAttackStore {
				db := new(smocks.IAttackStore)

				db.On("Add", mock.Anything).Return(nil)
				db.On("GetByID", mock.Anything).Return(models.AttackDetails{}, nil)

				return db
			},
			args: args{
				models.AttackParams{
					Rate: 10,
				},
			},
			want: &models.AttackResponse{},
		},
		{
			name: "Error - Task not found",
			db: func() models.IAttackStore {
				db := new(smocks.IAttackStore)

				db.On("Add", mock.Anything).Return(nil)
				db.On("GetByID", mock.Anything).Return(models.AttackDetails{}, fmt.Errorf("error"))

				return db
			},
			args: args{
				models.AttackParams{
					Rate: 10,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := setupDispatcher(tt.db())
			got, err := d.Dispatch(tt.args.params)
			if err != nil && !tt.wantErr {
				t.Fail()
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dispatcher.Dispatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dispatcher_Run(t *testing.T) {
	mockStore := &smocks.IAttackStore{}

	mockStore.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockStore.On("Add", mock.Anything).Return(nil)
	mockStore.On("GetByID", mock.Anything).Return(models.AttackDetails{}, nil)

	d := setupDispatcher(mockStore)

	quit := make(chan struct{})

	go d.Run(quit)

	resp, err := d.Dispatch(models.AttackParams{})
	if err != nil || resp == nil {
		t.Fail()
	}
}

func Test_dispatcher_Run_Error_GetByID(t *testing.T) {
	mockStore := &smocks.IAttackStore{}

	mockStore.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockStore.On("Add", mock.Anything).Return(nil)
	mockStore.On("GetByID", mock.Anything).Return(models.AttackDetails{}, fmt.Errorf("error"))

	d := setupDispatcher(mockStore)

	quit := make(chan struct{})

	go d.Run(quit)

	resp, err := d.Dispatch(models.AttackParams{})
	if err == nil || resp != nil {
		t.Fail()
	}
}

func Test_dispatcher_Cancel(t *testing.T) {
	mockStore := &smocks.IAttackStore{}

	mockStore.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockStore.On("Add", mock.Anything).Return(nil)
	mockStore.On("GetByID", mock.Anything).Return(models.AttackDetails{}, nil)

	d := NewDispatcher(mockStore, func(s string, params models.AttackParams, i chan struct{}) (reader io.Reader, e error) {
		<-i
		return strings.NewReader("hello world"), nil
	})

	quit := make(chan struct{})

	go d.Run(quit)

	resp, err := d.Dispatch(models.AttackParams{})
	if err != nil || resp == nil {
		t.Fail()
	}

	for _, task := range d.tasks {
		id := task.ID()
		err = d.Cancel(id, true)
		if err != nil {
			t.Fail()
		}
	}
}

func Test_dispatcher_Cancel_Error_completed(t *testing.T) {
	mockStore := &smocks.IAttackStore{}

	mockStore.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockStore.On("Add", mock.Anything).Return(nil)
	mockStore.On("GetByID", mock.Anything).Return(models.AttackDetails{}, nil)

	d := NewDispatcher(mockStore, func(s string, params models.AttackParams, i chan struct{}) (reader io.Reader, e error) {
		return strings.NewReader("hello world"), nil
	})

	quit := make(chan struct{})

	go d.Run(quit)

	resp, err := d.Dispatch(models.AttackParams{})
	if err != nil || resp == nil {
		t.Fail()
	}

	time.Sleep(time.Second * 2)

	for _, task := range d.tasks {
		id := task.ID()
		err := d.Cancel(id, true)
		if err == nil {
			t.Fail()
		}
	}
}

func Test_dispatcher_Cancel_Error_not_found(t *testing.T) {
	mockStore := &smocks.IAttackStore{}

	mockStore.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockStore.On("Add", mock.Anything).Return(nil)
	mockStore.On("GetByID", mock.Anything).Return(models.AttackDetails{}, nil)

	d := NewDispatcher(mockStore, func(s string, params models.AttackParams, i chan struct{}) (reader io.Reader, e error) {
		return nil, nil
	})

	quit := make(chan struct{})

	go d.Run(quit)

	err := d.Cancel("123", true)
	if err == nil {
		t.Fail()
	}
}

func Test_dispatcher_Get(t *testing.T) {
	mockStore := &smocks.IAttackStore{}

	mockStore.On("GetByID", mock.Anything).Return(models.AttackDetails{}, nil)

	d := setupDispatcher(mockStore)

	got, err := d.Get("123")
	if err != nil || got == nil {
		t.Fail()
	}
}

func Test_dispatcher_Get_Error_GetByID(t *testing.T) {
	mockStore := &smocks.IAttackStore{}

	mockStore.On("GetByID", mock.Anything).Return(models.AttackDetails{}, fmt.Errorf("error"))

	d := setupDispatcher(mockStore)

	got, err := d.Get("123")
	if err == nil || got != nil {
		t.Fail()
	}
}

func Test_dispatcher_List(t *testing.T) {
	mockStore := &smocks.IAttackStore{}

	mockStore.On("GetAll").Return([]models.AttackDetails{{}})

	d := setupDispatcher(mockStore)

	got := d.List()
	if len(got) == 0 {
		t.Fail()
	}
}

func Test_dispatcher_List_Empty(t *testing.T) {
	mockStore := &smocks.IAttackStore{}

	mockStore.On("GetAll").Return([]models.AttackDetails{})

	d := setupDispatcher(mockStore)

	got := d.List()
	if len(got) != 0 {
		t.Fail()
	}
}
