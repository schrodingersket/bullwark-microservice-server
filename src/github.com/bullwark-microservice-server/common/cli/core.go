package cli

import (
    "flag"
)

type CoreConfig struct {
    Client         *bool
    Port           *int
    Verbose        *bool
    BinSavePath    *string
    BaseURL        *string
}

func (c CoreConfig) Configure(configMap map[ConfigType]Config) {

  configMap[CoreConfigType] = CoreConfig{
      Client           :   flag.Bool("client", false, "Run in CLI mode."),
      Port             :   flag.Int("port", 8000, "Port to run on."),
      Verbose          :   flag.Bool("verbose", false, "Turn on verbose logging"),
      BinSavePath      :   flag.String("bin-save-path", "bin/binaries", "File path to which service binaries should be saved."),
      BaseURL          :   flag.String("base-url", "/microservices", "Determines the base URL which the application serves at."),
  }

}
