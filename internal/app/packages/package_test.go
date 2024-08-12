package packages

import (
	"io/fs"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/liquibase/liquibase-package-manager/internal/app/utils"
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
	type args struct {
		lb *version.Version
	}

	lbOld, _ := version.NewVersion("4.8.0")
	lbNew, _ := version.NewVersion("4.17.2")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   Version
	}{
		{
			name: "Can Get Latest Driver Version: LB Old",
			fields: fields{
				"test",
				"driver",
				[]Version{driverV1, driverV2},
			},
			args: args{lb: lbOld},
			want: driverV2,
		},
		{
			name: "Can Get Latest Driver Version: LB New",
			fields: fields{
				"test",
				"driver",
				[]Version{driverV1, driverV2},
			},
			args: args{lb: lbNew},
			want: driverV2,
		},
		{
			name: "Can Get Latest Extension Version: LB Old",
			fields: fields{
				"test",
				"extension",
				[]Version{extensionV1, extensionV2},
			},
			args: args{lb: lbOld},
			want: extensionV1,
		},
		{
			name: "Can Get Latest Extension Version: LB New",
			fields: fields{
				"test",
				"extension",
				[]Version{extensionV1, extensionV2},
			},
			args: args{lb: lbNew},
			want: extensionV2,
		},
		{
			name: "Can Get Latest Pro Version: LB Old",
			fields: fields{
				"test",
				"pro",
				[]Version{proV1, proV2},
			},
			args: args{lb: lbOld},
			want: proV1,
		},
		{
			name: "Can Get Latest Pro Version: LB New",
			fields: fields{
				"test",
				"pro",
				[]Version{proV1, proV2},
			},
			args: args{lb: lbNew},
			want: proV2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Package{
				Name:     tt.fields.Name,
				Category: tt.fields.Category,
				Versions: tt.fields.Versions,
			}
			if got := p.GetLatestVersion(tt.args.lb); !reflect.DeepEqual(got, tt.want) {
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
	var files, _ = utils.ReadDir(strings.TrimRight(string(rootPath), "\n") + "/tests/mocks/installed")

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

func TestPackage_DeleteVersion(t *testing.T) {
	type fields struct {
		Name     string
		Category string
		Versions []Version
	}
	type args struct {
		ver Version
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Version
	}{
		{
			name: "Can Delete Latest Extension Version",
			fields: fields{
				"test",
				"extension",
				[]Version{extensionV1, extensionV2},
			},
			args: args{ver: extensionV2},
			want: []Version{extensionV1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Package{
				Name:     tt.fields.Name,
				Category: tt.fields.Category,
				Versions: tt.fields.Versions,
			}
			if got := p.DeleteVersion(tt.args.ver); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
