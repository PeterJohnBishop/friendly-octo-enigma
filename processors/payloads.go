// Package processors handles data processing for endpoints
package processors

import (
	"encoding/json"
	"maps"
)

type PropertyDetail struct {
	Value        any    `json:"value"`
	InferredType string `json:"inferred_type"`
}

func MapAndMergeHeaders(headers chan map[string][]string) map[string][]string {
	mergedHeaders := make(map[string][]string)
	for h := range headers {
		maps.Copy(mergedHeaders, h)
	}

	return mergedHeaders
}

func inferType(val any) string {
	switch val.(type) {
	case string:
		return "string"
	case float64:
		return "number" // Go unmarshals all JSON numbers into float64
	case bool:
		return "boolean"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	case nil:
		return "null"
	default:
		return "unknown"
	}
}

func MapBody(payloadChan chan []byte) ([]map[string]PropertyDetail, error) {
	var results []map[string]PropertyDetail

	for rawBytes := range payloadChan {

		var originalData map[string]any

		if err := json.Unmarshal(rawBytes, &originalData); err != nil {
			// skip invalid JSON payloads
			continue
		}

		typedPayload := make(map[string]PropertyDetail)

		for key, val := range originalData {
			typedPayload[key] = PropertyDetail{
				Value:        val,
				InferredType: inferType(val),
			}
		}

		results = append(results, typedPayload)
	}

	return results, nil
}
