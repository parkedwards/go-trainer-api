package availability

import (
	"fmt"
	"strconv"
	"time"

	"github.com/parkedwards/go-trainer-api/pkg/models"
	"github.com/parkedwards/go-trainer-api/pkg/utils"
)

type DBClient interface {
	Select(q *models.Query) []models.Appointment
}

type AvailabilityManager struct {
	dbc DBClient
}

func New(dbc DBClient) *AvailabilityManager {
	return &AvailabilityManager{
		dbc: dbc,
	}
}

func (a *AvailabilityManager) GetTrainerAvailabilityForDateRange(trainerId string, startsAt string, endsAt string) ([]models.Availability, error) {
	now := time.Now()
	timeStartsAt, _ := time.Parse(time.RFC3339, startsAt)
	timeEndsAt, _ := time.Parse(time.RFC3339, endsAt)

	if timeEndsAt.Before(timeStartsAt) || timeEndsAt.Equal(timeStartsAt) {
		return nil, fmt.Errorf("Your start time must be before your end time.")
	}

	if timeEndsAt.Before(now) || timeStartsAt.Before(now) {
		return nil, fmt.Errorf("Your start time or end time is in the past.")
	}

	// convert url parameter (string) to int
	convertedTrainerId, _ := strconv.ParseInt(trainerId, 10, 64)

	// query "db" for all appointments for a given trainerId
	existingAppointments := a.dbc.Select(&models.Query{TrainerId: convertedTrainerId})
	appointmentMap := a.createAppointmentMap(existingAppointments)

	availableSlots := a.getAvailableTimeSlotsForRange(timeStartsAt, timeEndsAt, appointmentMap, convertedTrainerId)

	return availableSlots, nil
}

// Custom map for constant lookup of existing appointment times for a given trainerId
// eg. { "2021-05-03T08:00:00-08:00": true }
func (a *AvailabilityManager) createAppointmentMap(appointments []models.Appointment) map[string]bool {
	appointmentMap := make(map[string]bool)

	for _, a := range appointments {
		appointmentMap[a.StartsAt] = true
	}

	return appointmentMap
}

// Generates a continuous list of 30m slots from starts_at (rounded to nearest :30) -> ends_at
// Filters out:
// - non-business hours slots
// - occupied slots
func (a *AvailabilityManager) getAvailableTimeSlotsForRange(timeStartsAt time.Time, timeEndsAt time.Time, appointmentMap map[string]bool, trainerId int64) []models.Availability {

	increment_30m := 30 * time.Minute

	// rounding bookend times - every :30
	roundedStartsAt := timeStartsAt.Round(increment_30m)
	roundedEndsAt := timeEndsAt.Round(increment_30m)

	availability := []models.Availability{}

	pointer := roundedStartsAt
	for {
		if pointer.Equal(roundedEndsAt) {
			break
		}

		if utils.IsDuringBusinessHours(pointer) && appointmentMap[pointer.Format(time.RFC3339)] != true {
			av := models.Availability{}
			av.TrainerId = trainerId
			av.StartsAt = pointer.Format(time.RFC3339)
			av.EndsAt = pointer.Add(increment_30m).Format(time.RFC3339)
			availability = append(availability, av)
		}

		// advance the pointer +30 min
		pointer = pointer.Add(increment_30m)
	}

	return availability
}
