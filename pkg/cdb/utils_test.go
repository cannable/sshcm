package cdb

import (
	"testing"
)

func TestIsValidDefault(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "args",
			args: args{
				name: "args",
			},
			want: true,
		},
		{
			name: "command",
			args: args{
				name: "command",
			},
			want: true,
		},
		{
			name: "identity",
			args: args{
				name: "identity",
			},
			want: true,
		},
		{
			name: "user",
			args: args{
				name: "user",
			},
			want: true,
		},
		{
			name: "nonsense",
			args: args{
				name: "nonsense",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidDefault(tt.args.name); got != tt.want {
				t.Errorf("IsValidDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidProperty(t *testing.T) {
	type args struct {
		property string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nickname",
			args: args{
				property: "nickname",
			},
			want: true,
		},
		{
			name: "host",
			args: args{
				property: "host",
			},
			want: true,
		},
		{
			name: "user",
			args: args{
				property: "user",
			},
			want: true,
		},
		{
			name: "description",
			args: args{
				property: "description",
			},
			want: true,
		},
		{
			name: "args",
			args: args{
				property: "args",
			},
			want: true,
		},
		{
			name: "identity",
			args: args{
				property: "identity",
			},
			want: true,
		},
		{
			name: "command",
			args: args{
				property: "command",
			},
			want: true,
		},
		{
			name: "binary",
			args: args{
				property: "binary",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidProperty(tt.args.property); got != tt.want {
				t.Errorf("IsValidProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateNickname(t *testing.T) {
	type args struct {
		nickname string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				nickname: "demo",
			},
			wantErr: false,
		},
		{
			name: "id",
			args: args{
				nickname: "60",
			},
			wantErr: true,
		},
		{
			name: "empty string",
			args: args{
				nickname: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateNickname(tt.args.nickname); (err != nil) != tt.wantErr {
				t.Errorf("ValidateNickname() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateId(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "25",
			args: args{
				id: "25",
			},
			wantErr: false,
		},
		{
			name: "blarg",
			args: args{
				id: "blarg",
			},
			wantErr: true,
		},
		{
			name: "0",
			args: args{
				id: "0",
			},
			wantErr: true,
		},
		{
			name: "-25",
			args: args{
				id: "-25",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateId(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("ValidateId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidIdOrNickname(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "25",
			args: args{
				s: "25",
			},
			want: true,
		},
		{
			name: "blarg",
			args: args{
				s: "blarg",
			},
			want: true,
		},
		{
			name: "0",
			args: args{
				s: "0",
			},
			want: false,
		},
		{
			name: "nickname-spaces",
			args: args{
				s: "asdf fdsa",
			},
			want: true,
		},
		{
			name: "-25",
			args: args{
				s: "-25",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidIdOrNickname(tt.args.s); got != tt.want {
				t.Errorf("IsValidIdOrNickname() = %v, want %v", got, tt.want)
			}
		})
	}
}
