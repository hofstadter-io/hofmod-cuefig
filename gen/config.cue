package gen

import (

  hof "github.com/hofstadter-io/hof/schema"

  "github.com/hofstadter-io/hofmod-cuefig/schema"
)

HofGenerator :: hof.HofGenerator & {
  Config: schema.Config
  Outdir?: string

  In: {
    CONFIG: Config
  }

  Out: _OnceFiles

  PackageName: "github.com/hofstadter-io/hofmod-cuefig"

  // Files that are not repeatedly used, they are generated once for the whole CLI
  _OnceFiles: [...hof.HofGeneratorFile] & [
    {
      TemplateName:  "config.go"
      Filepath:  "\(Outdir)/lib/cuefig/\(In.CONFIG.Name).go"
    },
  ]

}


