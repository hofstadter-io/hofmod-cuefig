package gen

import (

  hof "github.com/hofstadter-io/hof/schema"

  "github.com/hofstadter-io/hofmod-cuefig/schema"
)

#HofGenerator: hof.#HofGenerator & {
  Config: schema.#Config
  Outdir?: string

  In: {
    CONFIG: Config
  }

  Out: _OnceFiles

  PackageName: "github.com/hofstadter-io/hofmod-cuefig"

  // OnceFiles that are not repeatedly used, they are generated once for the whole generator input
  _OnceFiles: [...hof.#HofGeneratorFile] & [
    {
      TemplateName:  "config.go"
      Filepath:  "\(Outdir)/cuefig/\(In.CONFIG.Name).go"
    },
  ]

}
