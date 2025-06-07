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
