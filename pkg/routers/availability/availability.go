package availability

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/parkedwards/go-trainer-api/pkg/models"
)

type AvailabilityManager interface {
	GetTrainerAvailabilityForDateRange(trainerId string, startsAt string, endsAt string) ([]models.Availability, error)
}

type AvailabilityRouter struct {
	availabilityMgr AvailabilityManager
}

// Manages the /availability route handlers
// Takes in the AvailabilityManager instance, which interfaces with the database
func New(am AvailabilityManager) *AvailabilityRouter {
	return &AvailabilityRouter{
		availabilityMgr: am,
	}
}

func (r *AvailabilityRouter) RegisterRoutes(c chi.Router) {
	c.Get("/availability/{trainerId}", r.getTrainerAvailability)
}

// getTrainerAvailability godoc
// @Summary Get Availability for Trainer
// @Description get all available time slots for trainer
// @ID get-availability-by-trainer-id
// @Produce  json
// @Param trainerId path string true "Trainer ID"
// @Param starts_at query string true "Starts At"
// @Param ends_at query string true "Ends At"
// @Success 200 {object} []models.Availability
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /availability/{trainerId} [get]
func (r *AvailabilityRouter) getTrainerAvailability(w http.ResponseWriter, req *http.Request) {
	trainerId := chi.URLParam(req, "trainerId")
	startsAt := req.URL.Query().Get("starts_at")
	endsAt := req.URL.Query().Get("ends_at")

	// If no starts_at OR ends_at is provided, we will return all avails from now -> 7 days
	if startsAt == "" || endsAt == "" {
		now := time.Now().Add(1 * time.Hour)
		startsAt = now.Format(time.RFC3339)
		endsAt = now.AddDate(0, 0, 7).Format(time.RFC3339)

		fmt.Println(startsAt)
		fmt.Println(endsAt)
	}

	availability, err := r.availabilityMgr.GetTrainerAvailabilityForDateRange(trainerId, startsAt, endsAt)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte(err.Error())))
		return
	}

	if len(availability) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write(([]byte("Your trainer is all booked up for those dates.")))
		return
	}

	result, _ := json.Marshal(availability)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
