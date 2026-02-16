package swagger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"time"
)

func fetchDocsJSON(microserviceAddress string) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", microserviceAddress+"/swagger/doc.json", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Microservice returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func mergeDocs(docs ...map[string]any) map[string]any {
	merged := make(map[string]any)

	merged["swagger"] = "2.0"
	merged["info"] = map[string]any{
		"title":   "Library backend API",
		"version": "1.0.0",
	}
	merged["securityDefinitions"] = map[string]any{
		"BearerAuth": map[string]any{
			"type": "apiKey",
			"name": "Authorization",
			"in":   "header",
		},
	}

	mergedPaths := make(map[string]any)
	mergedDefinitions := make(map[string]any)
	var mergedTags []any

	for _, doc := range docs {
		if doc == nil {
			continue
		}

		if paths, ok := doc["paths"].(map[string]any); ok {
			maps.Copy(mergedPaths, paths)
		}
		if defs, ok := doc["definitions"].(map[string]any); ok {
			maps.Copy(mergedDefinitions, defs)
		}
		if tags, ok := doc["tags"].([]any); ok {
			mergedTags = append(mergedTags, tags...)
		}
	}

	merged["paths"] = mergedPaths
	merged["definitions"] = mergedDefinitions
	merged["tags"] = mergedTags

	return merged
}
