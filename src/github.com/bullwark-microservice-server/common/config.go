package common

import (
    "github.com/bullwark-microservice-server/common/cli"
    "flag"
)

var flagsParsed = false

var Configs  = make(map[cli.ConfigType]cli.Config)

func ConfigureFromFlags() {

    // First, parse flags and ensure we haven't already done so.
    //
    if flagsParsed {
        return
    } else {

        // Parse all CLI options
        //
        for _, opts := range cli.ConfigList {
            opts.Configure(Configs)
        }

        flag.Parse()
        flagsParsed = true
    }
}

func GetListenPort() int {

    coreConfig, ok := Configs[cli.CoreConfigType].(cli.CoreConfig)

    if ok {
        return *coreConfig.Port
    }

    return -1
}

func IsClientMode() bool {

    coreConfig, ok := Configs[cli.CoreConfigType].(cli.CoreConfig)

    if ok {
        return *coreConfig.Client
    }

    return false
}
