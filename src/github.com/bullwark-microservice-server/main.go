package main

import (
    "fmt"
    "github.com/gorilla/handlers"
    "net/http"

    "github.com/bullwark-microservice-server/client"
    "github.com/bullwark-microservice-server/common"
    "github.com/bullwark-microservice-server/routes"
)

func main() {

    // Parse command-line flags
    //
    common.ConfigureFromFlags()


    if common.IsClientMode() {

        // Run client registration
        //
        client.RunRegistrarClient()

    } else {

        // Define API routes
        //
        router := handlers.CORS()(routes.NewRouter())

        // Start server
        //
        fmt.Printf("Microservice server started on :%d\n", common.GetListenPort())
        if err := http.ListenAndServe(fmt.Sprintf(":%d", common.GetListenPort()), router); err != nil {
            panic(err)
        }
    }
}
