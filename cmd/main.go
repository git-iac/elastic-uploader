package main

import (
	"context"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/git-iac/elastic-uploader/pkg"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	cfg := elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTIC_URL")},
		APIKey:    os.Getenv("ELASTIC_API_KEY"),
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}

	cl := pkg.NewElasticClient(client)

	if err := cl.GetInfo(ctx); err != nil {
		log.Fatal(err)
	}
	err = cl.CreateIndexForSections(ctx, pkg.SectionIndex, pkg.IndexBody)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}

}
