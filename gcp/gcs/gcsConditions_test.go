package gcs

import (
	"testing"

	"cloud.google.com/go/storage"
)

type FakeResource struct{}

func (f FakeResource) GetID() string {
	return "fake-resource"
}

func TestGCSBucketVersioningEnabled(t *testing.T) {
	// Create a test bucket with versioning enabled
	bucket := &GCSBucket{
		Bucket: storage.BucketAttrs{Name: "test-bucket", VersioningEnabled: true},
	}
	if !GCSBucketVersioningEnabled(bucket) {
		t.Error("Expected GCSBucketVersioningEnabled to return true for bucket with versioning enabled")
	}

	// Create a test bucket with versioning disabled
	bucket = &GCSBucket{
		Bucket: storage.BucketAttrs{Name: "test-bucket", VersioningEnabled: false},
	}
	if GCSBucketVersioningEnabled(bucket) {
		t.Error("Expected GCSBucketVersioningEnabled to return false for bucket with versioning disabled")
	}

	// Create a test resource that is not a GCSBucket
	resource := FakeResource{}
	if GCSBucketVersioningEnabled(resource) {
		t.Error("Expected GCSBucketVersioningEnabled to return false for non-GCSBucket resource")
	}
}
