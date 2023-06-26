package serverless

import (
	"context"
	"fmt"
	"log"

	"github.com/padok-team/yatas-gcp/internal"
	"google.golang.org/api/cloudfunctions/v1"
)

// Get the compute zones based on the list of compute regions provided in the config
func GetCloudFunctions(account internal.GCPAccount) []string {
	// Set up authentication and create the Cloud Functions client.
	ctx := context.Background()
	client, err := cloudfunctions.NewService(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// List functions in a specific project and location.
	parent := "projects/" + account.Project + "/locations/{location}"
	response, err := client.Projects.Locations.Functions.List(parent).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Failed to list functions: %v", err)
	}

	// Iterate over the functions and print their names.
	if len(response.Functions) > 0 {
		fmt.Println("Functions:")
		for _, function := range response.Functions {
			fmt.Printf("- %s\n", function.Name)
		}
	} else {
		fmt.Println("No functions found.")
	}
}
