resource "google_pubsub_topic" "trip_events" {
  name    = "trip-events"
  project = var.project_id
}

resource "google_pubsub_subscription" "trip_events_sub" {
  name  = "trip-events-sub"
  topic = google_pubsub_topic.trip_events.name

  ack_deadline_seconds  = 20
  retain_acked_messages = false
}