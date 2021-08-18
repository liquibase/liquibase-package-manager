package test

import (
	"io/ioutil"
	"os/exec"
	"package-manager/pkg/lpm"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var rootPath []byte
var d lpm.Dependencies
var contextArgs *lpm.ContextArgs

func init() {
	rootPath, _ = exec.Command("git", "rev-parse", "--show-toplevel").Output()
	contextArgs = &lpm.ContextArgs{
		Path:       "/",
		WorkingDir: strings.TrimRight(string(rootPath), "\n") + "/tests/mocks/liquibase.json",
	}
	d = lpm.Dependencies{}
}

func TestDependencies_CreateFile(t *testing.T) {
	ctx, err := lpm.NewInitializedContext(contextArgs)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = d.CreateManifestFile(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}
	_, file := filepath.Split(ctx.GetManifestFilepath())
	if file != "liquibase.json" {
		t.Fatalf("Expected %s but got %s", "liquibase.json", file)
	}
}

func TestDependencies_FileExists(t *testing.T) {
	ctx, err := lpm.NewInitializedContext(contextArgs)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if d.FileExists(ctx) != true {
		t.Fatalf("Unable to verify liquibase.json file exists.")
	}
}

func TestDependencies_Write(t *testing.T) {
	ctx, err := lpm.NewInitializedContext(contextArgs)
	if err != nil {
		t.Fatalf(err.Error())
	}
	d.Append(lpm.NewDependency("package", "tag"))
	d.Append(lpm.NewDependency("package", "tag2"))
	d.Append(lpm.NewDependency("package2", "tag"))
	d.Append(lpm.NewDependency("package2", "tag3"))
	err = d.WriteManifest(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	file, _ := ioutil.ReadFile(ctx.GetManifestFilepath())
	content := `{
 "dependencies": [
  {
   "package": "tag"
  },
  {
   "package": "tag2"
  },
  {
   "package2": "tag"
  },
  {
   "package2": "tag3"
  },
 ]
}`
	if string(file) != content {
		t.Fatalf("Unable to verify liquibase.json json contents.")
	}
}

func TestDependencies_Read(t *testing.T) {
	ctx, err := lpm.NewInitializedContext(contextArgs)
	dd := lpm.NewDependencies()
	err = dd.ReadManifest(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// TODO This seems strange. Shouldn't ReadManifest never return
	//      invalid types? If it can, it should be rewritten.
	if reflect.TypeOf(dd.Dependencies[0]) != reflect.TypeOf(lpm.Dependency{}) {
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
	if len(d.Dependencies) != 0 {
		t.Fatalf("Unable to remove dependency")
	}
}
