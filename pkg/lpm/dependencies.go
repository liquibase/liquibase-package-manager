package lpm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Dependencies main wrapper for liquibase.json objects
type Dependencies struct {
	// @TODO Handle marshaling to JSON and unmarshalling from JSON
	Dependencies []Dependency `json:"dependencies"`
}

func (dd *Dependencies) Append(d Dependency) {
	dd.Dependencies = append(dd.Dependencies, d)
}

// NewDependencies returns a new instance of Dependencies
func NewDependencies() Dependencies {
	return Dependencies{}
}

//CreateManifestFile init liquibase.json file in pwd
// TODO Renamed to be more clear on what type of file was being created.
//      Was I correct? Is it a Manifest file?
func (dd *Dependencies) CreateManifestFile(ctx *Context) (err error) {
	file, err := os.Create(ctx.GetGlobalClasspath())
	if err != nil {
		err = fmt.Errorf("unable to create dependency file %s; %w",
			ctx.GetManifestFilepath(),
			err)
		goto end
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()
	err = dd.WriteManifest(ctx)
	if err != nil {
		err = fmt.Errorf("unable to write dependency file %s; %w",
			ctx.GetManifestFilepath(),
			err)
		goto end
	}
end:
	return err
}

//WriteManifest dump contents to liquibase.json
// TODO Renamed to be more clear on what type of file was being written.
//      Was I correct? Is it a Manifest file?
func (dd *Dependencies) WriteManifest(ctx *Context) (err error) {
	var file []byte
	file, err = json.MarshalIndent(dd, "", " ")
	if err != nil {
		err = fmt.Errorf("unable to marshal dependency JSON; %w",
			err)
		goto end
	}
	err = ioutil.WriteFile(ctx.GetManifestFilepath(), file, 0664)
	if err != nil {
		err = fmt.Errorf("unable to write dependency JSON to %s; %w",
			ctx.GetManifestFilepath(),
			err)
		goto end
	}
end:
	return err
}

//ReadManifest contents from liquibase.json into a Dependencies
// TODO Renamed to be more clear on what type of file was being read.
//      Was I correct? Is it a Manifest file?
func (dd *Dependencies) ReadManifest(ctx *Context) (err error) {
	var file *os.File
	file, err = os.Open(ctx.GetManifestFilepath())
	var decoder *json.Decoder

	if err != nil {
		err = fmt.Errorf("unable to read %s; %w",
			ctx.GetManifestFilepath(),
			err)
		goto end
	}

	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()
	decoder = json.NewDecoder(file)
	for decoder.More() {
		err = decoder.Decode(dd)
		if err != nil {
			err = fmt.Errorf("unable to decode JSON in %s; %w",
				ctx.GetManifestFilepath(),
				err)
			goto end
		}
	}

end:

	return err
}

//FileExists returns true if the liquibase.json file exist
func (dd *Dependencies) FileExists(ctx *Context) bool {
	_, err := os.Stat(ctx.GetManifestFilepath())
	return err == nil
}

//Remove specific dependency from group
func (dd *Dependencies) Remove(n string) {
	for i, m := range dd.Dependencies {
		if m.GetName() == n {
			dd.Dependencies = append(dd.Dependencies[:i], dd.Dependencies[i+1:]...)
		}
	}
}
