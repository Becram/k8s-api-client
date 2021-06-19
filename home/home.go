package home

import (
	"net/http"

	"github.com/Becram/k8s-api-client/templates"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	templates.RenderTemplate(w, "home", nil)
}
