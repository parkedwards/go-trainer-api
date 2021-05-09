package appointment

import (
	"fmt"
	"strconv"
	"time"

	"github.com/parkedwards/go-trainer-api/pkg/models"
	"github.com/parkedwards/go-trainer-api/pkg/utils"
)

type DBClient interface {
	Select(q *models.Query) []models.Appointment
	Insert(a *models.Appointment)
}

type AppointmentManager struct {
	dbc DBClient
}

func New(dbc DBClient) *AppointmentManager {
	return &AppointmentManager{
		dbc: dbc,
	}
}

func (a *AppointmentManager) GetAllAppointmentsByTrainerId(trainerId string) []models.Appointment {
	// normalize string url param => int64, which is how the trainer_id is stored in the "db"
	intTrainerId, _ := strconv.ParseInt(trainerId, 10, 64)
	query := &models.Query{TrainerId: intTrainerId}
	appointments := a.dbc.Select(query)
	return appointments
}

// Assumption:
// Client provides start/end bookends that are 30m apart (since that's how we provide them from the GET /availability call)
func (a *AppointmentManager) MakeAppointmentWithTrainer(details *models.Appointment) error {

	// converts user-provide timed strings into time.Time, so we can perform some checks on them
	timeStartsAt, _ := time.Parse(time.RFC3339, details.StartsAt)

	if !utils.IsDuringBusinessHours(timeStartsAt) {
		return fmt.Errorf("Our trainers are available from 8a - 5p PT, Monday through Friday")
	}

	appointments := a.dbc.Select(&models.Query{TrainerId: details.TrainerId})

	for _, a := range appointments {
		if a.StartsAt == details.StartsAt || a.EndsAt == details.EndsAt {
			return fmt.Errorf("This trainer is already booked during the time you requested.")
		}
	}

	a.dbc.Insert(details)

	return nil
}
