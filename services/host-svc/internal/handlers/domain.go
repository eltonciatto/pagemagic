package handlers

import (
	"encoding/json"
	"net/http"

	"pagemagic/host-svc/internal/services"

	"github.com/gorilla/mux"
)

type DomainHandler struct {
	domainSvc *services.DomainService
}

func NewDomainHandler(domainSvc *services.DomainService) *DomainHandler {
	return &DomainHandler{
		domainSvc: domainSvc,
	}
}

func (h *DomainHandler) AddDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]

	var req struct {
		Domain string `json:"domain"`
		Type   string `json:"type"` // "subdomain" ou "custom"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	domain, err := h.domainSvc.AddDomain(r.Context(), siteID, req.Domain, req.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain)
}

func (h *DomainHandler) ListDomains(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]

	domains, err := h.domainSvc.ListDomains(r.Context(), siteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"domains": domains,
		"total":   len(domains),
	})
}

func (h *DomainHandler) UpdateDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]
	domainID := vars["domainId"]

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	domain, err := h.domainSvc.UpdateDomain(r.Context(), domainID, siteID, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain)
}

func (h *DomainHandler) DeleteDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]
	domainID := vars["domainId"]

	err := h.domainSvc.DeleteDomain(r.Context(), domainID, siteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DomainHandler) VerifyDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]
	domainID := vars["domainId"]

	err := h.domainSvc.VerifyDomain(r.Context(), domainID, siteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Domain verification started",
	})
}
