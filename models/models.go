// Package models provides shared data models
package models

type PropertyDetail struct {
	Value        any    `json:"value"`
	InferredType string `json:"inferred_type"`
}
