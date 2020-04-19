package schema

import (
  "strings"
)

Config :: {
  // binary name, often last part of Package
  Name: string
  configName: strings.ToCamel(Name)
  ConfigName: strings.ToTitle(Name)

  Entrypoint: string

  ConfigSchema: _
}
