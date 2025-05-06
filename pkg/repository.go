package pkg

import (
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ReqestsRepo struct{}

func (r *ReqestsRepo) composeCreateIndexRequest(index string, body string) *esapi.IndicesCreateRequest {
	createReq := esapi.IndicesCreateRequest{
		Index: index,
		Body:  strings.NewReader(body),
	}
	return &createReq
}
