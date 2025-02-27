package handlers

import (
	"fmt"
	"net/http"

	"gpsd-api-gateway/internal/gateway/pkg/config"

	"github.com/gorilla/mux"
)

func getIncidentMgmtBaseURL() string {
	return fmt.Sprintf(
		"http://%s:%s/api",
		config.ApiGatewayConfig.IncidentMgmtHost,
		config.ApiGatewayConfig.IncidentMgmtPort,
	)
}

func GetAllIncidentsHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getIncidentMgmtBaseURL()+"/incidents", nil)
}

func CreateIncidentHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getIncidentMgmtBaseURL()+"/incidents", nil)
}

func GetIncidentByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ForwardRequest(w, r, getIncidentMgmtBaseURL()+"/incidents/"+id, nil)
}

func DeleteIncidentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ForwardRequest(w, r, getIncidentMgmtBaseURL()+"/incidents/"+id, nil)
}

func ChangeIncidentStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	status := vars["status"]
	ForwardRequest(w, r, getIncidentMgmtBaseURL()+"/incidents/"+id+"/status/"+status, nil)
}
