package cdb

import "testing"

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
			if got := validateDbSchemaVersion(tt.args.version); got != tt.want {
				t.Errorf("validateDbSchemaVersion() error = %v, want %v", got, tt.want)
			}
		})
	}
}
