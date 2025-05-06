package pkg

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticClient struct {
	E *elasticsearch.Client
	ReqestsRepo
}

func NewElasticClient(e *elasticsearch.Client) *ElasticClient {
	return &ElasticClient{
		E: e,
	}
}

func (c *ElasticClient) GetInfo(ctx context.Context) error {
	res, err := c.E.Info(c.E.Info.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error getting response from Elasticsearch: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	log.Println("Elasticsearch client connected.")
	return nil
}

func (c *ElasticClient) CreateIndexForSections(ctx context.Context, index string, body string) error {
	if err := c.checkIfIndexExists(ctx, index); err != nil {
		return err
	}
	createReq := c.composeCreateIndexRequest(index, body)
	createRes, err := createReq.Do(ctx, c.E)
	if err != nil {
		return fmt.Errorf("error sending index create request: %w", err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		return fmt.Errorf("error creating index '%s': %s", index, createRes.String())
	}

	log.Printf("Index '%s' created successfully", index)

	return nil
}

func (c *ElasticClient) checkIfIndexExists(ctx context.Context, index string) error {
	res, err := c.E.Indices.Exists([]string{index}, c.E.Indices.Exists.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error checking if index exists: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	}
	if res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected status code %d while checking index existence", res.StatusCode)
	}
	return nil
}
