package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"vitaline-clinic/internal/appconfig"
	"vitaline-clinic/internal/logging"
	"vitaline-clinic/internal/storage/postgres"
	"vitaline-clinic/internal/usecase"
)

type Options struct {
	ConfigPath   string
	Debug        bool
	SkipDatabase bool
}

type UseCases interface {
	Patients() usecase.IPatientService
	Appointments() usecase.IAppointmentService
	Schedule() usecase.IScheduleService
	Visits() usecase.IVisitService
}

type HealthChecker interface {
	Ping(ctx context.Context) error
}

type closer interface {
	Close() error
}

type useCases struct {
	patients     usecase.IPatientService
	appointments usecase.IAppointmentService
	schedule     usecase.IScheduleService
	visits       usecase.IVisitService
}

func (s *useCases) Patients() usecase.IPatientService {
	return s.patients
}

func (s *useCases) Appointments() usecase.IAppointmentService {
	return s.appointments
}

func (s *useCases) Schedule() usecase.IScheduleService {
	return s.schedule
}

func (s *useCases) Visits() usecase.IVisitService {
	return s.visits
}

type databaseHealthChecker struct {
	db *sql.DB
}

func (h *databaseHealthChecker) Ping(ctx context.Context) error {
	return h.db.PingContext(ctx)
}

type App struct {
	ConfigPath    string
	AppName       string
	APIAddress    string
	GUIAddress    string
	GUIAPIBaseURL string
	Logger        *slog.Logger
	UseCases      UseCases
	Health        HealthChecker

	closers []closer
}

func New(ctx context.Context, opts Options) (*App, error) {
	cfg, configPath, err := appconfig.LoadWithOptions(opts.ConfigPath, appconfig.LoadOptions{
		RequireDatabase: !opts.SkipDatabase,
	})
	if err != nil {
		return nil, err
	}

	logs, err := logging.New(cfg.Logging, opts.Debug)
	if err != nil {
		return nil, err
	}

	if opts.SkipDatabase {
		app := &App{
			ConfigPath:    configPath,
			AppName:       cfg.App.Name,
			APIAddress:    cfg.API.Address,
			GUIAddress:    cfg.GUI.Address,
			GUIAPIBaseURL: cfg.GUI.APIBaseURL,
			Logger:        logs.Logger,
			closers:       []closer{logs},
		}
		app.Logger.Info(
			"application bootstrap completed",
			"config_path", app.ConfigPath,
			"app_name", cfg.App.Name,
			"environment", cfg.App.Environment,
			"debug", opts.Debug,
			"http_address", cfg.API.Address,
			"gui_address", cfg.GUI.Address,
			"gui_api_base_url", cfg.GUI.APIBaseURL,
		)
		return app, nil
	}

	db, err := postgres.Open(ctx, cfg.Database.URL)
	if err != nil {
		_ = logs.Close()
		return nil, err
	}

	patientRepo := postgres.NewPatientRepository(db)
	appointmentRepo := postgres.NewAppointmentRepository(db)
	scheduleRepo := postgres.NewScheduleRepository(db)
	visitRepo := postgres.NewVisitRepository(db)

	minLeadTime := time.Duration(cfg.Business.AppointmentMinLeadTimeMinutes) * time.Minute

	patientSvc := usecase.NewPatientService(patientRepo)
	appointmentSvc := usecase.NewAppointmentService(
		appointmentRepo,
		usecase.WithAppointmentMinLeadTime(minLeadTime),
	)
	scheduleSvc := usecase.NewScheduleService(scheduleRepo)
	visitSvc := usecase.NewVisitService(visitRepo, appointmentRepo)

	app := &App{
		ConfigPath:    configPath,
		AppName:       cfg.App.Name,
		APIAddress:    cfg.API.Address,
		GUIAddress:    cfg.GUI.Address,
		GUIAPIBaseURL: cfg.GUI.APIBaseURL,
		Logger:        logs.Logger,
		UseCases: &useCases{
			patients:     patientSvc,
			appointments: appointmentSvc,
			schedule:     scheduleSvc,
			visits:       visitSvc,
		},
		Health:  &databaseHealthChecker{db: db},
		closers: []closer{db, logs},
	}

	app.Logger.Info(
		"application bootstrap completed",
		"config_path", app.ConfigPath,
		"app_name", cfg.App.Name,
		"environment", cfg.App.Environment,
		"debug", opts.Debug,
		"http_address", cfg.API.Address,
		"gui_address", cfg.GUI.Address,
		"appointment_min_lead_time_minutes", cfg.Business.AppointmentMinLeadTimeMinutes,
	)

	return app, nil
}

func (a *App) Close() error {
	var result error

	for _, c := range a.closers {
		if c == nil {
			continue
		}
		if err := c.Close(); err != nil && result == nil {
			result = fmt.Errorf("close resource: %w", err)
		}
	}

	return result
}
