package middlewares

import (
    "net/http"
)

// This middleware ensures that content is properly JSON-encoded.
//
var Json Middleware = func(inner http.Handler, route Route) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        // Set basic JSON headers.
        //
        w.Header().Set("Content-Type", "application/json;charset=UTF-8")

        inner.ServeHTTP(w, r)
    })
}
