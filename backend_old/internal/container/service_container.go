package container

import (
	"log"

	"cloud.google.com/go/firestore"
	"github.com/ShalArl/trip-manager/internal/config"
	"github.com/ShalArl/trip-manager/internal/infrastructure"
	"github.com/ShalArl/trip-manager/internal/repository"
	"github.com/ShalArl/trip-manager/internal/service"
	"github.com/ShalArl/trip-manager/internal/storage"
	"github.com/jmoiron/sqlx"
)

type ServiceConfig struct {
	SQLDb           *sqlx.DB
	FirestoreClient *firestore.Client
	Logger          *log.Logger
	Config          *config.Config
	Storage         storage.Storage
}

type ServiceContainer struct {
	Trip          service.TripService
	Location      service.LocationService
	User          service.UserService
	Activity      service.ActivityService
	Media         infrastructure.MediaService
	Transport     service.TransportService
	Social        service.SocialService
	Accommodation service.AccommodationService
}

func NewServiceContainer(cfg *ServiceConfig) (*ServiceContainer, error) {
	// Initialize media service (needed by handlers for presigned URLs)
	mediaService := infrastructure.NewMediaService(cfg.Storage, cfg.Config.Storage.SignedURLTTL)
	socialRepo := repository.NewSocialRepository(cfg.FirestoreClient)

	// Initialize repositories with the database connection
	tripRepo := repository.NewTripRepository(cfg.SQLDb)
	locationRepo := repository.NewLocationRepository(cfg.SQLDb)
	userRepo := repository.NewUserRepository(cfg.SQLDb)
	activityRepo := repository.NewActivityRepository(cfg.SQLDb)
	transportRepo := repository.NewTransportRepository(cfg.SQLDb)
	accommodationRepo := repository.NewAccommodationRepository(cfg.SQLDb)

	// Initialize services
	tripService := service.NewTripService(tripRepo, locationRepo, activityRepo)
	locationService := service.NewLocationService(locationRepo, mediaService)

	// Initialize user service
	userService := service.NewUserService(userRepo, mediaService)
	activityService := service.NewActivityService(activityRepo)

	// Initialize transport and accomodation service
	transportService := service.NewTransportService(transportRepo)
	accommodationService := service.NewAccommodationService(accommodationRepo)

	// Initialize social service
	socialService := service.NewSocialService(socialRepo, userRepo, mediaService)

	return &ServiceContainer{
		Trip:          tripService,
		Location:      locationService,
		User:          userService,
		Activity:      activityService,
		Media:         mediaService,
		Transport:     transportService,
		Accommodation: accommodationService,
		Social:        socialService,
	}, nil
}
