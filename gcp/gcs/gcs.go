package gcs

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func  RunChecks(wa *sync.WaitGroup, s aws.Config, c *commons.Config, queue chan []commons.Check) {

	var checkConfig commons.CheckConfig
	checkConfig.Init(s, c)
	var checks []commons.Check
	buckets := GetListS3(s)

	S3ToEncryption := GetS3ToEncryption(s, buckets)

	go commons.CheckTest(checkConfig.Wg, c, "GCP_GCS_001", checkIfEncryptionEnabled)(checkConfig, S3ToEncryption, "AWS_GCS_001")
	// Wait for all the goroutines to finish

	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)

			checkConfig.Wg.Done()

		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
