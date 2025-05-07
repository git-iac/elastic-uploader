package main

import (
	"context"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/git-iac/elastic-uploader/cli"
	"github.com/git-iac/elastic-uploader/pkg"
)

func main() {
	pkg.WithMetrics(func() {
		opts := cli.Parse()
		ctx := context.Background()
		cfg := elasticsearch.Config{
			Addresses: []string{os.Getenv("ELASTIC_URL")},
			APIKey:    os.Getenv("ELASTIC_API_KEY"),
		}
		client, err := elasticsearch.NewClient(cfg)
		if err != nil {
			log.Fatalf("Error creating the Elasticsearch client: %s", err)
		}

		cltService := pkg.NewElasticService(client)

		switch opts.IndexAction {
		case cli.PopulateIndex:
			log.Printf("populating index [%s]...", pkg.SectionIndex)
			if err := cltService.UploadBulk(ctx, opts.ReadChunkSize); err != nil {
				//TODO: errors with special chars look ugly
				log.Fatalf("error populating index: %s", err)
			}
		case cli.CreateIndex:
			log.Printf("creating index [%s]...", pkg.SectionIndex)
			if err := cltService.CreateIndexForSections(ctx, pkg.IndexBody, pkg.SectionIndex); err != nil {
				log.Fatalf("error creating index: %s", err)
			}
		default:
			log.Fatal("no action flag was passed")
		}

	})
}
