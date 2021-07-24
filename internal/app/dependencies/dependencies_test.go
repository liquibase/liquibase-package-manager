package dependencies

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var d Dependencies

func init() {
	rootPath, _ := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	FileLocation = strings.TrimRight(string(rootPath), "\n") + "/tests/mocks/liquibase.json"
	d = Dependencies{}
}

func TestDependencies_CreateFile(t *testing.T) {
	d.CreateFile()
	_, file := filepath.Split(FileLocation)
	if file != "liquibase.json" {
		t.Fatalf("Expected %s but got %s", "liquibase.json", file )
	}
}

func TestDependencies_FileExists(t *testing.T) {
	if d.FileExists() != true {
		t.Fatalf( "Unable to verify liquibase.json file exists." )
	}
}

func TestDependencies_Write(t *testing.T) {
	d.Dependencies = append(d.Dependencies, Dependency{"package": "tag"})
	d.Write()

	file, _ := ioutil.ReadFile(FileLocation)
	content := `{
 "dependencies": [
  {
   "package": "tag"
  }
 ]
}`
	if string(file) != content {
		t.Fatalf( "Unable to verify liquibase.json json contents." )
	}
}

func TestDependencies_Read(t *testing.T) {
	dd := Dependencies{}
	dd.Read()
	if reflect.TypeOf(dd.Dependencies[0]) != reflect.TypeOf(Dependency{}) {
		t.Fatalf("Unable to load Dendency from file")
	}
	for k, v := range dd.Dependencies[0] {
		if k != "package" {
			t.Fatalf("Invalid Key")
		}
		if v != "tag" {
			t.Fatalf("Invalid Value")
		}
	}
}

func TestDependencies_Remove(t *testing.T) {
	d.Remove("package")
	if len(d.Dependencies) != 0  {
		t.Fatalf( "Unable to remove dependency" )
	}
}
