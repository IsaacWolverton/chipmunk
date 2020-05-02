resource "google_storage_bucket" "chipmunk-storage" {
    name     = "chipmunk-storage"
    project  = var.project
    location = "US"
}

resource "google_storage_bucket_object" "chipmunk-application" {
    bucket = google_storage_bucket.chipmunk-storage.name

    name   = "${var.application_image}/application.tar"
    source = var.application_path
}