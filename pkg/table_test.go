package pkg

import (
	"errors"
	"testing"
)

type testModel struct {
	ID   PrimaryKey
	Data string
}

func (m *testModel) GetID() PrimaryKey {
	return m.ID
}
func (m *testModel) SetID(id PrimaryKey) {
	m.ID = id
}

func TestTable_Insert(t *testing.T) {
	tests := []struct {
		name    string
		v       Model
		wantErr error
	}{
		{
			name:    "already has id",
			wantErr: ErrAlreadyHasID,
			v:       &testModel{ID: 1},
		},
		{
			name:    "insert",
			wantErr: nil,
			v:       &testModel{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			table := &table{
				name: "users",
				data: make(map[PrimaryKey]Model),
			}
			err := table.Insert(test.v)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("%s error = %v, wantErr %v", test.name, err, test.wantErr)
			}
			if err != nil {
				return
			}
			if test.v.GetID() == 0 {
				t.Errorf("%s id not set", test.name)
			}
		})
	}
}

func TestTable_Update(t *testing.T) {
	tests := []struct {
		name    string
		v       Model
		wantErr error
	}{
		{
			name:    "not found",
			wantErr: ErrNotFound,
			v:       &testModel{},
		},
		{
			name:    "update",
			wantErr: nil,
			v:       &testModel{ID: 1, Data: "test1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			table := &table{
				name: "users",
				data: map[PrimaryKey]Model{
					1: &testModel{ID: 1, Data: "test"},
				},
			}
			err := table.Update(test.v)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("%s error = %v, wantErr %v", test.name, err, test.wantErr)
			}
			if err != nil {
				return
			}
			if table.data[1].(*testModel).Data != "test1" {
				t.Errorf("%s data not updated", test.name)
			}
		})
	}
}

func TestTable_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      PrimaryKey
		wantErr error
	}{
		{
			name:    "not found",
			wantErr: ErrNotFound,
			id:      0,
		},
		{
			name:    "delete",
			wantErr: nil,
			id:      1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			table := &table{
				name: "users",
				data: map[PrimaryKey]Model{
					1: &testModel{ID: 1, Data: "test"},
				},
			}
			err := table.Delete(test.id)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("%s error = %v, wantErr %v", test.name, err, test.wantErr)
			}
			if err != nil {
				return
			}
			if _, ok := table.data[test.id]; ok {
				t.Errorf("%s data not deleted", test.name)
			}
		})
	}
}

func TestTable_Get(t *testing.T) {
	tests := []struct {
		name    string
		id      PrimaryKey
		want    Model
		wantErr error
	}{
		{
			name:    "not found",
			wantErr: ErrNotFound,
			id:      0,
		},
		{
			name:    "get",
			wantErr: nil,
			id:      1,
			want:    &testModel{ID: 1, Data: "test"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			table := &table{
				name: "users",
				data: map[PrimaryKey]Model{
					1: &testModel{ID: 1, Data: "test"},
				},
			}
			got, err := table.Get(test.id)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("%s error = %v, wantErr %v", test.name, err, test.wantErr)
			}
			if err != nil {
				return
			}
			if got.GetID() != test.want.GetID() {
				t.Errorf("%s id = %v, want %v", test.name, got.GetID(), test.want.GetID())
			}
		})
	}
}

func TestTable_Find(t *testing.T) {
	tests := []struct {
		name    string
		f       func(Model) bool
		want    []Model
		wantErr error
	}{
		{
			name:    "find all",
			wantErr: nil,
			f:       func(m Model) bool { return true },
			want: []Model{
				&testModel{ID: 1, Data: "test"},
				&testModel{ID: 2, Data: "test"},
			},
		},
		{
			name:    "find 1",
			wantErr: nil,
			f:       func(m Model) bool { return m.GetID() == 1 },
			want: []Model{
				&testModel{ID: 1, Data: "test"},
			},
		},
		{
			name:    "find none",
			wantErr: nil,
			f:       func(m Model) bool { return false },
			want:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			table := &table{
				name: "users",
				data: map[PrimaryKey]Model{
					1: &testModel{ID: 1, Data: "test"},
					2: &testModel{ID: 2, Data: "test"},
				},
			}
			got, err := table.Find(test.f)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("%s error = %v, wantErr %v", test.name, err, test.wantErr)
			}
			if err != nil {
				return
			}
			if len(got) != len(test.want) {
				t.Errorf("%s len = %v, want %v", test.name, len(got), len(test.want))
			}
		})
	}
}

func TestTable_Name(t *testing.T) {
	table := &table{
		name: "users",
	}
	if table.Name() != "users" {
		t.Errorf("name = %v, want %v", table.Name(), "users")
	}
}
