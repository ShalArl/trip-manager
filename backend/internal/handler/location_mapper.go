package handler

import (
	"context"

	"github.com/ShalArl/trip-manager/internal/domain"
	"github.com/ShalArl/trip-manager/internal/generated"
	"github.com/ShalArl/trip-manager/internal/infrastructure"
	"github.com/ShalArl/trip-manager/pkg/ptr"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

func mapLocationToLocationResponse(ctx context.Context, media infrastructure.MediaService, l *domain.Location) *generated.LocationResponse {
	id, _ := uuid.Parse(l.ID)

	images := make([]generated.LocationImageResponse, 0, len(l.Images))
	for _, img := range l.Images {
		url, err := media.GetDownloadURL(ctx, img.ImageKey)
		if err == nil {
			id, _ := uuid.Parse(img.ID)
			locId, _ := uuid.Parse(img.LocationID)
			images = append(images, generated.LocationImageResponse{
				Id:         id,
				LocationId: locId,
				ImageUrl:   url,
				Sequence:   &img.Sequence,
				CreatedAt:  &img.CreatedAt,
			})
		}
	}

	return &generated.LocationResponse{
		Id:               ptr.ToPtr(id),
		Name:             l.Name,
		City:             l.City,
		Country:          l.Country,
		ShortDescription: l.ShortDescription,
		DateFrom:         openapitypes.Date{Time: l.DateFrom},
		DateTo:           openapitypes.Date{Time: l.DateTo},
		Latitude:         ptr.ToPtr(float32(l.Coordinates.Lat)),
		Longitude:        ptr.ToPtr(float32(l.Coordinates.Lon)),
		Notes:            ptr.ToPtr(l.Notes),
		Sequence:         ptr.ToPtr(l.Sequence),
		CreatedAt:        &l.CreatedAt,
		UpdatedAt:        &l.UpdatedAt,
		Images:           &images, CreatedBy: &generated.UserSummary{
			Id:    uuid.MustParse(l.CreatedBy.ID),
			Name:  l.CreatedBy.Name,
			Email: openapitypes.Email(l.CreatedBy.Email),
		},
	}
}

func mapLocationsToLocationListResponse(ctx context.Context, media infrastructure.MediaService, locations []*domain.Location, limit int, offset int, total int) *generated.LocationListResponse {
	loc := make([]generated.LocationResponse, len(locations))
	for i, location := range locations {
		loc[i] = *mapLocationToLocationResponse(ctx, media, location)
	}
	return &generated.LocationListResponse{
		Data:   loc,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
}
