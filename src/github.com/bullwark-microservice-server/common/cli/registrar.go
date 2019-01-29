package cli

import (
    "flag"
)

type RegistrarConfig struct {
    Host     *string
    Port     *int
    Scheme   *string
    BaseURL  *string
}

func (c RegistrarConfig) Configure(configMap map[ConfigType]Config) {

    configMap[RegistrarConfigType] = RegistrarConfig{
        Host:        flag.String("registrar-host", "127.0.0.1", "Registrar host"),
        Port:        flag.Int("registrar-port", 8000, "Registrar port"),
        Scheme:      flag.String("registrar-scheme", "http", "Registrar scheme"),
        BaseURL:     flag.String("registrar-base-url", "", "Registrar base URL path"),
    }

}
