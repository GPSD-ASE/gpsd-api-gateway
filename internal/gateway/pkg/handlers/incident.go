package handlers

import (
	"fmt"
	"net/http"

	"gpsd-api-gateway/internal/gateway/pkg/config"

	"github.com/gorilla/mux"
)

func getIncidentMgmtBaseURL(cc *config.Config) string {
	return fmt.Sprintf(
		"http://%s:%s/api",
		cc.IncidentMgmtHost,
		cc.IncidentMgmtPort,
	)
}

func (h *Handler) GetAllIncidentsHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getIncidentMgmtBaseURL(h.Config)+"/incidents", nil)
}

func (h *Handler) PostIncidentHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getIncidentMgmtBaseURL(h.Config)+"/incidents", nil)
}

func (h *Handler) GetIncidentByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ForwardRequest(w, r, getIncidentMgmtBaseURL(h.Config)+"/incidents/"+id, nil)
}

func (h *Handler) DeleteIncidentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ForwardRequest(w, r, getIncidentMgmtBaseURL(h.Config)+"/incidents/"+id, nil)
}

func (h *Handler) UpdateIncidentStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	status := vars["status"]
	ForwardRequest(w, r, getIncidentMgmtBaseURL(h.Config)+"/incidents/"+id+"/status/"+status, nil)
}
