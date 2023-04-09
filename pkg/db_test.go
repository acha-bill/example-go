package pkg

import (
	"errors"
	"testing"
)

func TestDb_AddTable(t *testing.T) {
	tests := []struct {
		name      string
		tableName string
		d         *db
		wantErr   error
	}{
		{
			name:      "no table name",
			tableName: "",
			wantErr:   ErrorNoTableName,
		},
		{
			name:      "table already exists",
			tableName: "users",
			wantErr:   ErrTableExists,
			d: &db{
				tables: map[string]Table{
					"users": nil,
				},
			},
		},
		{
			name:      "add table",
			tableName: "users",
			wantErr:   nil,
			d: &db{
				tables: map[string]Table{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.d.AddTable(test.tableName)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("%s error = %v, wantErr %v", test.name, err, test.wantErr)
			}
			if err != nil {
				return
			}
			if _, ok := test.d.tables[test.tableName]; !ok {
				t.Errorf("%s table not found", test.name)
			}
		})
	}
}

func TestDb_Tables(t *testing.T) {
	tests := []struct {
		name string
		d    *db
		want []string
	}{
		{
			name: "empty",
			d:    &db{},
			want: []string{},
		},
		{
			name: "one table",
			d: &db{
				tables: map[string]Table{
					"users": nil,
				},
			},
			want: []string{"users"},
		},
		{
			name: "two tables",
			d: &db{
				tables: map[string]Table{
					"users":  nil,
					"groups": nil,
				},
			},
			want: []string{"groups", "users"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.d.Tables()
			if len(got) != len(test.want) {
				t.Errorf("%s got = %v, want %v", test.name, got, test.want)
			}
			for i := range got {
				if got[i] != test.want[i] {
					t.Errorf("%s got = %v, want %v", test.name, got, test.want)
				}
			}
		})
	}
}

func TestDb_Table(t *testing.T) {
	tests := []struct {
		name      string
		tableName string
		d         *db
		want      Table
		wantErr   error
	}{
		{
			name:      "no table name",
			tableName: "",
			wantErr:   ErrorNoTableName,
		},
		{
			name:      "table not found",
			tableName: "users",
			d: &db{
				tables: map[string]Table{},
			},
			wantErr: ErrorNoTable,
		},
		{
			name:      "table found",
			tableName: "users",
			want:      &table{},
			d: &db{
				tables: map[string]Table{
					"users": &table{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.d.Table(test.tableName)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("%s error = %v, wantErr %v", test.name, err, test.wantErr)
			}
			if err != nil {
				return
			}
			if got != nil && test.want == nil {
				t.Errorf("%s got = %v, want %v", test.name, got, test.want)
			}
		})
	}
}
