package reservations

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/http/security"
)

// Handler is the HTTP handler used to handle reservation operations.
type Handler struct {
	*mux.Router
	ReservationService portainer.ReservationService
}

// NewHandler creates a handler to manage reservation operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/reservations",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.reservationCreate))).Methods(http.MethodPost)
	h.Handle("/reservations",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.reservationList))).Methods(http.MethodGet)
	h.Handle("/reservations/{id}",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.reservationUpdate))).Methods(http.MethodPut)
	h.Handle("/reservations/{id}",
		bouncer.AuthorizedAccess(httperror.LoggerHandler(h.reservationDelete))).Methods(http.MethodDelete)

	return h
}
