package reservation

import (
	"time"

	"github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/bolt/internal"

	"github.com/boltdb/bolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "reservations"
)

// Service represents a service for managing endpoint data.
type Service struct {
	db *bolt.DB
}

// NewService creates a new instance of a service.
func NewService(db *bolt.DB) (*Service, error) {
	err := internal.CreateBucket(db, BucketName)
	if err != nil {
		return nil, err
	}

	return &Service{
		db: db,
	}, nil
}

// Reservation returns a reservation by ID
func (service *Service) Reservation(ID portainer.ReservationID) (*portainer.Reservation, error) {
	var reservation portainer.Reservation
	identifier := internal.Itob(int(ID))

	err := internal.GetObject(service.db, BucketName, identifier, &reservation)
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

// Reservations return an array containing all the reservations.
func (service *Service) Reservations() ([]portainer.Reservation, error) {
	var reservations = make([]portainer.Reservation, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var reservation portainer.Reservation
			err := internal.UnmarshalObject(v, &reservation)
			if err != nil {
				return err
			}
			if !reservation.Revoked {
				reservations = append(reservations, reservation)
			}
		}

		return nil
	})

	return reservations, err
}

// CreateReservation creates a new reservation.
func (service *Service) CreateReservation(reservation *portainer.Reservation) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		id, _ := bucket.NextSequence()
		reservation.ID = portainer.ReservationID(id)

		reservation.CreatedAt = time.Now().UnixNano() / int64(time.Millisecond)
		reservation.Revoked = false

		data, err := internal.MarshalObject(reservation)
		if err != nil {
			return err
		}

		return bucket.Put(internal.Itob(int(reservation.ID)), data)
	})
}

// UpdateReservation saves a reservation.
func (service *Service) UpdateReservation(ID portainer.ReservationID, reservation *portainer.Reservation) error {
	identifier := internal.Itob(int(ID))
	return internal.UpdateObject(service.db, BucketName, identifier, reservation)
}

// DeleteReservation deletes a reservation.
func (service *Service) DeleteReservation(ID portainer.ReservationID) error {

	reservation, err := Reservation(ID)
	if err != nil {
		return err
	}

	reservation.Revoked = true

	return UpdateReservation(ID, reservation)
}
