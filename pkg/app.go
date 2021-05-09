package app

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	database "github.com/parkedwards/go-trainer-api/pkg/database"
	AppointmentManager "github.com/parkedwards/go-trainer-api/pkg/managers/appointment"
	AvailabilityManager "github.com/parkedwards/go-trainer-api/pkg/managers/availability"
	AppointmentRouter "github.com/parkedwards/go-trainer-api/pkg/routers/appointment"
	AvailabilityRouter "github.com/parkedwards/go-trainer-api/pkg/routers/availability"
)

func Init() *chi.Mux {
	router := chi.NewRouter()

	// init database client
	db := database.New()

	// init managers
	availabilityManager := AvailabilityManager.New(db)
	appointmentManager := AppointmentManager.New(db)

	// init route handlers
	availabilityRouter := AvailabilityRouter.New(availabilityManager)
	appointmentRouter := AppointmentRouter.New(appointmentManager)

	// register routes
	availabilityRouter.RegisterRoutes(router)
	appointmentRouter.RegisterRoutes(router)

	return router
}

func Boot(router *chi.Mux) {
	port := ":8000"
	fmt.Printf("bootin it up at %v", port)
	http.ListenAndServe(port, router)
}
