package handler

import (
	"net/http"

	"github.com/ShalArl/trip-manager/internal/app"
)

func LikeActivityHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func UnlikeActivityHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func ListCommentsHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func CreateCommentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func UpdateCommentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func DeleteCommentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func LikeCommentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func UnlikeCommentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func ListRepliesHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func CreateReplyHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func UpdateReplyHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func DeleteReplyHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func LikeReplyHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func UnlikeReplyHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
