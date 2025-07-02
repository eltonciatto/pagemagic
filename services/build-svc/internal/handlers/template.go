package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"pagemagic/build-svc/internal/services"

	"github.com/gorilla/mux"
)

type TemplateHandler struct {
	templateSvc *services.TemplateService
}

func NewTemplateHandler(templateSvc *services.TemplateService) *TemplateHandler {
	return &TemplateHandler{
		templateSvc: templateSvc,
	}
}

func (h *TemplateHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	category := r.URL.Query().Get("category")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
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

	templates, err := h.templateSvc.ListTemplates(r.Context(), category, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"templates": templates,
		"total":     len(templates),
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *TemplateHandler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateID := vars["id"]

	template, err := h.templateSvc.GetTemplate(r.Context(), templateID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

func (h *TemplateHandler) GetFeaturedTemplates(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	templates, err := h.templateSvc.GetFeaturedTemplates(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"templates": templates,
		"total":     len(templates),
	})
}

func (h *TemplateHandler) SearchTemplates(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	framework := r.URL.Query().Get("framework")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
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

	templates, err := h.templateSvc.SearchTemplates(r.Context(), query, category, framework, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"templates": templates,
		"total":     len(templates),
		"query":     query,
		"filters": map[string]string{
			"category":  category,
			"framework": framework,
		},
	})
}

func (h *TemplateHandler) DownloadTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateID := vars["id"]

	// Incrementar contador de downloads
	err := h.templateSvc.IncrementDownloads(r.Context(), templateID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Obter template completo
	template, err := h.templateSvc.GetTemplate(r.Context(), templateID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"template": template,
		"message":  "Template downloaded successfully",
	})
}

func (h *TemplateHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.templateSvc.GetCategories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"categories": categories,
	})
}

func (h *TemplateHandler) GetFrameworks(w http.ResponseWriter, r *http.Request) {
	frameworks, err := h.templateSvc.GetFrameworks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"frameworks": frameworks,
	})
}
