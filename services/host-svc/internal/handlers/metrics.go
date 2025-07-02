package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"pagemagic/host-svc/internal/services"

	"github.com/gorilla/mux"
)

type MetricsHandler struct {
	metricsSvc *services.MetricsService
}

func NewMetricsHandler(metricsSvc *services.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		metricsSvc: metricsSvc,
	}
}

func (h *MetricsHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]

	// Parse query parameters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start_date format", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -30) // Last 30 days
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end_date format", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now()
	}

	analytics, err := h.metricsSvc.GetAnalytics(r.Context(), siteID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

func (h *MetricsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]

	stats, err := h.metricsSvc.GetStats(r.Context(), siteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *MetricsHandler) PrometheusHandler() http.Handler {
	return h.metricsSvc.PrometheusHandler()
}
