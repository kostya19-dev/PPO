package main

import (
	"os"

	"vitaline-clinic/internal/delivery/cli"
	"vitaline-clinic/internal/storage/memory"
	"vitaline-clinic/internal/usecase"
)

func main() {
	patientRepo := memory.NewPatientRepository()
	appointmentRepo := memory.NewAppointmentRepository()
	visitRepo := memory.NewVisitRepository()
	scheduleRepo := memory.NewScheduleRepository()

	app := cli.NewApp(
		usecase.NewPatientService(patientRepo),
		usecase.NewAppointmentService(appointmentRepo),
		usecase.NewVisitService(visitRepo, appointmentRepo),
		usecase.NewScheduleService(scheduleRepo),
		os.Stdin,
		os.Stdout,
	)

	app.Run()
}
