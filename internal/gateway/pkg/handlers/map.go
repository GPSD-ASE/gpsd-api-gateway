package handlers

import (
	"fmt"
	"net/http"

	"gpsd-api-gateway/internal/gateway/pkg/config"
)

func getMapMgmtBaseURL(cc *config.Config) string {
	return fmt.Sprintf(
		"http://%s:%s",
		cc.MapMgmtHost,
		cc.MapMgmtPort,
	)
}

func (h *Handler) GetZonesHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL(h.Config)+"/zones", nil)
}

func (h *Handler) GetRouteHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL(h.Config)+"/route", nil)
}

func (h *Handler) GetRoutingHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL(h.Config)+"/routing", nil)
}

func (h *Handler) PostEvacuationHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL(h.Config)+"/evacuation", nil)
}

func (h *Handler) TrafficHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL(h.Config)+"/traffic", nil)
}

func (h *Handler) GetSafezonesHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL(h.Config)+"/safezones", nil)
}

func (h *Handler) PostSafezonesHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL(h.Config)+"/safezones", nil)
}
