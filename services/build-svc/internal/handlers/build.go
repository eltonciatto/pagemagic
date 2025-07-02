package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"pagemagic/build-svc/internal/models"
	"pagemagic/build-svc/internal/services"

	"github.com/gorilla/mux"
)

type BuildHandler struct {
	buildSvc *services.BuildService
}

func NewBuildHandler(buildSvc *services.BuildService) *BuildHandler {
	return &BuildHandler{
		buildSvc: buildSvc,
	}
}

func (h *BuildHandler) CreateBuildJob(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBuildJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	job, err := h.buildSvc.CreateBuildJob(r.Context(), userID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.BuildJobResponse{Job: job})
}

func (h *BuildHandler) GetBuildJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobId"]

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	job, err := h.buildSvc.GetBuildJob(r.Context(), jobID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := models.BuildJobResponse{Job: job}
	if job.Status == "completed" {
		response.Output = &job.Output
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *BuildHandler) ListBuildJobs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["siteId"]

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	jobs, err := h.buildSvc.ListBuildJobs(r.Context(), siteID, userID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"jobs":   jobs,
		"total":  len(jobs),
		"limit":  limit,
		"offset": offset,
	})
}

func (h *BuildHandler) GetBuildLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobId"]

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	job, err := h.buildSvc.GetBuildJob(r.Context(), jobID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"job_id": jobID,
		"logs":   job.Logs,
		"status": job.Status,
	})
}

func (h *BuildHandler) CancelBuildJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobId"]

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	// Verificar se o job existe e pertence ao usuário
	job, err := h.buildSvc.GetBuildJob(r.Context(), jobID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Verificar se o job pode ser cancelado
	if job.Status == "completed" || job.Status == "failed" {
		http.Error(w, "Job cannot be cancelled", http.StatusBadRequest)
		return
	}

	// Cancelar o job
	err = h.buildSvc.UpdateBuildJobStatus(r.Context(), jobID, "cancelled", "Cancelled by user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Build job cancelled successfully",
	})
}

func (h *BuildHandler) RetryBuildJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobId"]

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	// Obter job original
	originalJob, err := h.buildSvc.GetBuildJob(r.Context(), jobID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Criar novo job baseado no original
	req := &models.CreateBuildJobRequest{
		SiteID:      originalJob.SiteID,
		Type:        originalJob.Type,
		SourceType:  originalJob.SourceType,
		SourceData:  originalJob.SourceData,
		BuildConfig: originalJob.BuildConfig,
	}

	newJob, err := h.buildSvc.CreateBuildJob(r.Context(), userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.BuildJobResponse{Job: newJob})
}
