// Package transformers handles execution of scripts natively in memory
package transformers

import (
	"fmt"
	"reflect"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func RunSecureGo(script string, name string) (map[string]interface{}, error) {
	i := interp.New(interp.Options{})

	safeSymbols := map[string]map[string]reflect.Value{
		"fmt": stdlib.Symbols["fmt/fmt"],
	}
	i.Use(safeSymbols)

	if _, err := i.Eval(script); err != nil {
		return nil, fmt.Errorf("go parsing error: %w", err)
	}

	callStr := fmt.Sprintf(`main.Process("%s")`, name)
	val, err := i.Eval(callStr)
	if err != nil {
		return nil, fmt.Errorf("go execution error: %w", err)
	}

	result, ok := val.Interface().(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid return type: expected map[string]interface{}, got %T", val.Interface())
	}

	return result, nil
}

func GoExample() {
	goScript := `
		package main
		import "fmt"
		
		func Process(name string) map[string]interface{} {
			return map[string]interface{}{
				"message": fmt.Sprintf("Hello World, I'm %s", name),
				"status":  200,
				"active":  true,
			}
		}
	`

	data, err := RunSecureGo(goScript, "Peter")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("GO Message: %s (Status: %v)\n", data["message"], data["status"])
}
