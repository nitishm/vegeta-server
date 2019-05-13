package models

import (
	"reflect"
	"testing"
)

func TestTaskMap_Add(t *testing.T) {
	type args struct {
		attack AttackDetails
	}
	tests := []struct {
		name    string
		tm      TaskMap
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			tm:   make(TaskMap),
			args: args{
				attack: AttackDetails{
					AttackInfo: AttackInfo{
						ID: "123",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tm.Add(tt.args.attack); (err != nil) != tt.wantErr {
				t.Errorf("TaskMap.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// argsAll are args for the GetAll test case
type argsAll struct {
	filterParams FilterParams
}

// testAll is test struct for the GetAll test case
type testAll struct {
	name string
	tm   TaskMap
	args argsAll
	want []AttackDetails
}

func TestTaskMap_GetAll(t *testing.T) {
	tests := make([]testAll, 0)

	ok := testAll{
		name: "OK",
		tm: TaskMap{
			"1": AttackDetails{
				AttackInfo: AttackInfo{
					ID: "1",
				},
			},
		},
		args: argsAll{
			filterParams: make(FilterParams),
		},
		want: []AttackDetails{
			{
				AttackInfo: AttackInfo{
					ID: "1",
				},
			},
		},
	}
	tests = append(tests, ok)
	tests = append(tests, dataStatus()...)
	tests = append(tests, dataBefore()...)
	tests = append(tests, dataAfter()...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tm.GetAll(tt.args.filterParams); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskMap.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func dataStatus() []testAll {
	t := make([]testAll, 0)

	match := testAll{
		name: "OK - With Status filter match",
		tm: TaskMap{
			"1": AttackDetails{
				AttackInfo: AttackInfo{
					ID:     "1",
					Status: AttackResponseStatusCompleted,
				},
			},
		},
		args: argsAll{
			filterParams: FilterParams{
				"status": "completed",
			},
		},
		want: []AttackDetails{
			{
				AttackInfo: AttackInfo{
					ID:     "1",
					Status: AttackResponseStatusCompleted,
				},
			},
		},
	}

	mismatch := testAll{
		name: "OK - With Status filter mismatch",
		tm: TaskMap{
			"1": AttackDetails{
				AttackInfo: AttackInfo{
					ID:     "1",
					Status: AttackResponseStatusCompleted,
				},
			},
		},
		args: argsAll{
			filterParams: FilterParams{
				"status": "failed",
			},
		},
		want: []AttackDetails{},
	}

	empty := testAll{
		name: "OK - With Status filter empty",
		tm: TaskMap{
			"1": AttackDetails{
				AttackInfo: AttackInfo{
					ID:     "1",
					Status: AttackResponseStatusCompleted,
				},
			},
		},
		args: argsAll{
			filterParams: FilterParams{
				"status": "",
			},
		},
		want: []AttackDetails{
			{
				AttackInfo: AttackInfo{
					ID:     "1",
					Status: AttackResponseStatusCompleted,
				},
			},
		},
	}
	t = append(t, match, mismatch, empty)

	return t
}

func dataBefore() []testAll {
	t := make([]testAll, 0)
	match := testAll{
		name: "OK - With Created_Before filter match",
		tm: TaskMap{
			"1": AttackDetails{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Wed, 02 Jan 2019 01:00:00 MST",
				},
			},
			"2": AttackDetails{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Sun, 02 Jan 2022 01:00:00 MST",
				},
			},
		},
		args: argsAll{
			filterParams: FilterParams{
				"created_before": "2020-02-05 01:00:02",
			},
		},
		want: []AttackDetails{
			{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Wed, 02 Jan 2019 01:00:00 MST",
				},
			},
		},
	}

	failed := testAll{
		name: "OK - With Created_Before filter failed",
		tm: TaskMap{
			"1": AttackDetails{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Wed, 02 Jan 2019 01:00:00 MST",
				},
			},
		},
		args: argsAll{
			filterParams: FilterParams{
				"created_before": "bad date",
			},
		},
		want: []AttackDetails{
			{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Wed, 02 Jan 2019 01:00:00 MST",
				},
			},
		},
	}

	empty := testAll{
		name: "OK - With Created_Before filter empty",
		tm: TaskMap{
			"1": AttackDetails{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Wed, 02 Jan 2019 01:00:00 MST",
				},
			},
		},
		args: argsAll{
			filterParams: FilterParams{
				"created_before": "",
			},
		},
		want: []AttackDetails{
			{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Wed, 02 Jan 2019 01:00:00 MST",
				},
			},
		},
	}
	t = append(t, match, failed, empty)

	return t
}

func dataAfter() []testAll {
	t := make([]testAll, 0)
	match := testAll{
		name: "OK - With Created_After filter match",
		tm: TaskMap{
			"1": AttackDetails{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Wed, 02 Jan 2019 01:00:00 MST",
				},
			},
			"2": AttackDetails{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Sun, 02 Jan 2022 01:00:00 MST",
				},
			},
		},
		args: argsAll{
			filterParams: FilterParams{
				"created_after": "2020-05-17 01:02:03",
			},
		},
		want: []AttackDetails{
			{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Sun, 02 Jan 2022 01:00:00 MST",
				},
			},
		},
	}

	failed := testAll{
		name: "OK - With Created_After filter failed",
		tm: TaskMap{
			"1": AttackDetails{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Wed, 02 Jan 2019 01:00:00 MST",
				},
			},
		},
		args: argsAll{
			filterParams: FilterParams{
				"created_after": "bad date",
			},
		},
		want: []AttackDetails{
			{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Wed, 02 Jan 2019 01:00:00 MST",
				},
			},
		},
	}

	empty := testAll{
		name: "OK - With Created_After filter empty",
		tm: TaskMap{
			"1": AttackDetails{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Wed, 02 Jan 2019 01:00:00 MST",
				},
			},
		},
		args: argsAll{
			filterParams: FilterParams{
				"created_after": "",
			},
		},
		want: []AttackDetails{
			{
				AttackInfo: AttackInfo{
					ID:        "1",
					CreatedAt: "Wed, 02 Jan 2019 01:00:00 MST",
				},
			},
		},
	}
	t = append(t, match, failed, empty)

	return t
}

func TestTaskMap_GetByID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		tm      TaskMap
		args    args
		want    AttackDetails
		wantErr bool
	}{
		{
			name: "OK",
			tm: TaskMap{
				"1": AttackDetails{
					AttackInfo: AttackInfo{
						ID: "1",
					},
				},
			},
			args: args{
				id: "1",
			},
			want: AttackDetails{
				AttackInfo: AttackInfo{
					ID: "1",
				},
			},
			wantErr: false,
		},
		{
			name: "OK",
			tm: TaskMap{
				"1": AttackDetails{
					AttackInfo: AttackInfo{
						ID: "1",
					},
				},
			},
			args: args{
				id: "2",
			},
			want:    AttackDetails{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tm.GetByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskMap.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskMap.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskMap_Update(t *testing.T) {
	type args struct {
		id     string
		attack AttackDetails
	}
	tests := []struct {
		name    string
		tm      TaskMap
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			tm: TaskMap{
				"1": AttackDetails{
					AttackInfo: AttackInfo{
						ID:     "1",
						Status: AttackResponseStatusRunning,
					},
				},
			},
			args: args{
				id: "1",
				attack: AttackDetails{
					AttackInfo: AttackInfo{
						ID:     "1",
						Status: AttackResponseStatusCompleted,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error : args id mismatch",
			tm: TaskMap{
				"1": AttackDetails{
					AttackInfo: AttackInfo{
						ID:     "1",
						Status: AttackResponseStatusRunning,
					},
				},
			},
			args: args{
				id: "1",
				attack: AttackDetails{
					AttackInfo: AttackInfo{
						ID:     "2",
						Status: AttackResponseStatusCompleted,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error : ID not found",
			tm: TaskMap{
				"1": AttackDetails{
					AttackInfo: AttackInfo{
						ID:     "1",
						Status: AttackResponseStatusRunning,
					},
				},
			},
			args: args{
				id: "2",
				attack: AttackDetails{
					AttackInfo: AttackInfo{
						ID:     "2",
						Status: AttackResponseStatusCompleted,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tm.Update(tt.args.id, tt.args.attack); (err != nil) != tt.wantErr {
				t.Errorf("TaskMap.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskMap_Delete(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		tm      TaskMap
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			tm: TaskMap{
				"1": AttackDetails{
					AttackInfo: AttackInfo{
						ID:     "1",
						Status: AttackResponseStatusRunning,
					},
				},
			},
			args: args{
				id: "1",
			},
			wantErr: false,
		},
		{
			name: "OK",
			tm: TaskMap{
				"1": AttackDetails{
					AttackInfo: AttackInfo{
						ID:     "1",
						Status: AttackResponseStatusRunning,
					},
				},
			},
			args: args{
				id: "2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tm.Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("TaskMap.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewTaskMap(t *testing.T) {
	tests := []struct {
		name string
		want TaskMap
	}{
		{
			name: "OK",
			want: make(TaskMap),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTaskMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaskMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
