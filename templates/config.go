package cuefig

// Name: {{ .CONFIG.Name }}

import (
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue"
	// "cuelang.org/go/cue/build"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"

	"github.com/hofstadter-io/hof/lib/util"
)

var {{ .CONFIG.ConfigName }} map[string]interface{}

const {{ .CONFIG.ConfigName }}Filepath = "{{ .CONFIG.Entrypoint }}"

func LoadDefault() (map[string]interface{}, error) {
	return LoadConfig({{ .CONFIG.ConfigName }}Filepath)
}

func LoadConfig(entry string) (map[string]interface{}, error) {

	_, err := os.Lstat(DmaFilepath)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && ( strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file") ) {
			// error is worse than non-existant
			return nil, err
		}
		// otherwise, does not exist, so we should init
		return nil, nil
	}

	var errs []error
	cfg := map[string]interface{}{}

	CueRT := &cue.Runtime{}

	BIS := load.Instances([]string{entry}, nil)


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

		// Get top level struct from cuelang
		S, err := V.Struct()
		if err != nil {
			// fmt.Println("STRUCT ERR", err)
		  es := errors.Errors(err)
			for _, e := range es {
				errs = append(errs, e.(error))
			}
			continue
		}

		iter := S.Fields()
		for iter.Next() {

			label := iter.Label()
			value := iter.Value()

			// Now decode
			val := map[string]interface{}{}
			err = value.Decode(&val)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			cfg[label] = val

			/* XXX leave as example
			fmt.Println("  -", label, value)
			for attrKey, attrVal := range value.Attributes() {
				fmt.Println("  --", attrKey)
				for i := 0; i < 5; i++ {
					str, err := attrVal.String(i)
					if err != nil {
					  break
					}
					fmt.Println("  ---", str)
				}
			}
			*/

		}

	}

	if len(errs) > 0 {
		for _, e := range errs {
			util.PrintCueError(e)
		}
		return nil, fmt.Errorf("Errors while reading DMA config file")
	}

	{{ .CONFIG.ConfigName }} = cfg

	return cfg, nil
}
