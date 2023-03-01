package gcs

import (
	"fmt"
	"sync"

	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, client internal.Client_Account, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	// checkConfig.Init(s, c)
	var checks []commons.Check
	buckets := GetListGCS(client)
	fmt.Printf("%v", buckets)

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
