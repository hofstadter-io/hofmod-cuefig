package cuefig

// Name: {{ .CONFIG.Name }}

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	// "cuelang.org/go/cue/build"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	"github.com/kirsle/configdir"

	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

const (
	{{ .CONFIG.ConfigName }}Entrypoint = "{{ .CONFIG.Entrypoint }}"
	{{ .CONFIG.ConfigName }}Workpath   = "{{ .CONFIG.Workpath }}"
	{{ .CONFIG.ConfigName }}Location    = "{{ .CONFIG.Location }}"
)

func Load{{ .CONFIG.ConfigName }}Default() (cue.Value, error) {
	// default const value
	workpath, err := calc{{ .CONFIG.ConfigName }}Workpath()
	if err != nil {
		return cue.Value{}, err
	}
	return Load{{ .CONFIG.ConfigName }}Config(workpath, {{ .CONFIG.ConfigName }}Entrypoint)
}

func Save{{ .CONFIG.ConfigName }}Default(val cue.Value) (error) {
	// default const value
	workpath, err := calc{{ .CONFIG.ConfigName }}Workpath()
	if err != nil {
		return err
	}
	return Save{{ .CONFIG.ConfigName }}Config(workpath, {{ .CONFIG.ConfigName }}Entrypoint, val)
}

func calc{{ .CONFIG.ConfigName }}Workpath() (string, error) {
	workpath := {{ .CONFIG.ConfigName }}Workpath
	switch {{ .CONFIG.ConfigName }}Location {

		case "home":
			dir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			workpath = filepath.Join(dir, workpath)

		case "user":
			dir, err := os.UserConfigDir()
			if err != nil {
				return "", err
			}
			workpath = filepath.Join(dir, workpath)

		case "cache":
			dir, err := os.UserCacheDir()
			if err != nil {
				return "", err
			}
			workpath = filepath.Join(dir, workpath)

		case "system":
			// TODO, add some preference for well known directories here?
			workpath = filepath.Join(configdir.SystemConfig()[0], workpath)
	}

	return workpath, nil
}

func Load{{ .CONFIG.ConfigName }}Config(workpath, entrypoint string) (val cue.Value, err error) {

	fpath := filepath.Join(workpath, entrypoint)

	// possibly, check for workpath
	if workpath != "" {
		_, err = os.Lstat(workpath)
		if err != nil {
			if _, ok := err.(*os.PathError); !ok && ( strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file") ) {
				// error is worse than non-existant
				return val, err
			}
			// otherwise, does not exist, so we should init?
			// XXX want to let applications decide how to handle this
			return val, err
		}
	}

	// check for entrypoint
	_, err = os.Lstat(fpath)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && ( strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file") ) {
			// error is worse than non-existant
			return val, err
		}
		// otherwise, does not exist, so we should init?
		// XXX want to let applications decide how to handle this
		return val, err
	}

	var errs []error

	CueRT := &cue.Runtime{}

	loadConfig := &load.Config{
		Dir: workpath,
		Package: "{{ .CONFIG.Defaults.Package }}",
	}

	BIS := load.Instances([]string{entrypoint}, loadConfig)
	for _, bi := range BIS {

		if bi.Err != nil {
			// fmt.Println("BI ERR", bi.Err, bi.Incomplete, bi.DepsErrors)
		  es := errors.Errors(bi.Err)
			for _, e := range es {
				errs = append(errs, e.(error))
			}
			continue
		}

		// Build the Instance
		I, err := CueRT.Build(bi)
		if err != nil {
		  es := errors.Errors(err)
			// fmt.Println("BUILD ERR", es, I)
			for _, e := range es {
				errs = append(errs, e.(error))
			}
			continue
		}

		// Get top level value from cuelang and persist
		V := I.Value()
		val = V
	}

	if len(errs) > 0 {
		for _, e := range errs {
			cuetils.PrintCueError(e)
		}
		return val, fmt.Errorf("Errors while reading {{.CONFIG.ConfigName}} file: %q", fpath)
	}

	return val, nil
}

func Save{{ .CONFIG.ConfigName }}Config(workpath, entrypoint string, val cue.Value) (err error) {

	fpath := filepath.Join(workpath, entrypoint)

	// possibly, check for workpath
	if workpath != "" {
		_, err = os.Lstat(workpath)
		if err != nil {
			if _, ok := err.(*os.PathError); !ok && ( strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file") ) {
				// error is worse than non-existant
				return err
			}
			// otherwise, does not exist, so we should init
			err = yagu.Mkdir(workpath)
			if err != nil {
				return err
			}
		}
	}

	// check for entrypoint
	_, err = os.Lstat(fpath)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && ( strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file") ) {
			// error is worse than non-existant
			return err
		}
		// otherwise, does not exist, so we should init?
		err = yagu.Mkdir(filepath.Dir(fpath))
		if err != nil {
			return err
		}
	}

	// get string version of value
	bytes, err := format.Node(val.Syntax())
	if err != nil {
		return err
	}

	// TODO, temp print
	str := string(bytes)
	fmt.Println(str)

	// write the file
	err = ioutil.WriteFile(fpath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
