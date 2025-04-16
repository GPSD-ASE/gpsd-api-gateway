package handlers

import (
	"fmt"
	"net/http"

	"gpsd-api-gateway/internal/gateway/pkg/config"
)

func getDecisionEngineBaseURL(cc *config.Config) string {
	return fmt.Sprintf(
		"http://%s:%s",
		cc.DecisionEngineHost,
		cc.DecisionEnginePort,
	)
}

func getEscalationMgmtBaseURL(cc *config.Config) string {
	return fmt.Sprintf(
		"http://%s:%s",
		cc.EscalationMgmtHost,
		cc.EscalationMgmtPort,
	)
}

func (h *Handler) PostDecisionHandler(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getDecisionEngineBaseURL(h.Config)+"/decision/incident", nil)
}

func (h *Handler) PostIncidentAnalysis(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, getEscalationMgmtBaseURL(h.Config)+"/incident-analysis/", nil)
}
