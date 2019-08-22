package reservations

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	"github.com/portainer/portainer/api"
)

type reservationUpdatePayload struct {
	Cpu 	int
	Memory  int
}

func (payload *reservationUpdatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Cpu) {
		return portainer.Error("Invalid reservation cpu value")
	}
	if govalidator.IsNull(payload.Memory) {
		return portainer.Error("Invalid reservation memory value")
	}

	return nil
}

// PUT request on /api/reservation/:id
func (handler *Handler) reservationUpdate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	reservationID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid reservation identifier route variable", err}
	}

	var payload reservationUpdatePayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	reservation, err := handler.ReservationService.Reservation(portainer.ReservationID(reservationID))
	if err == portainer.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a reservation with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a reservation with the specified identifier inside the database", err}
	}

	reservation.Cpu = payload.Cpu
	reservation.Memory = payload.Memory

	err = handler.ReservationService.UpdateReservation(reservation.ID, reservation)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist reservation changes inside the database", err}
	}

	return response.JSON(w, reservation)
}
