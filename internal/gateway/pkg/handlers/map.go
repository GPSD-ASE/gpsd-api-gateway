package handlers

import (
	"fmt"
	"net/http"

	"gpsd-api-gateway/internal/gateway/pkg/config"
)

func getMapMgmtBaseURL(cc *config.Config) string {
	return fmt.Sprintf(
		"http://%s:%s/api/v1",
		cc.MapMgmtHost,
		cc.MapMgmtPort,
	)
}

func (h *Handler) ZonesHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL(h.Config)+"/zones", nil)
}

func (h *Handler) RoutingHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL(h.Config)+"/routing", nil)
}

func (h *Handler) EvacuationHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL(h.Config)+"/evacuation", nil)
}

func (h *Handler) TrafficHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL(h.Config)+"/traffic", nil)
}
