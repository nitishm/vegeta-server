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

func TestTaskMap_GetAll(t *testing.T) {
	type args struct {
		filterParams FilterParams
	}
	tests := []struct {
		name string
		tm   TaskMap
		args args
		want []AttackDetails
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
				filterParams: make(FilterParams),
			},
			want: []AttackDetails{
				{
					AttackInfo: AttackInfo{
						ID: "1",
					},
				},
			},
		},
		{
			name: "OK - With Status filter match",
			tm: TaskMap{
				"1": AttackDetails{
					AttackInfo: AttackInfo{
						ID:     "1",
						Status: AttackResponseStatusCompleted,
					},
				},
			},
			args: args{
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
		},
		{
			name: "OK - With Status filter mismatch",
			tm: TaskMap{
				"1": AttackDetails{
					AttackInfo: AttackInfo{
						ID:     "1",
						Status: AttackResponseStatusCompleted,
					},
				},
			},
			args: args{
				filterParams: FilterParams{
					"status": "failed",
				},
			},
			want: []AttackDetails{},
		},
		{
			name: "OK - With Status filter empty",
			tm: TaskMap{
				"1": AttackDetails{
					AttackInfo: AttackInfo{
						ID:     "1",
						Status: AttackResponseStatusCompleted,
					},
				},
			},
			args: args{
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tm.GetAll(tt.args.filterParams); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskMap.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
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
