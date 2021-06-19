package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/Becram/k8s-api-client/app"
	"github.com/Becram/k8s-api-client/callback"
	"github.com/Becram/k8s-api-client/home"
	"github.com/Becram/k8s-api-client/k8s"
	"github.com/Becram/k8s-api-client/login"
	"github.com/Becram/k8s-api-client/logout"
	"github.com/Becram/k8s-api-client/middlewares"
	"github.com/Becram/k8s-api-client/user"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		home.HomeHandler,
	},
	Route{
		"restartDeployment",
		"POST",
		"/restart",
		k8s.RestartDeployment,
	},

	Route{
		"login",
		"get",
		"/login",
		login.LoginHandler,
	},

	Route{
		"callback",
		"get",
		"/callback",
		callback.CallbackHandler,
	},

	Route{
		"logout",
		"get",
		"/logout",
		logout.LogoutHandler,
	},
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {

	app.Init()
	router := NewRouter()
	router.Handle("/user", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(user.UserHandler)),
	))
	log.Fatal(http.ListenAndServe(":8080", router))

}
