package packages

import (
	"io/fs"
	"io/ioutil"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

var driver = Package{
	"driver",
	"driver",
	[]Version{driverV1, driverV2},
}
var extension = Package{
	"extension",
	"extension",
	[]Version{extensionV1, extensionV2},
}
var pro = Package{
	"pro",
	"pro",
	[]Version{proV1, proV2},
}

func TestPackage_GetLatestVersion(t *testing.T) {
	type fields struct {
		Name     string
		Category string
		Versions []Version
	}
	tests := []struct {
		name   string
		fields fields
		want   Version
	}{
		{
			name: "Can Get Latest Version",
			fields: fields{
				"test",
				"driver",
				[]Version{driverV1, driverV2},
			},
			want: driverV2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Package{
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
		Versions []Version
	}
	type args struct {
		v string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Version
	}{
		{
			name: "Can get Specific Version",
			fields: fields{
				"test",
				"driver",
				[]Version{driverV1},
			},
			args: args{"0.0.1"},
			want: driverV1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Package{
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
		Versions []Version
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
		want   Version
	}{
		{
			name: "Can Get Installed Version",
			fields: fields{
				"test",
				"driver",
				[]Version{driverV1, driverV2},
			},
			args: args{files},
			want: driverV1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Package{
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