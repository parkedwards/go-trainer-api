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

	_ "github.com/parkedwards/go-trainer-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	port = ":8000"
)

// @title Trainer API
// @version 1.0
// @description Find and schedule a time with a trainer.
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

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost%v/swagger/doc.json", port)), //The url pointing to API definition"
	))

	return router
}

func Boot(router *chi.Mux) {
	fmt.Printf("bootin it up at %v", port)
	err := http.ListenAndServe(port, router)

	fmt.Println(err)
}
