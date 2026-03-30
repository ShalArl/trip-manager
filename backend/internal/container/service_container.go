package container

import (
	"log"
	"time"

	"github.com/ShalArl/trip-manager/internal/auth"
	"github.com/ShalArl/trip-manager/internal/config"
	"github.com/ShalArl/trip-manager/internal/repository"
	"github.com/ShalArl/trip-manager/internal/service"
	"github.com/jmoiron/sqlx"
)

type ServiceConfig struct {
	DB     *sqlx.DB
	Logger *log.Logger
	Config *config.Config
}

type ServiceContainer struct {
	Trip     service.TripService
	Location service.LocationService
	User     service.UserService
	Activity service.ActivityService
	Auth     service.AuthService
}

func NewServiceContainer(cfg *ServiceConfig) *ServiceContainer {
	var tripRepo repository.TripRepository
	var locationRepo repository.LocationRepository
	var userRepo repository.UserRepository
	var activityRepo repository.ActivityRepository

	// Initialize repositories with the database connection
	tripRepo = repository.NewTripRepository(cfg.DB)
	locationRepo = repository.NewLocationRepository(cfg.DB)
	userRepo = repository.NewUserRepository(cfg.DB)
	activityRepo = repository.NewActivityRepository(cfg.DB)

	// Initialize services
	tripService := service.NewTripService(tripRepo, locationRepo, activityRepo)
	locationService := service.NewLocationService(locationRepo)
	userService := service.NewUserService(userRepo)
	activityService := service.NewActivityService(activityRepo)

	// Initialize auth manager (7 day token expiration)
	authManager := auth.NewAuthManager(cfg.Config.JWTSecret, 7*24*time.Hour)

	// Initialize auth service
	authService := service.NewAuthService(userRepo, authManager, userService)

	return &ServiceContainer{
		Trip:     tripService,
		Location: locationService,
		User:     userService,
		Activity: activityService,
		Auth:     authService,
	}
}
