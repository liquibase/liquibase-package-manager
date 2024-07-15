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

var ps = Packages{driver, extension, pro}
var installedFiles []fs.FileInfo
var missingFiles []fs.FileInfo

func init() {
	rootPath, _ := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	testPath = strings.TrimRight(string(rootPath), "\n")
	installedFiles, _ = utils.ReadDir(testPath + "/tests/mocks/installed")
	missingFiles, _ = utils.ReadDir(testPath + "/tests/mocks/liquibase")
}

func TestPackages_FilterByCategory(t *testing.T) {
	type args struct {
		c string
	}
	tests := []struct {
		name string
		ps   Packages
		args args
		want Packages
	}{
		{
			name: "Can Filter by Driver",
			ps:   ps,
			args: args{"driver"},
			want: []Package{driver},
		},
		{
			name: "Can Filter by Extension",
			ps:   ps,
			args: args{"extension"},
			want: []Package{extension},
		},
		{
			name: "Can Filter by Pro",
			ps:   ps,
			args: args{"pro"},
			want: []Package{pro},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ps.FilterByCategory(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterByCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackages_GetByName(t *testing.T) {
	type args struct {
		n string
	}
	tests := []struct {
		name string
		ps   Packages
		args args
		want Package
	}{
		{
			name: "Can Get Package (Driver) By Name",
			ps:   ps,
			args: args{"driver"},
			want: driver,
		},
		{
			name: "Can Get Package (Extension) By Name",
			ps:   ps,
			args: args{"extension"},
			want: extension,
		},
		{
			name: "Can Get Package (Pro) by Name",
			ps:   ps,
			args: args{"pro"},
			want: pro,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ps.GetByName(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackages_Display(t *testing.T) {
	type args struct {
		files []fs.FileInfo
	}

	rootPath, _ := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	var installed, _ = utils.ReadDir(strings.TrimRight(string(rootPath), "\n") + "/tests/mocks/installed")
	var files, _ = utils.ReadDir(strings.TrimRight(string(rootPath), "\n") + "/tests/mocks/classpath")

	tests := []struct {
		name string
		ps   Packages
		args args
		want []string
	}{
		{
			name: "Can Display Installed Formatted Lists",
			ps:   ps,
			args: args{installed},
			want: []string{
				"     Package                                Category",
				"├──  driver@0.0.1                           driver",
				"├──  extension@1.0.0                        extension",
				"└──  pro@0.0.1                              pro",
			},
		},
		{
			name: "Can Display Uninstalled Formatted Lists",
			ps:   ps,
			args: args{files},
			want: []string{
				"     Package                                Category",
				"├──  driver                                 driver",
				"├──  extension                              extension",
				"└──  pro                                    pro",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ps.Display(tt.args.files); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackages_GetInstalled(t *testing.T) {
	type args struct {
		cpFiles []fs.FileInfo
	}
	tests := []struct {
		name string
		ps   Packages
		args args
		want Packages
	}{
		{
			name: "Can Get Installed Packages",
			ps:   ps,
			args: args{installedFiles},
			want: ps,
		},
		{
			name: "Can Confirm No Installed Packages",
			ps:   ps,
			args: args{missingFiles},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ps.GetInstalled(tt.args.cpFiles); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackages_GetOutdated(t *testing.T) {
	type args struct {
		lb      *version.Version
		cpFiles []fs.FileInfo
	}
	lb, _ := version.NewVersion("4.6.2")
	tests := []struct {
		name string
		ps   Packages
		args args
		want Packages
	}{
		{
			name: "Can Get Outdated Packages",
			ps:   ps,
			args: args{lb, installedFiles},
			want: Packages{
				Package{
					Name:     "driver",
					Category: "driver",
					Versions: []Version{driverV1, driverV2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ps.GetOutdated(tt.args.lb, tt.args.cpFiles); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOutdated() = %v, want %v", got, tt.want)
			}
		})
	}
}
