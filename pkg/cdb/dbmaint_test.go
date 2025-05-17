package cdb

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// newMockConnDb: Returns a new ConnectionDB with a sqlmock connection
// (and associated Sqlmock).
func newMockConnDb() (*ConnectionDB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()

	conndb := &ConnectionDB{
		connection: db,
	}

	if err != nil {
		panic(err)
	}

	return conndb, mock
}

func Test_validateDbSchemaVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		//wantErr bool
		want error
	}{
		{
			name: "v1.1",
			args: args{
				version: "v1.1",
			},
			want: nil,
		},
		{
			name: "empty",
			args: args{
				version: "",
			},
			want: ErrSchemaVerInvalid,
		},
		{
			name: "1.0",
			args: args{
				version: "1.0",
			},
			want: ErrSchemaUpgradeNeeded,
		},
		{
			name: "v100.99.88",
			args: args{
				version: "v100.99.88",
			},
			want: ErrSchemaTooNew,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateDbSchemaVersion(tt.args.version); got != tt.want {
				t.Errorf("validateDbSchemaVersion() error = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnectionDB_GetDbSchemaVersion(t *testing.T) {
	tests := []struct {
		name        string
		want        string
		wantErr     bool
		expectQuery string
		cols        []string
		row         string
	}{
		{
			name:    "backwards_compat-1.1",
			want:    "1.1",
			wantErr: false,
			row:     "1.1",
		},
		{
			name:    "v1.1",
			want:    "v1.1",
			wantErr: false,
			row:     "v1.1",
		},
		{
			name:    "weird_garbage",
			want:    "asdf fdsa",
			wantErr: false,
			row:     "asdf fdsa",
		},
		{
			name:    "no_rows",
			want:    "",
			wantErr: true,
			row:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare test mocking
			conndb, mock := newMockConnDb()

			if !tt.wantErr {
				rows := sqlmock.NewRows([]string{"value"}).AddRow(tt.row)
				mock.ExpectQuery("SELECT value").WillReturnRows(rows)
			}

			got, err := conndb.GetDbSchemaVersion()

			defer conndb.connection.Close()

			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectionDB.GetDbSchemaVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ConnectionDB.GetDbSchemaVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
