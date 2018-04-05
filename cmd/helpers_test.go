package cmd

import (
	"reflect"
	"testing"
)

func Test_getSemverFromString(t *testing.T) {
	type args struct {
		semver string
	}
	tests := []struct {
		name    string
		args    args
		want    semanticVersion
		wantErr bool
	}{
		{
			name: "correct semver retrieved for three parts",
			args: args{
				semver: "1.2.3",
			},
			want: semanticVersion{
				major: 1,
				minor: 2,
				patch: 3,
			},
		},
		{
			name: "correct semver retrieved for two parts",
			args: args{
				semver: "1.2",
			},
			want: semanticVersion{
				major: 1,
				minor: 2,
				patch: 0,
			},
		},
		{
			name: "correct semver retrieved for one part",
			args: args{
				semver: "1",
			},
			want: semanticVersion{
				major: 1,
				minor: 0,
				patch: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSemverFromString(tt.args.semver)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSemverFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSemverFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSemverFromSlice(t *testing.T) {
	type args struct {
		semver []string
	}
	tests := []struct {
		name    string
		args    args
		want    semanticVersion
		wantErr bool
	}{
		{
			name: "correct semver retrieved",
			args: args{
				semver: []string{"1", "2", "3"},
			},
			want: semanticVersion{
				major: 1,
				minor: 2,
				patch: 3,
			},
		},
		{
			name: "correct semver retrieved with zeros",
			args: args{
				semver: []string{"1", "2", "30"},
			},
			want: semanticVersion{
				major: 1,
				minor: 2,
				patch: 30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSemverFromSlice(tt.args.semver)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSemverFromSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSemverFromSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
