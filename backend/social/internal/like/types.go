package like

type EntityLikeResponse struct {
	LikeCount int  `json:"likeCount"`
	HasLiked  bool `json:"hasLiked"`
}

type TargetType int

const (
	TargetTypeUnknown TargetType = iota
	TargetTypeComment
	TargetTypeTrip
)
