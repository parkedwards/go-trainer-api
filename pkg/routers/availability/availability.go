package availability

import (
	"encoding/json"
	"net/http"

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

func (r *AvailabilityRouter) getTrainerAvailability(w http.ResponseWriter, req *http.Request) {
	trainerId := chi.URLParam(req, "trainerId")
	startsAt := req.URL.Query().Get("starts_at")
	endsAt := req.URL.Query().Get("ends_at")

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
