package domain

import "time"

type ResourceMeta struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy UserSummary
}

type UserSummary struct {
	ID        string
	Name      string
	Email     string
	AvatarKey *string
}

type ListResult[T any] struct {
	Data       []T
	TotalCount int
}
