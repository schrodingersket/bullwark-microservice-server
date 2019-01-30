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

        allowedOrigins := handlers.AllowedOrigins([]string{"*"})
        allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
        allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})

        // Define API routes
        //
        router := handlers.CORS(allowedOrigins, allowedMethods,
            allowedHeaders)(routes.NewRouter())

        // Start server
        //
        fmt.Printf("Microservice server started on :%d\n", common.GetListenPort())
        if err := http.ListenAndServe(fmt.Sprintf(":%d", common.GetListenPort()), router); err != nil {
            panic(err)
        }
    }
}
