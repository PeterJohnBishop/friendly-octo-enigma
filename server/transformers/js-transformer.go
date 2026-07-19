package transformers

import (
	"fmt"

	"github.com/dop251/goja"
)

func RunJavaScript(script string, name string) (map[string]interface{}, error) {
	vm := goja.New()

	vm.Set("ClientName", name)

	val, err := vm.RunString(script)
	if err != nil {
		return nil, fmt.Errorf("js execution error: %w", err)
	}

	exported := val.Export()

	result, ok := exported.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid return type: expected a JSON object, got %T", exported)
	}

	return result, nil
}

func JSExample() {
	jsScript := `
		function process() {
			return {
				message: "Hello World, I'm " + ClientName,
				status: 200,
				active: true
			};
		}
		process(); // Returns the object
	`

	data, err := RunJavaScript(jsScript, "Peter")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("JS Message: %s (Status: %v)\n", data["message"], data["status"])
}
