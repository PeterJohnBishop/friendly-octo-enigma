// Package runners handles execution of scripts natively in memory
package runners

import (
	"fmt"
	"reflect"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// RunSecureGo runs a script written in go natively with variadic list of vars
func RunSecureGo(script string, args ...any) (map[string]any, error) {
	i := interp.New(interp.Options{})

	// definie the packages the script has access to
	safeSymbols := map[string]map[string]reflect.Value{
		"fmt":     stdlib.Symbols["fmt/fmt"],
		"strings": stdlib.Symbols["strings/strings"],
		"math":    stdlib.Symbols["math/math"],
		"time":    stdlib.Symbols["time/time"],
	}

	i.Use(safeSymbols)

	// evaluate the script to load the packages and function
	if _, err := i.Eval(script); err != nil {
		return nil, fmt.Errorf("go parsing error: %w", err)
	}

	// extract the function itself as a reflect.Value
	fnVal, err := i.Eval("main.Process")
	if err != nil {
		return nil, fmt.Errorf("failed to find main.Process: %w", err)
	}

	if fnVal.Kind() != reflect.Func {
		return nil, fmt.Errorf("main.Process is not a function")
	}

	// convert variadic args into a slice of reflect.Value
	inArgs := make([]reflect.Value, len(args))
	for j, arg := range args {
		inArgs[j] = reflect.ValueOf(arg)
	}

	// call the function using reflection
	out := fnVal.Call(inArgs)

	// extract and validate the return value
	if len(out) == 0 {
		return nil, fmt.Errorf("function returned no values")
	}

	result, ok := out[0].Interface().(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid return type: expected map[string]any, got %T", out[0].Interface())
	}

	return result, nil
}

func GoExample() {
	goScript := `
		package main
		import "fmt"
		
		// signature must match the types passed in via RunSecureGo
		func Process(name string, age int, tags []string, metadata map[string]any) map[string]any {
			return map[string]any{
				"message": fmt.Sprintf("Hello %s, age %d", name, age),
				"tags":    tags,
				"meta":    metadata,
				"status":  200,
			}
		}
	`

	// arbitrary types for example: string, int, slice, map
	tags := []string{"admin", "user"}
	meta := map[string]any{"session": "xyz123", "active": true}

	data, err := RunSecureGo(goScript, "Peter", 30, tags, meta)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("GO Message: %s\n", data["message"])
	fmt.Printf("Tags: %v\n", data["tags"])
	fmt.Printf("Meta: %v\n", data["meta"])
}
