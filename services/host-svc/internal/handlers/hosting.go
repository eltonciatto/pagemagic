package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"pagemagic/host-svc/internal/models"
	"pagemagic/host-svc/internal/services"

	"github.com/gorilla/mux"
)

type HostingHandler struct {
	hostSvc   *services.HostingService
	deploySvc *services.DeploymentService
	domainSvc *services.DomainService
}

func NewHostingHandler(hostSvc *services.HostingService, deploySvc *services.DeploymentService, domainSvc *services.DomainService) *HostingHandler {
	return &HostingHandler{
		hostSvc:   hostSvc,
		deploySvc: deploySvc,
		domainSvc: domainSvc,
	}
}

func (h *HostingHandler) CreateSite(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	site, err := h.hostSvc.CreateSite(r.Context(), userID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(site)
}

func (h *HostingHandler) ListSites(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	sites, err := h.hostSvc.ListSites(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sites": sites,
		"total": len(sites),
	})
}

func (h *HostingHandler) GetSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	site, err := h.hostSvc.GetSite(r.Context(), siteID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(site)
}

func (h *HostingHandler) UpdateSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	var req models.UpdateSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	site, err := h.hostSvc.UpdateSite(r.Context(), siteID, userID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(site)
}

func (h *HostingHandler) DeleteSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	err := h.hostSvc.DeleteSite(r.Context(), siteID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *HostingHandler) GetSiteStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteID := vars["id"]

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	site, err := h.hostSvc.GetSite(r.Context(), siteID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Obter deployments recentes
	deployments, err := h.deploySvc.ListDeployments(r.Context(), siteID, 5, 0)
	if err != nil {
		deployments = []*models.Deployment{}
	}

	// Obter domínios
	domains, err := h.domainSvc.ListDomains(r.Context(), siteID)
	if err != nil {
		domains = []*models.DomainConfig{}
	}

	status := map[string]interface{}{
		"site":        site,
		"deployments": deployments,
		"domains":     domains,
		"last_check":  time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (h *HostingHandler) ServeStaticFiles() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extrair domínio do host
		host := r.Host

		// Buscar site pelo domínio
		site, err := h.hostSvc.GetSiteByDomain(r.Context(), host)
		if err != nil {
			h.serve404(w, r)
			return
		}

		// Verificar se site está ativo
		if site.Status != "active" {
			h.serveMaintenancePage(w, r, site)
			return
		}

		// Aplicar configurações de segurança/headers
		h.applySecurityHeaders(w, site)

		// Servir arquivo
		h.serveFile(w, r, site)
	})
}

func (h *HostingHandler) serveFile(w http.ResponseWriter, r *http.Request, site *models.Site) {
	// Implementar lógica de servir arquivos do storage
	// Por enquanto, servir página padrão

	path := r.URL.Path
	if path == "/" || path == "" {
		path = "/index.html"
	}

	// Simular conteúdo
	content := `<!DOCTYPE html>
<html>
<head>
    <title>` + site.Domain + `</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
    <h1>Welcome to ` + site.Domain + `</h1>
    <p>This site is hosted by PageMagic</p>
    <p>Site ID: ` + site.ID + `</p>
    <p>Last Deploy: ` + site.LastDeploy.Format("2006-01-02 15:04:05") + `</p>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(content))
}

func (h *HostingHandler) serve404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	content := `<!DOCTYPE html>
<html>
<head>
    <title>Site Not Found</title>
</head>
<body>
    <h1>404 - Site Not Found</h1>
    <p>The requested site could not be found.</p>
</body>
</html>`
	w.Write([]byte(content))
}

func (h *HostingHandler) serveMaintenancePage(w http.ResponseWriter, r *http.Request, site *models.Site) {
	w.WriteHeader(http.StatusServiceUnavailable)
	content := `<!DOCTYPE html>
<html>
<head>
    <title>Site Under Maintenance</title>
</head>
<body>
    <h1>Site Under Maintenance</h1>
    <p>This site is currently being updated. Please check back later.</p>
</body>
</html>`
	w.Write([]byte(content))
}

func (h *HostingHandler) applySecurityHeaders(w http.ResponseWriter, site *models.Site) {
	// Aplicar headers de segurança configurados
	for key, value := range site.Config.Headers {
		w.Header().Set(key, value)
	}

	// Headers padrão de segurança
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
}
