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
		name string
		args args
		want bool
	}{
		{
			name: "true",
			args: args{
				id: 20,
			},
			want: true,
		},
		{
			name: "false",
			args: args{
				id: 3,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conndb, mock := newMockConnDb()

			if tt.want {
				rows := sqlmock.NewRows([]string{"id"}).AddRow("20")
				mock.ExpectQuery("SELECT id").WithArgs(tt.args.id).WillReturnRows(rows)
			} else {
				mock.ExpectQuery("SELECT id").WithArgs(tt.args.id)
			}

			if got := conndb.Exists(tt.args.id); got != tt.want {
				t.Errorf("ConnectionDB.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}
