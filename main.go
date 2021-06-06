package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"example.com/k8s-go-client/k8s"
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
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"restartDeployment",
		"POST",
		"/restart",
		restartDeployment,
	},
}

func restartDeployment(w http.ResponseWriter, r *http.Request) {

	// updatedStatus := &status{Deployment: "demo", RestartedAt: "2021-06-06T00:04:34+07:00"}
	// b, err := json.Marshal(updatedStatus)
	// if err != nil {
	// fmt.Fprintf(w, "formdata, %q", html.EscapeString(r.PostFormValue("Name")))

	// }
	// // return string(b)
	statuses := Statuses{
		Status{Deployment: r.PostFormValue("Name"), RestartedAt: k8s.DeploymentRestart("apps", r.PostFormValue("Name"))["kubectl.kubernetes.io/restartedAt"]},
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		panic(err)
	}
	// fmt.Fprintln(w, "Status:", statuses)
	// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// func handleRequests() {
// 	http.HandleFunc("/", homePage)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

func main() {

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

	// fmt.Printf("Deployment restarted at %s", k8s.DeploymentRestart("apps", "demo-deployment")["kubectl.kubernetes.io/restartedAt"])

}
