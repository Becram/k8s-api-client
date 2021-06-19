package main

import (
	"encoding/json"
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
	"github.com/gorilla/mux"
)

type Status struct {
	Deployment  string `json:"Name"`
	RestartedAt string `json:"RestartedAt"`
}

type Statuses []Status

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
		restartDeployment,
	},

	Route{
		"login",
		"get",
		"/login",
		login.LoginHandler,
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



func restartDeployment(w http.ResponseWriter, r *http.Request) {
	statuses := Statuses{
		Status{Deployment: r.PostFormValue("Name"), RestartedAt: k8s.DeploymentRestart("apps", r.PostFormValue("Name"))["kubectl.kubernetes.io/restartedAt"]},
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		panic(err)
	}

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

}
