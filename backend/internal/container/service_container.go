package container

import (
	"log"

	"github.com/ShalArl/trip-manager/internal/config"
	"github.com/ShalArl/trip-manager/internal/infrastructure"
	"github.com/ShalArl/trip-manager/internal/repository"
	"github.com/ShalArl/trip-manager/internal/service"
	"github.com/ShalArl/trip-manager/internal/storage"
	"github.com/jmoiron/sqlx"
)

type ServiceConfig struct {
	DB      *sqlx.DB
	Logger  *log.Logger
	Config  *config.Config
	Storage storage.Storage
}

type ServiceContainer struct {
	Trip     service.TripService
	Location service.LocationService
	User     service.UserService
	Activity service.ActivityService
	Media    infrastructure.MediaService
}

func NewServiceContainer(cfg *ServiceConfig) (*ServiceContainer, error) {
	// Initialize media service (needed by handlers for presigned URLs)
	mediaService := infrastructure.NewMediaService(cfg.Storage, cfg.Config.Storage.SignedURLTTL)

	// Initialize repositories with the database connection
	tripRepo := repository.NewTripRepository(cfg.DB)
	locationRepo := repository.NewLocationRepository(cfg.DB)
	userRepo := repository.NewUserRepository(cfg.DB)
	activityRepo := repository.NewActivityRepository(cfg.DB)

	// Initialize services
	tripService := service.NewTripService(tripRepo, locationRepo, activityRepo)
	locationService := service.NewLocationService(locationRepo)

	// Initialize user service
	userService := service.NewUserService(userRepo, mediaService)
	activityService := service.NewActivityService(activityRepo)

	return &ServiceContainer{
		Trip:     tripService,
		Location: locationService,
		User:     userService,
		Activity: activityService,
		Media:    mediaService,
	}, nil
}
