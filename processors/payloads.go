// Package processors identify data in the webhook payload that can be passed to the transformers.
package processors

import (
	"encoding/json"
	"fmt"
	"maps"

	"github.com/peterjohnbishop/friendly-octo-enigma/models"
)

func MapAndMergeHeaders(headers chan map[string][]string, mappedHeaders chan map[string][]string) {
	mergedHeaders := make(map[string][]string)

	for h := range headers {
		maps.Copy(mergedHeaders, h)

		outMap := make(map[string][]string)
		maps.Copy(outMap, mergedHeaders)

		mappedHeaders <- outMap
	}
}

func ProcessValue(val any) models.PropertyDetail {
	switch v := val.(type) {
	case string:
		return models.PropertyDetail{Value: v, InferredType: "string"}
	case float64:
		return models.PropertyDetail{Value: v, InferredType: "number"}
	case bool:
		return models.PropertyDetail{Value: v, InferredType: "boolean"}
	case []any:
		var parsedArray []models.PropertyDetail
		for _, item := range v {
			parsedArray = append(parsedArray, ProcessValue(item))
		}
		return models.PropertyDetail{Value: parsedArray, InferredType: "array"}
	case map[string]any:
		parsedObject := make(map[string]models.PropertyDetail)
		for k, item := range v {
			parsedObject[k] = ProcessValue(item)
		}
		return models.PropertyDetail{Value: parsedObject, InferredType: "object"}
	case nil:
		return models.PropertyDetail{Value: nil, InferredType: "null"}
	default:
		return models.PropertyDetail{Value: v, InferredType: "unknown"}
	}
}

func MapBody(payloadChan chan []byte, typedPayloadChan chan map[string]models.PropertyDetail) {
	for rawBytes := range payloadChan {
		var originalData map[string]any

		if err := json.Unmarshal(rawBytes, &originalData); err != nil {
			fmt.Printf("error processing request body: %s\n", err)
			continue
		}

		typedPayload := make(map[string]models.PropertyDetail)

		for key, val := range originalData {
			typedPayload[key] = ProcessValue(val)
		}

		typedPayloadChan <- typedPayload

	}
}
