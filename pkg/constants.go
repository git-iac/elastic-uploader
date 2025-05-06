package pkg

const SectionIndex string = "sections_index"

const IndexBody string = `{
		"settings": {
			"analysis": {
				"char_filter": {
					"stop_word_filter": {
						"type": "pattern_replace",
						"pattern": "\\|",
						"replacement": "###SPLIT###"
					}
				},
				"analyzer": {
					"cv_analyzer": {
						"type": "custom",
						"tokenizer": "standard",
						"char_filter": ["stop_word_filter"],
						"filter": ["lowercase", "stop"]
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"public": { "type": "boolean" },
				"sections": {
					"type": "nested",
					"properties": {
						"sectionName": { "type": "keyword" },
						"content": {
							"type": "text",
							"analyzer": "cv_analyzer",
							"index_phrases": true
						}
					}
				},
				"metadata": {
					"properties": {
						"dateReceived": { "type": "date" },
						"updatedReceived": { "type": "date" }
					}
				}
			}
		}
	}`
