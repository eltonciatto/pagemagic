package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"pagemagic/host-svc/internal/models"
	"pagemagic/host-svc/internal/services"

	"github.com/gorilla/mux"
)

type DeploymentHandler struct {
	deploySvc *services.DeploymentService
	cacheSvc  *services.CacheService
}

func NewDeploymentHandler(deploySvc *services.DeploymentService, cacheSvc *services.CacheService) *DeploymentHandler {
	return &DeploymentHandler{
		deploySvc: deploySvc,
		cacheSvc:  cacheSvc,
	}
}

func (h *DeploymentHandler) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]

	var req struct {
		BuildID string                  `json:"build_id"`
		Files   []models.DeploymentFile `json:"files"`
		Config  models.SiteConfig       `json:"config"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	deployment, err := h.deploySvc.CreateDeployment(r.Context(), siteID, req.BuildID, req.Files, req.Config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployment)
}

func (h *DeploymentHandler) ListDeployments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	deployments, err := h.deploySvc.ListDeployments(r.Context(), siteID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deployments": deployments,
		"total":       len(deployments),
		"limit":       limit,
		"offset":      offset,
	})
}

func (h *DeploymentHandler) GetDeployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]
	deploymentID := vars["deployId"]

	deployment, err := h.deploySvc.GetDeployment(r.Context(), deploymentID, siteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployment)
}

func (h *DeploymentHandler) RollbackDeployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]
	deploymentID := vars["deployId"]

	err := h.deploySvc.RollbackToDeployment(r.Context(), siteID, deploymentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Rollback completed successfully",
	})
}

func (h *DeploymentHandler) PurgeCache(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]

	err := h.cacheSvc.PurgeCache(r.Context(), siteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Cache purged successfully",
	})
}

func (h *DeploymentHandler) GetCacheStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]

	status, err := h.cacheSvc.GetCacheStatus(r.Context(), siteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
