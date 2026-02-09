package swagger

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"

	"go.uber.org/zap"
)

func fetchDocsJSON(serviceURL string) map[string]any {
	resp, err := http.Get(fmt.Sprintf("%s/swagger/doc.json", serviceURL))
	if err != nil {
		zap.S().Errorf("Error fetching swagger from %s: %v\n", serviceURL, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		zap.S().Errorf("Error fetching swagger from %s: %v\n, wrong status code: %d", serviceURL, resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.S().Errorf("Error reading swagger body from %s: %v\n", serviceURL, err)
		return nil
	}

	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		zap.S().Errorf("Error unmarshalling swagger from %s: %v\n", serviceURL, err)
		return nil
	}
	return m
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
