package cli

import (
    "flag"
)

type ClientConfig struct {
    GenerateServiceId  *bool
    ConfigPath         *string
}

func (c ClientConfig) Configure(configMap map[ConfigType]Config) {

    configMap[ClientConfigType] = ClientConfig{
        GenerateServiceId :   flag.Bool("generate-service-ids", false, "Set this to true to enable automatic generation of service ids"),
        ConfigPath        :   flag.String("config", "services.yml", "Path to CLI config, if running in CLI mode"),
    }
}
