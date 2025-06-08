package cdb

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestConnectionDB_AddSimple(t *testing.T) {
	wantErr := false

	c := &Connection{
		Nickname: "something",
		Host:     "somewhere",
	}

	conndb, mock := newMockConnDb()

	// Pretend that there are no existing connections with this nickname
	rows := sqlmock.NewRows([]string{"id"})
	mock.ExpectQuery("SELECT id").WithArgs(c.Nickname).WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO connections").WithArgs(
		"something",
		"somewhere",
		"",
		"",
		"",
		"",
		"",
	).WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := conndb.Add(c)

	conndb.Close()

	if (err != nil) != wantErr {
		t.Errorf("ConnectionDB.Add() error = %v, wantErr %v", err, wantErr)
		return
	}
}

func TestConnectionDB_AddNoNickname(t *testing.T) {
	conndb, _ := newMockConnDb()

	wantErr := ErrConnNoNickname
	_, err := conndb.Add(&Connection{
		Host: "somewhere",
	})

	conndb.Close()

	if err != wantErr {
		t.Errorf("ConnectionDB.Add() error = %v, want %v", err, wantErr)
		return
	}
}

func TestConnectionDB_AddBadNickname(t *testing.T) {
	conndb, _ := newMockConnDb()

	wantErr := ErrNicknameLetter
	_, err := conndb.Add(&Connection{
		Nickname: "700",
	})

	conndb.Close()

	if err != wantErr {
		t.Errorf("ConnectionDB.Add() error = %v, want %v", err, wantErr)
		return
	}
}

func TestConnectionDB_AddNoHost(t *testing.T) {
	conndb, _ := newMockConnDb()

	wantErr := ErrConnNoHost
	_, err := conndb.Add(&Connection{
		Nickname: "something",
	})

	conndb.Close()

	if err != wantErr {
		t.Errorf("ConnectionDB.Add() error = %v, want %v", err, wantErr)
		return
	}
}

func TestConnectionDB_Exists(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "true",
			args: args{
				id: 20,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "false",
			args: args{
				id: 3,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test DB mocking
			conndb, mock := newMockConnDb()

			var rows *sqlmock.Rows

			if tt.want {
				rows = sqlmock.NewRows([]string{"id"}).AddRow("20")
			} else {
				rows = sqlmock.NewRows([]string{"id"})
			}

			mock.ExpectQuery("SELECT id").WithArgs(tt.args.id).WillReturnRows(rows)

			got, err := conndb.Exists(tt.args.id)

			if err != nil && !tt.wantErr {
				// Got an error when we didn't want one
				t.Errorf("ConnectionDB.Exists() error = %v, wantErr %v", err, tt.wantErr)
			} else if got != tt.want {
				// Got the wrong result
				t.Errorf("ConnectionDB.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnectionDB_ExistsByProperty(t *testing.T) {
	type args struct {
		property string
		value    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "true-nickname",
			args: args{
				property: "nickname",
				value:    "something",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "false-nickname",
			args: args{
				property: "nickname",
				value:    "something",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "bad-property",
			args: args{
				property: "blarg",
				value:    "won't work",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test DB mocking
			conndb, mock := newMockConnDb()

			var rows *sqlmock.Rows

			if tt.want {
				rows = sqlmock.NewRows([]string{"id"}).AddRow("17")
			} else {
				rows = sqlmock.NewRows([]string{"id"})
			}

			// We only expect the value arg, as the property gets written into the query
			mock.ExpectQuery("SELECT id").WithArgs(tt.args.value).WillReturnRows(rows)

			got, err := conndb.ExistsByProperty(tt.args.property, tt.args.value)

			if err != nil && !tt.wantErr {
				// Got an error when we didn't want one
				t.Errorf("ConnectionDB.ExistsByProperty() error = %v, wantErr %v", err, tt.wantErr)
			} else if got != tt.want {
				// Got the wrong result
				t.Errorf("ConnectionDB.ExistsByProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}
