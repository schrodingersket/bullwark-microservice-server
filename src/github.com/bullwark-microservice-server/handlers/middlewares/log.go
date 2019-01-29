package middlewares

import (
    "log"
    "net/http"
    "time"
)

// This middleware logs every request and time the request was made.
//
var Logger Middleware = func(inner http.Handler, route Route) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        inner.ServeHTTP(w, r)

        log.Printf(
            "%s\t%s\t%s\t%s",
            r.Method,
            r.RequestURI,
            route.Name,
            time.Since(start),
        )
    })
}