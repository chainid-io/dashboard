package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/chainid-io/dashboard"
	httperror "github.com/chainid-io/dashboard/http/error"
	"github.com/chainid-io/dashboard/http/security"
)

// TemplatesHandler represents an HTTP API handler for managing templates.
type TemplatesHandler struct {
	*mux.Router
	Logger          *log.Logger
	SettingsService chainid.SettingsService
}

const (
	containerTemplatesURLLinuxServerIo = "https://tools.linuxserver.io/chainid.json"
)

// NewTemplatesHandler returns a new instance of TemplatesHandler.
func NewTemplatesHandler(bouncer *security.RequestBouncer) *TemplatesHandler {
	h := &TemplatesHandler{
		Router: mux.NewRouter(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	h.Handle("/templates",
		bouncer.AuthenticatedAccess(http.HandlerFunc(h.handleGetTemplates))).Methods(http.MethodGet)
	return h
}

// handleGetTemplates handles GET requests on /templates?key=<key>
func (handler *TemplatesHandler) handleGetTemplates(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if key == "" {
		httperror.WriteErrorResponse(w, ErrInvalidQueryFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	var templatesURL string
	switch key {
	case "containers":
		settings, err := handler.SettingsService.Settings()
		if err != nil {
			httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
			return
		}
		templatesURL = settings.TemplatesURL
	case "linuxserver.io":
		templatesURL = containerTemplatesURLLinuxServerIo
	default:
		httperror.WriteErrorResponse(w, ErrInvalidQueryFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	resp, err := http.Get(templatesURL)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
