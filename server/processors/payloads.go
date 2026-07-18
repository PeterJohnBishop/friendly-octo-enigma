// Package processors handles data processing for endpoints
package processors

import (
	"encoding/json"
	"fmt"
	"maps"
)

type PropertyDetail struct {
	Value        any    `json:"value"`
	InferredType string `json:"inferred_type"`
}

func MapAndMergeHeaders(headers chan map[string][]string) {
	mergedHeaders := make(map[string][]string)

	for h := range headers {
		maps.Copy(mergedHeaders, h)

		prettyHeaders, _ := json.MarshalIndent(mergedHeaders, "", "  ")
		fmt.Printf("--- Payload Headers -->\n%s\n\n", string(prettyHeaders))
	}
}

func ProcessValue(val any) PropertyDetail {
	switch v := val.(type) {
	case string:
		return PropertyDetail{Value: v, InferredType: "string"}
	case float64:
		return PropertyDetail{Value: v, InferredType: "number"}
	case bool:
		return PropertyDetail{Value: v, InferredType: "boolean"}
	case []any:
		var parsedArray []PropertyDetail
		for _, item := range v {
			parsedArray = append(parsedArray, ProcessValue(item))
		}
		return PropertyDetail{Value: parsedArray, InferredType: "array"}
	case map[string]any:
		parsedObject := make(map[string]PropertyDetail)
		for k, item := range v {
			parsedObject[k] = ProcessValue(item)
		}
		return PropertyDetail{Value: parsedObject, InferredType: "object"}
	case nil:
		return PropertyDetail{Value: nil, InferredType: "null"}
	default:
		return PropertyDetail{Value: v, InferredType: "unknown"}
	}
}

func MapBody(payloadChan chan []byte) {
	for rawBytes := range payloadChan {
		var originalData map[string]any

		if err := json.Unmarshal(rawBytes, &originalData); err != nil {
			fmt.Printf("error processing request body: %s\n", err)
			continue
		}

		typedPayload := make(map[string]PropertyDetail)

		for key, val := range originalData {
			typedPayload[key] = ProcessValue(val)
		}

		prettyBody, _ := json.MarshalIndent(typedPayload, "", "  ")
		fmt.Printf("--- Processed Body ---\n%s\n\n", string(prettyBody))
	}
}
