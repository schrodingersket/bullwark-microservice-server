package client

import (
    "fmt"
    "os"
)

func RunRegistrarClient() {

    //var cliConfig = ParseClientConfig()
    //clientConfig, ok := common.Configs[cli.ClientConfigType].(cli.ClientConfig)
    //
    //if !ok {
    //    terminate("Unable to parse configuration.")
    //}
}

func terminate(message interface{}) {
    fmt.Println(message)
    fmt.Println("Exiting.")
    os.Exit(1)
}
