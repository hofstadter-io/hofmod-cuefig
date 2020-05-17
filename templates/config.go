package cuefig

// Name: {{ .CONFIG.Name }}

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	// "cuelang.org/go/cue/build"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"

	"github.com/hofstadter-io/hof/lib/util"
)

const (
	{{ .CONFIG.ConfigName }}Entrypoint = "{{ .CONFIG.Entrypoint }}"
	{{ .CONFIG.ConfigName }}Workpath   = "{{ .CONFIG.Workpath }}"
)

func Load{{ .CONFIG.ConfigName }}Default(cfg interface{}) (cue.Value, error) {
	return Load{{ .CONFIG.ConfigName }}Config({{.CONFIG.ConfigName}}Workpath, {{ .CONFIG.ConfigName }}Entrypoint, cfg)
}

func Load{{ .CONFIG.ConfigName }}Config(workpath, entrypoint string, cfg interface{}) (val cue.Value, err error) {

	// TODO Fallback order: local / user / global
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
	BIS := load.Instances([]string{fpath}, nil)
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

		// Get top level value from cuelang
		V := I.Value()

		err = V.Decode(&cfg)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		val = V

	}

	if len(errs) > 0 {
		for _, e := range errs {
			util.PrintCueError(e)
		}
		return val, fmt.Errorf("Errors while reading {{.CONFIG.ConfigName}} file: %q", fpath)
	}

	return val, nil
}
