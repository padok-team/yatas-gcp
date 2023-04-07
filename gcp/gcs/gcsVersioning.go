package gcs

import (
	"cloud.google.com/go/storage"
	"github.com/padok-team/yatas/plugins/commons"
)

func CheckIfVersioningEnabled(queue chan commons.Check, buckets []storage.BucketAttrs, testName string) {
	var check commons.Check
	check.InitCheck("GCS buckets are versioned", "Check if GCS buckets are using object versioning", testName, []string{"Security", "Good Practice"})
	for _, bucket := range buckets {
		if !bucket.VersioningEnabled {
			Message := "GCS bucket " + bucket.Name + " is not using object versioning"
			result := commons.Result{Status: "FAIL", Message: Message, ResourceID: bucket.Name}
			check.AddResult(result)
		} else {
			Message := "GCS bucket " + bucket.Name + " is using object versioning"
			result := commons.Result{Status: "OK", Message: Message, ResourceID: bucket.Name}
			check.AddResult(result)
		}
	}
	queue <- check
}
