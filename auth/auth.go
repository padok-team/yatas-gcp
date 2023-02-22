import (
        "context"
        "fmt"
        "io"

        "cloud.google.com/go/storage"
        "google.golang.org/api/iterator"
)


func authenticateImplicitWithAdc(w io.Writer, projectId string) error {
         projectId := "padok-lab"

        ctx := context.Background()

        // NOTE: Replace the client created below with the client required for your application.
        // Note that the credentials are not specified when constructing the client.
        // The client library finds your credentials using ADC.
        client, err := storage.NewClient(ctx)
        if err != nil {
                return fmt.Errorf("NewClient: %v", err)
        }
        defer client.Close()

        it := client.Buckets(ctx, projectId)
        for {
                bucketAttrs, err := it.Next()
                if err == iterator.Done {
                        break
                }
                if err != nil {
                        return err
                }
                fmt.Fprintf(w, "Bucket: %v\n", bucketAttrs.Name)
        }

        fmt.Fprintf(w, "Listed all storage buckets.\n")

        return nil
}
