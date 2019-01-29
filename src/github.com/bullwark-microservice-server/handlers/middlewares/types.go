package middlewares

import "net/http"

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
    Middlewares []Middleware
}

type Middleware func(inner http.Handler, route Route) http.Handler

type Routes []Route
