package handlers

import (
	"fmt"
	"net/http"

	"gpsd-api-gateway/internal/gateway/pkg/config"
)

func getMapMgmtBaseURL() string {
	return fmt.Sprintf(
		"http://%s:%s/api/v1",
		config.ApiGatewayConfig.MapMgmtHost,
		config.ApiGatewayConfig.MapMgmtPort,
	)
}

func ZonesHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL()+"/zones", nil)
}

func RoutingHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL()+"/routing", nil)
}

func EvacuationHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL()+"/evacuation", nil)
}

func TrafficHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getMapMgmtBaseURL()+"/traffic", nil)
}
