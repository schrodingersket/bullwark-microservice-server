package cli

type ConfigType int

const (
    CoreConfigType    ConfigType = iota
    RegistrarConfigType
    ClientConfigType
)
