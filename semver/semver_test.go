package semver

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

func Test_setSemverFromSlice(t *testing.T) {
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
			got, err := setSemverFromSlice(tt.args.semver)
			if (err != nil) != tt.wantErr {
				t.Errorf("setSemverFromSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setSemverFromSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCorrectSemver(t *testing.T) {
	type args struct {
		minimum string
		current string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "correct minimum version",
			args: args{
				minimum: "17.03.0",
				current: "18.03.0",
			},
			want: true,
		},
		{
			name: "correct minimum version",
			args: args{
				minimum: "18.3.0",
				current: "18.13.0",
			},
			want: true,
		},
		{
			name: "correct minimum version",
			args: args{
				minimum: "18.03.0",
				current: "18.03.10",
			},
			want: true,
		},
		{
			name: "correct minimum version",
			args: args{
				minimum: "19.03.0",
				current: "18.03.10",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CorrectSemver(tt.args.minimum, tt.args.current)
			if (err != nil) != tt.wantErr {
				t.Errorf("CorrectSemver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CorrectSemver() = %v, want %v", got, tt.want)
			}
		})
	}
}
