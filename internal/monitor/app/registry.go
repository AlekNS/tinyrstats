package app

import (
	"context"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/infrastructure"
	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/alekns/tinyrstats/internal/monitor/services"
	"github.com/go-kit/kit/log"
)

// Registry factory of services and applications.
type Registry interface {
	// Events get service.
	Events() monitor.Events
	// TaskRepository get service.
	TaskRepository() monitor.TaskRepository
	// ScheduleTaskService get service.
	ScheduleTaskService() monitor.ScheduleTaskService

	// TaskApp get application.
	TaskApp() monitor.TaskApp
	// StatsApp get application.
	StatsApp() monitor.StatsApp
}

// RegistryImpl .
type RegistryImpl struct {
	Registry
	events         monitor.Events
	taskRepository monitor.TaskRepository
	taskConsumer   *services.ConcurrentTaskExecutor
	scheduler      *services.IntervalScheduler

	taskApp  *taskAppImpl
	statsApp *statsAppImpl
}

// Events .
func (app *RegistryImpl) Events() monitor.Events {
	return app.events
}

// TaskRepository .
func (app *RegistryImpl) TaskRepository() monitor.TaskRepository {
	return app.taskRepository
}

// ScheduleTaskService .
func (app *RegistryImpl) ScheduleTaskService() monitor.ScheduleTaskService {
	return app.scheduler
}

// TaskApp .
func (app *RegistryImpl) TaskApp() monitor.TaskApp {
	return app.taskApp
}

// StatsApp .
func (app *RegistryImpl) StatsApp() monitor.StatsApp {
	return app.statsApp
}

// Init .
func (app *RegistryImpl) init(
	settings *config.Settings, logger log.Logger) *RegistryImpl {

	app.events = NewSyncEventsImpl()

	app.taskRepository = infrastructure.NewTaskRepositoryInMemory(settings.Tasks.RepositoryBucketsCount)

	healthTaskConsumer := services.NewHealthServiceConsumer(
		services.NewHTTPHealthService(settings.Tasks, logger),
		services.NewResultsTaskRepositoryConsumer(app.taskRepository))

	app.taskConsumer = services.NewConcurretTaskExecutor(settings.Tasks, logger, healthTaskConsumer)

	app.scheduler = services.NewIntervalScheduler(logger, settings.Scheduler, app.taskConsumer)

	app.taskApp = newTaskApp(settings, logger, app.events, app.scheduler, app.taskRepository)

	statsService := services.NewStatsServiceImpl(settings.Stats)
	app.statsApp = newStatsApp(settings, logger, statsService, app.events)

	return app
}

// Start is method to start necessary services.
func (app *RegistryImpl) Start(ctx context.Context) error {
	err := app.taskConsumer.Start(ctx)
	if err != nil {
		return err
	}

	err = app.scheduler.Start(ctx)
	if err != nil {
		app.taskConsumer.Stop(ctx)
		return nil
	}

	return nil
}

// Stop necessary services.
func (app *RegistryImpl) Stop(ctx context.Context) error {
	app.events.TaskQueriedByResource().OffAll()
	app.events.TaskQueriedByMinResponse().OffAll()
	app.events.TaskQueriedByMaxResponse().OffAll()
	app.scheduler.Stop(ctx)
	app.taskConsumer.Stop(ctx)

	return nil
}

// NewRegistryImpl creates default registry.
func NewRegistryImpl(settings *config.Settings, logger log.Logger) *RegistryImpl {
	instance := &RegistryImpl{}
	return instance.init(settings, logger)
}
