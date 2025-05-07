package pkg

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticsearchDocument struct {
	ID       string           `json:"id"`
	Public   bool             `json:"public"`
	Sections []Section        `json:"sections"`
	Metadata DocumentMetadata `json:"metadata"`
}
type Section struct {
	SectionName string `json:"sectionName"`
	Content     string `json:"content"`
}
type DocumentMetadata struct {
	DateReceived    time.Time `json:"dateReceived"`
	UpdatedReceived time.Time `json:"updatedReceived"`
}

type ElasticService struct {
	Clt *elasticsearch.Client
}

func NewElasticService(e *elasticsearch.Client) *ElasticService {
	return &ElasticService{
		Clt: e,
	}
}

func (c *ElasticService) doRequest(ctx context.Context, req esapi.Request) error {
	createRes, err := req.Do(ctx, c.Clt)
	if err != nil {
		return fmt.Errorf("error sending index create request: %w", err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		return fmt.Errorf("%s", createRes.String())
	}
	return nil
}

func (c *ElasticService) GetInfo(ctx context.Context) error {
	res, err := c.Clt.Info(c.Clt.Info.WithContext(ctx))
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

func (c *ElasticService) UploadBulk(ctx context.Context, chunkSize int) error {
	f, err := GetFileSections(chunkSize)
	if err != nil {
		return err
	}
	bulkReq := NewUploadBulkRequest(f)
	if err := c.doRequest(ctx, bulkReq); err != nil {
		return fmt.Errorf("error doing bulk create: %w", err)
	}

	return nil
}

func (c *ElasticService) CreateIndexForSections(ctx context.Context, index string, body string) error {
	if err := c.checkIfIndexExists(ctx, index); err != nil {
		return err
	}
	createReq := NewCreateIndexRequest(index, body)
	if err := c.doRequest(ctx, createReq); err != nil {
		return fmt.Errorf("error creating index: '%s': %w", index, err)
	}
	log.Printf("Index '%s' created successfully", index)

	return nil
}

func (c *ElasticService) checkIfIndexExists(ctx context.Context, index string) error {
	res, err := c.Clt.Indices.Exists([]string{index}, c.Clt.Indices.Exists.WithContext(ctx))
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
