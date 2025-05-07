package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func NewCreateIndexRequest(index string, body string) *esapi.IndicesCreateRequest {
	createReq := esapi.IndicesCreateRequest{
		Index: index,
		Body:  strings.NewReader(body),
	}
	return &createReq
}

func NewUploadBulkRequest(fs *FileSections) *esapi.BulkRequest {
	reader, err := createBulkRequestBody(convertToElasticDocument(fs), SectionIndex)
	if err != nil {
		log.Fatal(err)
	}
	bulkReq := esapi.BulkRequest{
		Body:    reader,
		Refresh: "wait_for",
	}
	return &bulkReq
}

func convertToElasticDocument(fSec *FileSections) *[]ElasticsearchDocument {
	documents := make([]ElasticsearchDocument, 0, len((*fSec)))
	var ed ElasticsearchDocument
	for k, v := range *fSec {
		ed.ID = k
		ed.Metadata.DateReceived = time.Now()
		ed.Metadata.UpdatedReceived = time.Now()
		ed.Public = true
		ed.Sections = v
		documents = append(documents, ed)

		ed = ElasticsearchDocument{}
	}
	return &documents
}

func createBulkRequestBody(documents *[]ElasticsearchDocument, indexName string) (io.Reader, error) {
	if documents == nil {
		return nil, fmt.Errorf("document slice is nil")
	}

	var buf strings.Builder

	for _, doc := range *documents {
		meta := map[string]map[string]string{
			"index": {
				"_index": indexName,
				"_id":    doc.ID,
			},
		}
		metaJSON, err := json.Marshal(meta)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata for doc ID %s: %w", doc.ID, err)
		}

		// Append metadata line and a newline
		buf.Write(metaJSON)
		buf.WriteByte('\n')

		// --- Construct Document Source Line ---
		docJSON, err := json.Marshal(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal document ID %s: %w", doc.ID, err)
		}

		// Append document source line and a newline
		buf.Write(docJSON)
		buf.WriteByte('\n')
	}

	return strings.NewReader(buf.String()), nil
}
