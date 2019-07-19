package reservations

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	"github.com/portainer/portainer/api"
)

type reservationCreatePayload struct {
	Name string
	Cpu int
	Memory int
}

func (payload *reservationCreatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Name) {
		return portainer.Error("Invalid reservation name")
	}
	if govalidator.IsNull(payload.Cpu) {
		return portainer.Error("Invalid reservation cpu value")
	}
	if govalidator.IsNull(payload.Memory) {
		return portainer.Error("Invalid reservation memory value")
	}
	return nil
}

// POST request on /api/reservations
func (handler *Handler) reservationCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload reservationCreatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	reservations, err := handler.ReservationService.Reservations()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve reservations from the database", err}
	}

	for _, reservation := range reservations {
		if reservation.Name == payload.Name && !reservation.Revoked {
			return &httperror.HandlerError{http.StatusConflict, "This name is already associated to a reservation", portainer.ErrReservationAlreadyExists}
		}
	}

	reservation := &portainer.Reservation{
		Name: payload.Name,
		Cpu: payload.Cpu,
		Memory: payload.Memory
	}

	err = handler.ReservationService.CreateReservation(reservation)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist the reservation inside the database", err}
	}

	return response.JSON(w, reservation)
}
