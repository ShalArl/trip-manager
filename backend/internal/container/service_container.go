package container

import (
	"log"

	"github.com/ShalArl/trip-manager/internal/auth"
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
	Auth     service.AuthService
	Media    infrastructure.MediaService
}

func NewServiceContainer(cfg *ServiceConfig) (*ServiceContainer, error) {
	// Initialize repositories with the database connection
	tripRepo := repository.NewTripRepository(cfg.DB)
	locationRepo := repository.NewLocationRepository(cfg.DB)
	userRepo := repository.NewUserRepository(cfg.DB)
	activityRepo := repository.NewActivityRepository(cfg.DB)

	// Initialize services
	tripService := service.NewTripService(tripRepo, locationRepo, activityRepo)
	locationService := service.NewLocationService(locationRepo)

	// Initialize user service
	userService := service.NewUserService(userRepo)
	activityService := service.NewActivityService(activityRepo)

	// Initialize media service (needed by handlers for presigned URLs)
	mediaService := infrastructure.NewMediaService(cfg.Storage)

	// Initialize auth manager (7 day token expiration)
	authManager := auth.NewAuthManager(cfg.Config.JWTSecret, cfg.Config.TokenExpiration)

	// Initialize auth service
	authService := service.NewAuthService(authManager, userService)

	return &ServiceContainer{
		Trip:     tripService,
		Location: locationService,
		User:     userService,
		Activity: activityService,
		Auth:     authService,
		Media:    mediaService,
	}, nil
}
