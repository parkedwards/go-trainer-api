package appointment

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/parkedwards/go-trainer-api/pkg/models"
)

type AppointmentManager interface {
	GetAllAppointmentsByTrainerId(trainerId string) []models.Appointment
	MakeAppointmentWithTrainer(appointmentDetails *models.Appointment) error
}

type AppointmentRouter struct {
	appointmentsManager AppointmentManager
}

// Manages the /appointment route handlers
// Takes in the AppointmentManager instance, which interfaces with the database
func New(am AppointmentManager) *AppointmentRouter {
	return &AppointmentRouter{
		appointmentsManager: am,
	}
}

func (r *AppointmentRouter) RegisterRoutes(c chi.Router) {
	c.Get("/appointment/{trainerId}", r.getTrainerAppointments)
	c.Post("/appointment", r.makeAppointment)
}

// getTrainerAppointments godoc
// @Summary Get Trainer's Appointments
// @Description Gets all existing appointments for a {trainerId}
// @ID get-appointments-by-trainer-id
// @Produce  json
// @Param trainerId path string true "Trainer ID"
// @Success 200 {object} []models.Appointment
// @Router /appointment/{trainerId} [get]
func (r *AppointmentRouter) getTrainerAppointments(w http.ResponseWriter, req *http.Request) {
	trainerId := chi.URLParam(req, "trainerId")
	appointments := r.appointmentsManager.GetAllAppointmentsByTrainerId(trainerId)

	result, _ := json.Marshal(appointments)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// makeAppointment godoc
// @Summary Make appointment
// @Description Creates an appointment for {trainerId}, based on startsAt -> endsAt
// @ID make-appointment
// @Accept  json
// @Produce  json
// @Param appointmentDetails body models.Appointment true "Appointment Details"
// @Success 204
// @Failure 404 {string} string
// @Router /appointment [post]
func (r *AppointmentRouter) makeAppointment(w http.ResponseWriter, req *http.Request) {
	appointmentDetails := &models.Appointment{}
	json.NewDecoder(req.Body).Decode(&appointmentDetails)

	err := r.appointmentsManager.MakeAppointmentWithTrainer(appointmentDetails)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
