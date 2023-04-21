package gcs

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func GCSBucketVersioningEnabled(resource commons.Resource) bool {
	bucket, ok := resource.(*GCSBucket)
	if !ok {
		return false
	}
	return bucket.Bucket.VersioningEnabled
}

func GCSBucketEncryptionEnabled(resource commons.Resource) bool {
	bucket, ok := resource.(*GCSBucket)
	if !ok {
		return false
	}
	if bucket.Bucket.Encryption != nil {
		return true
	} else {
		return false
	}
}

// TODO: avoid to do one API call per bucket
func GCSBucketPublicAccess(resource commons.Resource) bool {
	bucket, ok := resource.(*GCSBucket)
	if !ok {
		return false
	}

}
