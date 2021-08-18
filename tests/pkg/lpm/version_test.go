package test

import (
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"package-manager/pkg/lpm"
	"strings"
	"testing"
)

var driverV1 = lpm.Version{
	Tag:       "0.0.1",
	Path:      "tests/mocks/files/driver-0.0.1.txt",
	Algorithm: "SHA256",
	CheckSum:  "",
}
var driverV2 = lpm.Version{
	Tag:       "0.2.0",
	Path:      "tests/mocks/files/driver-0.2.0.txt",
	Algorithm: "SHA1",
	CheckSum:  "",
}
var extensionV1 = lpm.Version{
	Tag:       "0.0.2",
	Path:      "tests/mocks/files/extension-0.0.2.txt",
	Algorithm: "SHA1",
	CheckSum:  "",
}
var extensionV2 = lpm.Version{
	Tag:       "1.0.0",
	Path:      "tests/mocks/files/extension-1.0.0.txt",
	Algorithm: "SHA1",
	CheckSum:  "",
}

var testPath string

func init() {
	rootPath, _ := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	testPath = strings.TrimRight(string(rootPath), "\n")
}

func TestVersion_CopyToClassPath(t *testing.T) {
	var files []fs.FileInfo

	extensionV1.Path = testPath + "/tests/mocks/files/extension-0.0.2.txt"
	err := extensionV1.CopyToClassPath(testPath + "/tests/mocks/classpath/")
	if err != nil {
		t.Errorf(err.Error())
		goto end
	}

	files, err = ioutil.ReadDir(testPath + "/tests/mocks/classpath/")
	if err != nil {
		t.Errorf(err.Error())
		goto end
	}

	if files[1].Name() != "extension-0.0.2.txt" {
		t.Fatalf("Expected %s but got %s", "extension-0.0.2.txt", files[0].Name())
	}
	t.Cleanup(func() {
		_ = os.Remove(testPath + "/tests/mocks/classpath/extension-0.0.2.txt")
	})
end:
}

func TestVersion_GetFilename(t *testing.T) {
	tests := []struct {
		name    string
		version lpm.Version
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
			v := lpm.Version{
				Tag:       tt.version.Tag,
				Path:      tt.version.Path,
				Algorithm: tt.version.Algorithm,
				CheckSum:  tt.version.CheckSum,
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
		version lpm.Version
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
			v := lpm.Version{
				Tag:       tt.version.Tag,
				Path:      tt.version.Path,
				Algorithm: tt.version.Algorithm,
				CheckSum:  tt.version.CheckSum,
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
		version lpm.Version
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
			v := lpm.Version{
				Tag:       tt.version.Tag,
				Path:      tt.version.Path,
				Algorithm: tt.version.Algorithm,
				CheckSum:  tt.version.CheckSum,
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
		version lpm.Version
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
			v := lpm.Version{
				Tag:       tt.version.Tag,
				Path:      tt.version.Path,
				Algorithm: tt.version.Algorithm,
				CheckSum:  tt.version.CheckSum,
			}
			got, err := v.CalcChecksum(tt.args.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			if got != tt.want {
				t.Errorf("CalcChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}
