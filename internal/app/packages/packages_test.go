package packages

import (
	"io/fs"
	"io/ioutil"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

var ps = Packages{driver, extension}
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
			ps: ps,
			args: args{"driver"},
			want: []Package{driver},
		},
		{
			name: "Can Filter by Extension",
			ps: ps,
			args: args{"extension"},
			want: []Package{extension},
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
	var files, _ = ioutil.ReadDir(strings.TrimRight(string(rootPath), "\n") + "/tests/mocks/installed")

	tests := []struct {
		name string
		ps   Packages
		args args
		want []string
	}{
		{
			name: "Can Display Formatted Lists",
			ps: ps,
			args: args{files},
			want: []string{
				"     Package                                Category",
				"├──  driver@0.0.1                           driver",
				"└──  extension@1.0.0                        extension",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ps.Display(tt.args.files); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Display() = %v, want %v", got, tt.want)
			}
		})
	}
}