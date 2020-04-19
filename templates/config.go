package cuefig

// Name: {{ .CONFIG.Name }}

import (
	"fmt"

	"cuelang.org/go/cue"
	// "cuelang.org/go/cue/build"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"

	"github.com/hofstadter-io/hof/lib/util"
)

var {{ .CONFIG.ConfigName }} map[string]interface{}

const entrypoint = "{{ .CONFIG.Entrypoint }}"

func Init() {
	fmt.Println("reading cuefig:", entrypoint)
}

func LoadDefault() (map[string]interface{}, error) {
	return LoadConfig(entrypoint)
}

func LoadConfig(entry string) (map[string]interface{}, error) {

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
		// R.CueInstances = append(R.CueInstances, I)

		// Get top level value from cuelang
		V := I.Value()
		// R.TopLevelValues = append(R.TopLevelValues, V)

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

		l, b := V.Label()
		fmt.Println("Cfg Field:", l, b)
		// R.TopLevelStructs = append(R.TopLevelStructs, S)
		iter := S.Fields()
		for iter.Next() {

			label := iter.Label()
			value := iter.Value()
			fmt.Println("  -", label, value)

			// Now decode
			val := map[string]interface{}{}
			err = value.Decode(&val)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			cfg[label] = val

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
