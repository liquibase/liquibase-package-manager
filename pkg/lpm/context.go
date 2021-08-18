package lpm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FileSource byte

const GlobalFiles FileSource = 'g'
const LocalFiles FileSource = 'l'

const PackageFilePerms = 0664
const DirectoryPerms = 0775
const DefaultPackageFile = "packages.json"

// ClasspathFiles represents either Global or Local Classpath ClasspathFiles
type ClasspathFiles []fs.FileInfo

func (files ClasspathFiles) VersionExists(v Version) bool {
	r := false
	for _, f := range files {
		if f.Name() == v.GetFilename() {
			r = true
			break
		}
	}
	return r
}

// Context contains
type Context struct {
	context.Context

	WorkingDir string

	manifestFilepath string
	FileSource       FileSource
	Category         string
	HomeDir          string
	Path             string

	//PackageFile to use
	//ex: package.json
	PackageFile string

	// bytes for storing JSON read from file
	packageBytes []byte

	classpath string

	classpathFiles ClasspathFiles

	packages Packages

	errors []error
}

func (ctx *Context) GetPackageByName(n string) Package {
	return ctx.packages.GetByName(n)
}

// PrintJavaOptsHelper
// TODO Test this on windows
func (ctx *Context) PrintJavaOptsHelper() {
	jo := fmt.Sprintf("-cp liquibase_libs%s*:%s*:%sliquibase.jar",
		string(os.PathSeparator),
		ctx.GetGlobalClasspath(),
		ctx.HomeDir)
	fmt.Println()
	fmt.Println("---------- IMPORTANT ----------")
	fmt.Println("Run the following JAVA_OPTS command:")
	fmt.Printf("\n\texport JAVA_OPTS=\"%s\"", jo)
}

func ContextFromCobraCommand(cmd *cobra.Command) *Context {
	ctx, ok := cmd.Context().(*Context)
	if !ok {
		// When !ok it is a programming error. So okay to just exit
		fmt.Printf("Unable to assert *app.Context from *cobra.Context:\n")
		os.Exit(1)
	}
	return ctx
}

func NewContext(path string) *Context {
	return &Context{
		Context:        context.Background(),
		PackageFile:    DefaultPackageFile,
		classpathFiles: make(ClasspathFiles, 0),
		errors:         make([]error, 0),
		Path:           path,
	}
}

type ContextArgs struct {
	Path       string
	WorkingDir string
}

// NewInitializedContext creates a new context and initializes it,
// with optional arguments and returning an error on failure
func NewInitializedContext(args *ContextArgs) (ctx *Context, err error) {
	ctx = NewContext(ctx.Path)
	err = ctx.Initialize()
	if args.WorkingDir != "" {
		ctx.WorkingDir = args.WorkingDir
	}
	return ctx, err
}

// GetManifestFilepath gets filepath for `liquibase.json` manifest file
func (ctx *Context) GetManifestFilepath() string {
	if ctx.manifestFilepath != "" {
		goto end
	}
	ctx.manifestFilepath = fmt.Sprintf("%s%sliquibase.json",
		string(os.PathSeparator),
		ctx.WorkingDir)
end:
	return ctx.manifestFilepath
}

// GetGlobalPackageFilepath gets the global class filepath
func (ctx *Context) GetGlobalPackageFilepath() string {
	return ctx.GetClassPackageFilepath(ctx.GetGlobalClasspath())
}

// GetLocalPackageFilepath gets the local class filepath
func (ctx *Context) GetLocalPackageFilepath() (pfp string, err error) {
	var cp string
	cp, err = ctx.GetLocalClasspath()
	if err != nil {
		err = fmt.Errorf("unable to get local package filepath; %w", err)
		goto end
	}
	pfp = ctx.GetClassPackageFilepath(cp)
end:
	return pfp, err
}

// GetClassPackageFilepath gets the class filepath for classpath specified
func (ctx *Context) GetClassPackageFilepath(path string) string {
	return fmt.Sprintf("%s%s",
		path,
		ctx.PackageFile)
}

func (ctx *Context) GetGlobalClasspath() string {
	return fmt.Sprintf("%slib%s",
		ctx.HomeDir,
		ctx.Path)
}

func (ctx *Context) maybeCreateInitialPackageFile() (err error) {
	if ctx.PackagesInClassPath() {
		goto end
	}
	// TODO Should this be copying to GLOBAL classpath, or current classpath?
	err = ctx.CopyPackageBytesToClassPath(PackagesJSON)
	if err != nil {
		err = fmt.Errorf("unable to create initial packages.json in %s; %w",
			ctx.GetGlobalClasspath(),
			err)
		goto end
	}
end:
	return err
}

// Initialize ensure the Context object has required values
func (ctx *Context) Initialize() (err error) {

	ctx.HomeDir, err = GetHomeDir()
	if err != nil {
		err = fmt.Errorf("unable to get home directory for lbm; %w", err)
		goto end
	}

	ctx.WorkingDir, err = GetWorkingDir()
	if err != nil {
		err = fmt.Errorf("unable to get current working directory for lbm; %w", err)
		goto end
	}

	//Install Embedded Package File
	err = ctx.maybeCreateInitialPackageFile()
	if err != nil {
		goto end
	}

	//Load Packages
	//@TODO Should this be loaded from Global Package Filepath, or
	//      loaded from Local?
	err = ctx.LoadPackages(ctx.GetGlobalPackageFilepath())

	// Set global vs local classpath
	ctx.FileSource = LocalFiles

end:

	return err

}

// GetLocalClasspath returns the current working directory
// with the path `/liquibase_libs/` added.
//
func (ctx *Context) GetLocalClasspath() (cp string, err error) {

	var pwd string
	pwd, err = os.Getwd()
	if err != nil {
		err = fmt.Errorf("unable to get working directory for local classpath; %w",
			err)
		goto end
	}

	cp = fmt.Sprintf("%s%sliquibase_libs%s",
		string(os.PathSeparator),
		pwd,
		string(os.PathSeparator))

end:

	return cp, err
}

//GetFiles returns list of files
func (ctx *Context) GetFiles(path string) (files ClasspathFiles, err error) {
	err = os.Mkdir(path, DirectoryPerms)
	if err != nil {
		err = fmt.Errorf("unable to make directory %s; %w",
			path,
			err)
		goto end
	}

	files, err = ioutil.ReadDir(path)
	if err != nil {
		err = fmt.Errorf("unable to list files in directory %s; %w",
			path,
			err)
		goto end
	}
end:
	return files, err
}

//GetClasspath returns global or local depending on FileSource
func (ctx *Context) GetClasspath() (cp string, err error) {

	switch ctx.FileSource {

	case GlobalFiles:
		cp = ctx.GetGlobalClasspath()

	case LocalFiles:
		cp, err = ctx.GetLocalClasspath()
		if err != nil {
			err = fmt.Errorf("unable to get local classpath; %w",
				err)
		}
	}

	return cp, err
}

//GetClasspathFiles returns global or local files depending on FileSource
// @TODO consider catching the files to a private instance variable on 1st call
func (ctx *Context) GetClasspathFiles() (files ClasspathFiles, cp string, err error) {
	cp, err = ctx.GetClasspath()
	if err != nil {
		err = fmt.Errorf("unable to get classpath; %w",
			err)
		goto end
	}
	files, err = ctx.GetFiles(cp)
	if err != nil {
		err = fmt.Errorf("unable to get classpath files; %w",
			err)
		goto end
	}

end:
	return files, cp, err
}

//PackagesInClassPath is the packages.json file in global classpath
func (ctx *Context) PackagesInClassPath() bool {
	// TODO Should this be Global classpath or current classpath?
	_, err := os.Stat(ctx.GetGlobalClasspath())
	return err == nil
}

//CopyPackageBytesToClassPath install packages.json to provided classpath
func (ctx *Context) copyPackagesToClassPath(p []byte, cp string) (err error) {
	err = ioutil.WriteFile(cp, p, PackageFilePerms)
	if err != nil {
		err = fmt.Errorf("unable to write to %s; %w", cp, err)
	}
	return err
}

//CopyPackageBytesToClassPath install packages.json to provided classpath
func (ctx *Context) CopyPackageBytesToClassPath(b []byte) (err error) {
	var cp string
	cp, err = ctx.GetClasspath()
	if err != nil {
		err = fmt.Errorf("unable to copy packages to classpath %s; %w",
			cp,
			err)
		goto end
	}
	err = ctx.copyPackagesToClassPath(b, cp)
end:
	return err
}

//CopyPackagesToGlobalClassPath install packages.json to global classpath
func (ctx *Context) CopyPackagesToGlobalClassPath() (err error) {
	return ctx.copyPackagesToClassPath(
		ctx.packageBytes,
		ctx.GetGlobalClasspath())
}

//UnmarshalPackages from bytes
func (ctx *Context) UnmarshalPackages(b []byte) (packs Packages, err error) {
	packs = make(Packages, 0)

	//Unmarshal the JSON into array of Packages
	err = json.Unmarshal(b, &packs)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal packages from JSON; %w",
			err)
		goto end
	}
end:
	return packs, err
}

//LoadPackages from filepath
func (ctx *Context) LoadPackages(fp string) (err error) {
	// JSON Package File
	var jpf *os.File

	// error action
	var action string

	//Open the JSON Package File
	jpf, err = os.Open(fp)
	if err != nil {
		action = "open"
		goto end
	}

	//Get Bytes from JSON Package File
	ctx.packageBytes, err = ioutil.ReadAll(jpf)
	if err != nil {
		action = "read from"
		goto end
	}

	//Unmarshal the JSON into array of Packages
	ctx.packages, err = ctx.UnmarshalPackages(ctx.packageBytes)
	if err != nil {
		action = "load packages from"
		goto end
	}

end:
	if action != "" {
		err = fmt.Errorf("unable to %s %s; %w",
			action,
			jpf.Name(),
			err)
	}
	return err
}

//SearchPackages returns previously loaded packages filtered by `name`
func (ctx *Context) SearchPackages(name string) (p Packages) {
	pp := make(Packages, 0)
	if name == "" {
		copy(pp, ctx.GetPackages())
		goto end
	}
	for _, p := range ctx.GetPackages() {
		if !strings.Contains(p.Name, name) {
			continue
		}
		pp = append(pp, p)
	}
end:
	return pp
}

//GetPackages returns previously loaded packages
func (ctx *Context) GetPackages() (p Packages) {
	return ctx.packages
}

//GetFilteredPackages returns previously loaded packages filtered by category
func (ctx *Context) GetFilteredPackages() (p Packages) {
	return ctx.GetPackages().FilterByCategory(ctx.Category)
}

//WritePackages write packages back to file
func (ctx *Context) WritePackages(p Packages) (err error) {
	var b []byte
	var pwd string
	var jf string

	b, err = json.MarshalIndent(p, "", "  ")
	if err != nil {
		err = fmt.Errorf("unable to marshall JSON for packages; %w", err)
		goto end
	}

	pwd, err = os.Getwd()
	if err != nil {
		err = fmt.Errorf("unable to get working directory for packages; %w", err)
		goto end
	}
	jf = fmt.Sprintf("%s/embeds/packages.json", pwd)

	// @see https://stackoverflow.com/a/9373342/102699
	jf = filepath.FromSlash(jf)

	err = ioutil.WriteFile(jf, b, PackageFilePerms)
	if err != nil {
		err = fmt.Errorf("unable to get working directory for packages; %w", err)
		goto end
	}
end:
	return err
}

// GetWorkingDir returns current working directory
func GetWorkingDir() (workdir string, err error) {
	workdir, err = os.Getwd()
	if err != nil {
		err = fmt.Errorf("failed to get current working directory; %w", err)
	}
	return workdir, err
}

// GetHomeDir returns the Liquibase home directory
func GetHomeDir() (homedir string, err error) {
	var out []byte
	var loc string
	var fi os.FileInfo
	var link string

	if _, ok := os.LookupEnv("LIQUIBASE_HOME"); ok {
		homedir = os.Getenv("LIQUIBASE_HOME")
		goto end
	}

	out, err = exec.Command("which", "liquibase").CombinedOutput()
	if err != nil {
		// @see https://blog.golang.org/go1.13-errors#TOC_3.3.
		err = fmt.Errorf("unable to locate Liquibase: %w", err)
		goto end
	}

	// Determine if Command is Symlink
	loc = strings.TrimRight(string(out), "\n")
	fi, err = os.Lstat(loc)
	if err != nil {
		err = fmt.Errorf("cannot stat %s: %w", loc, err)
		goto end
	}

	if fi.Mode()&os.ModeSymlink == 0 {
		// Not Symlink
		homedir, _ = filepath.Split(loc)
		goto end
	}

	link, err = os.Readlink(loc)
	if err != nil {
		err = fmt.Errorf("cannot read %s: %w", loc, err)
		goto end
	}

	// Is Symlink
	homedir, _ = filepath.Split(link)

end:
	if !strings.HasSuffix(homedir, "/") {
		homedir = homedir + "/"
	}
	return homedir, err

}

func (ctx *Context) GetPackageAndVersion(name string) (p Package, v Version, err error) {
	var parts []string
	if !strings.Contains(name, "@") {
		p = ctx.GetPackageByName(name)
		v = p.GetLatestVersion()
		goto end
	}
	parts = strings.Split(name, "@")
	if parts[0] == "" {
		err = fmt.Errorf("no package name provided in '%s'", parts[0])
		goto end
	}
	p = ctx.GetPackageByName(parts[0])
	if parts[1] == "" {
		err = fmt.Errorf("no VersionNumber name provided in '%s'", parts[1])
		goto end
	}
	v = p.GetVersion(parts[1])
	if v.Tag == "" {
		err = fmt.Errorf("VersionNumber %s not available for package %s",
			parts[1],
			name)
		goto end
	}
end:
	return p, v, err

}

func (ctx *Context) Error(params ...interface{}) {
	var msg string
	var err error
	err, ok := params[1].(error)
	if ok {
		goto end
	}
	msg, ok = params[1].(string)
	if !ok {
		err = fmt.Errorf("error in parameters passed to ctx.Error(); first parameter must be of type `string` or `error`")
	}
	if len(params) == 1 {
		err = fmt.Errorf(msg)
		goto end
	}
	err = fmt.Errorf(msg, params[1:]...)
end:
	ctx.errors = append(ctx.errors, err)
}

// ShowUserError display any error to a user
// TODO Make the output generated more attractive
func ShowUserError(err error) {
	fmt.Printf("%s\n", err)

}
