package runners

import (
	"fmt"

	"github.com/dop251/goja"
)

// RunJavaScript runs a script natively in JavaScript and accepts a varadic liist of vars
func RunJavaScript(script string, args ...any) (map[string]any, error) {
	vm := goja.New()

	// run the script to define the functions in the JS environment
	if _, err := vm.RunString(script); err != nil {
		return nil, fmt.Errorf("js parsing error: %w", err)
	}

	// retrieve the target function
	fnVal := vm.Get("process")
	if fnVal == nil {
		return nil, fmt.Errorf("function 'process' not found in script")
	}

	// assert that the retrieved value is actually a callable function
	processFn, ok := goja.AssertFunction(fnVal)
	if !ok {
		return nil, fmt.Errorf("'process' is not a function")
	}

	// convert Go arguments into goja.Value arguments
	jsArgs := make([]goja.Value, len(args))
	for i, arg := range args {
		jsArgs[i] = vm.ToValue(arg)
	}

	// execute the JavaScript function
	// The first parameter is 'this' (pass undefined), followed by arguments
	resultVal, err := processFn(goja.Undefined(), jsArgs...)
	if err != nil {
		return nil, fmt.Errorf("js execution error: %w", err)
	}

	// export the result back to standard Go types
	exported := resultVal.Export()

	result, ok := exported.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid return type: expected a JSON object, got %T", exported)
	}

	return result, nil
}

func JSExample() {
	jsScript := `
		function process(name, age, tags, metadata) {
			return {
				message: "Hello " + name + ", age " + age,
				tags: tags,
				meta: metadata,
				status: 200,
				active: true
			};
		}
	`

	// Arbitrary Go types
	tags := []string{"admin", "user"}
	meta := map[string]any{"session": "xyz123", "active": true}

	data, err := RunJavaScript(jsScript, "Peter", 30, tags, meta)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("JS Message: %s (Status: %v)\n", data["message"], data["status"])
	fmt.Printf("Tags: %v\n", data["tags"])
	fmt.Printf("Meta: %v\n", data["meta"])
}
