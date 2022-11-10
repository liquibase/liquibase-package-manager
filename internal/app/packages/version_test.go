package packages

import (
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var driverV1 = Version{
	Tag:           "0.0.1",
	Path:          "tests/mocks/files/driver-0.0.1.txt",
	Algorithm:     "SHA256",
	CheckSum:      "",
	LiquibaseCore: "v4.6.2",
}
var driverV2 = Version{
	Tag:           "0.2.0",
	Path:          "tests/mocks/files/driver-0.2.0.txt",
	Algorithm:     "SHA1",
	CheckSum:      "",
	LiquibaseCore: "",
}
var extensionV1 = Version{
	Tag:           "0.0.2",
	Path:          "tests/mocks/files/extension-0.0.2.txt",
	Algorithm:     "SHA1",
	CheckSum:      "",
	LiquibaseCore: "",
}
var extensionV2 = Version{
	Tag:           "1.0.0",
	Path:          "tests/mocks/files/extension-1.0.0.txt",
	Algorithm:     "SHA1",
	CheckSum:      "",
	LiquibaseCore: "",
}
var proV1 = Version{
	Tag:           "0.0.1",
	Path:          "tests/mocks/files/pro-0.0.1.txt",
	Algorithm:     "SHA1",
	CheckSum:      "",
	LiquibaseCore: "",
}
var proV2 = Version{
	Tag:           "0.0.2",
	Path:          "tests/mocks/files/pro-0.0.2.txt",
	Algorithm:     "SHA1",
	CheckSum:      "",
	LiquibaseCore: "",
}

var testPath string

func init() {
	rootPath, _ := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	testPath = strings.TrimRight(string(rootPath), "\n")
}

func TestVersion_CopyToClassPath(t *testing.T) {
	extensionV1.Path = testPath + "/tests/mocks/files/extension-0.0.2.txt"
	extensionV1.CopyToClassPath(testPath + "/tests/mocks/classpath/")
	var files, _ = ioutil.ReadDir(testPath + "/tests/mocks/classpath/")

	if files[1].Name() != "extension-0.0.2.txt" {
		t.Fatalf("Expected %s but got %s", "extension-0.0.2.txt", files[0].Name())
	}
	t.Cleanup(func() {
		os.Remove(testPath + "/tests/mocks/classpath/extension-0.0.2.txt")
	})
}

func TestVersion_GetFilename(t *testing.T) {
	tests := []struct {
		name    string
		version Version
		want    string
	}{
		{
			name:    "Can Get Filename",
			version: driverV1,
			want:    "driver-0.0.1.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				Tag:           tt.version.Tag,
				Path:          tt.version.Path,
				Algorithm:     tt.version.Algorithm,
				CheckSum:      tt.version.CheckSum,
				LiquibaseCore: tt.version.LiquibaseCore,
			}
			if got := v.GetFilename(); got != tt.want {
				t.Errorf("GetFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_InClassPath(t *testing.T) {
	type args struct {
		files []fs.FileInfo
	}

	var files, _ = ioutil.ReadDir(testPath + "/tests/mocks/installed")

	tests := []struct {
		name    string
		version Version
		args    args
		want    bool
	}{
		{
			name:    "Can Determine if Package in Classpath",
			version: extensionV2,
			args:    args{files},
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				Tag:           tt.version.Tag,
				Path:          tt.version.Path,
				Algorithm:     tt.version.Algorithm,
				CheckSum:      tt.version.CheckSum,
				LiquibaseCore: tt.version.LiquibaseCore,
			}
			if got := v.InClassPath(tt.args.files); got != tt.want {
				t.Errorf("InClassPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_PathIsHttp(t *testing.T) {
	tests := []struct {
		name    string
		version Version
		want    bool
	}{
		{
			name:    "Can Determine if Path is Http",
			version: driverV2,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				Tag:           tt.version.Tag,
				Path:          tt.version.Path,
				Algorithm:     tt.version.Algorithm,
				CheckSum:      tt.version.CheckSum,
				LiquibaseCore: tt.version.LiquibaseCore,
			}
			if got := v.PathIsHTTP(); got != tt.want {
				t.Errorf("PathIsHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_calcChecksum(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		version Version
		args    args
		want    string
	}{
		{
			name:    "Can Calc Checksum SHA1",
			version: driverV2,
			args:    args{[]byte("DriverSHA1")},
			want:    "70daefe06dd19c073920273e02cfc712951795ea",
		},
		{
			name:    "Can Calc Checksum SHA256",
			version: driverV1,
			args:    args{[]byte("DriverSHA256")},
			want:    "94c74ea180983e2ec16451fed233c9f6d3d47572133cae84a0adc7c9fd7e1dd4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				Tag:           tt.version.Tag,
				Path:          tt.version.Path,
				Algorithm:     tt.version.Algorithm,
				CheckSum:      tt.version.CheckSum,
				LiquibaseCore: tt.version.LiquibaseCore,
			}
			if got := v.calcChecksum(tt.args.b); got != tt.want {
				t.Errorf("calcChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}
