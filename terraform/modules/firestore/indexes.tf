resource "google_firestore_index" "comments_entity_tenant" {
  project    = var.project_id
  collection = "comments"
  database   = "(default)"

  fields {
    field_path = "entityId"
    order      = "ASCENDING"
  }
  fields {
    field_path = "tenantId"
    order      = "ASCENDING"
  }
  fields {
    field_path = "createdAt"
    order      = "DESCENDING"
  }
}

resource "google_firestore_index" "entity_likes_entity_tenant" {
  project    = var.project_id
  collection = "entityLikes"
  database   = "(default)"

  fields {
    field_path = "entityId"
    order      = "ASCENDING"
  }
  fields {
    field_path = "tenantId"
    order      = "ASCENDING"
  }
  fields {
    field_path = "__name__"
    order      = "ASCENDING"
  }
}