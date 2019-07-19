package reservations

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
)

// GET request on /api/tags
func (handler *Handler) tagList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	tags, err := handler.ReservationService.Reservations()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve reservations from the database", err}
	}

	return response.JSON(w, tags)
}