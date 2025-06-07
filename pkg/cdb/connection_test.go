package cdb

import "testing"

func TestConnection_Validate(t *testing.T) {
	tests := []struct {
		name      string
		fields    Connection
		wantErr   bool
		wantedErr error
	}{
		{
			name: "simple-nickname",
			fields: Connection{
				Nickname: "something",
				Host:     "somewhere",
			},
			wantErr: false,
		},
		{
			name: "no-nickname",
			fields: Connection{
				Host: "somewhere",
			},
			wantErr:   true,
			wantedErr: ErrConnNoNickname,
		},
		{
			name: "no-host",
			fields: Connection{
				Nickname: "something",
			},
			wantErr:   true,
			wantedErr: ErrConnNoHost,
		},
		{
			name: "invalid-nickname",
			fields: Connection{
				Nickname: "23",
				Host:     "somewhere",
			},
			wantErr:   true,
			wantedErr: ErrNicknameLetter,
		},
		{
			name: "invalid-id",
			fields: Connection{
				Id:       -10,
				Nickname: "something",
				Host:     "somewhere",
			},
			wantErr:   true,
			wantedErr: ErrInvalidId,
		},
		{
			name: "zero-id",
			fields: Connection{
				Id:       0,
				Nickname: "something",
				Host:     "somewhere",
			},
			wantErr:   true,
			wantedErr: ErrConnIdZero,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields

			err := c.Validate()

			if err != nil && tt.wantErr && err != tt.wantedErr {
				t.Errorf("Connection.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
