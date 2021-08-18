package test

import (
	"io/fs"
	"io/ioutil"
	"os/exec"
	"package-manager/pkg/lpm"
	"reflect"
	"strings"
	"testing"
)

var driver = lpm.Package{
	"driver",
	"driver",
	[]lpm.Version{driverV1, driverV2},
}
var extension = lpm.Package{
	"extension",
	"extension",
	[]lpm.Version{extensionV1, extensionV2},
}

func TestPackage_GetLatestVersion(t *testing.T) {
	type fields struct {
		Name     string
		Category string
		Versions []lpm.Version
	}
	tests := []struct {
		name   string
		fields fields
		want   lpm.Version
	}{
		{
			name: "Can Get Latest Version",
			fields: fields{
				"test",
				"driver",
				[]lpm.Version{driverV1, driverV2},
			},
			want: driverV2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := lpm.Package{
				Name:     tt.fields.Name,
				Category: tt.fields.Category,
				Versions: tt.fields.Versions,
			}
			if got := p.GetLatestVersion(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLatestVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackage_GetVersion(t *testing.T) {
	type fields struct {
		Name     string
		Category string
		Versions []lpm.Version
	}
	type args struct {
		v string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   lpm.Version
	}{
		{
			name: "Can get Specific Version",
			fields: fields{
				"test",
				"driver",
				[]lpm.Version{driverV1},
			},
			args: args{"0.0.1"},
			want: driverV1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := lpm.Package{
				Name:     tt.fields.Name,
				Category: tt.fields.Category,
				Versions: tt.fields.Versions,
			}
			if got := p.GetVersion(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackage_GetInstalledVersion(t *testing.T) {
	type fields struct {
		Name     string
		Category string
		Versions []lpm.Version
	}
	type args struct {
		files []fs.FileInfo
	}

	rootPath, _ := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	var files, _ = ioutil.ReadDir(strings.TrimRight(string(rootPath), "\n") + "/tests/mocks/installed")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   lpm.Version
	}{
		{
			name: "Can Get Installed Version",
			fields: fields{
				"test",
				"driver",
				[]lpm.Version{driverV1, driverV2},
			},
			args: args{files},
			want: driverV1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := lpm.Package{
				Name:     tt.fields.Name,
				Category: tt.fields.Category,
				Versions: tt.fields.Versions,
			}
			if got := p.GetInstalledVersion(tt.args.files); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInstalledVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
