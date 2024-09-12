package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/yomorun/yomo/serverless"
)

// Implement DataTags() to observe data with the given tags
func DataTags() []uint32 {
	return []uint32{0x12}
}

// Implement Init() for state initialization, such as loading LLM Model to GPU memory.
func Init() error {
	return nil
}

// Parameters needed for OpenAI Function Calling
// ref: https://platform.openai.com/docs/guides/function-calling
type Parameter struct {
	State bool `json:"state" jsonschema:"description=To set the PLC state to True or False"`
}

// Implement Description() to define the description of OpenAI Function Calling
// ref: https://platform.openai.com/docs/guides/function-calling
func Description() string {
	return "A function that sets PLC state."
}

// Implement InputSchema() to define the input schema of the function
func InputSchema() any {
	return &Parameter{}
}

// Implement Handler() to handle the function call
func Handler(ctx serverless.Context) {
	fmt.Println("start running handler")
	ch := make(chan string)

	go func() {
		var msg Parameter
		err := ctx.ReadLLMArguments(&msg)
		if err != nil {
			ch <- "error: " + err.Error()
			return
		}

		value := 0
		if msg.State {
			value = 1
		}

		for i := 0; i < 4; i++ {
			url := fmt.Sprintf("http://localhost:30080/deviceshifu-plc/sendsinglebit?rootaddress=Q&address=0&start=0&digit=%d&value=%d", i, value)

			_, err = httpGet(url)
			if err != nil {
				ch <- "error: " + err.Error()
				return
			}
		}

		ch <- "success"
	}()

	for res := range ch {
		fmt.Println("res:", res)
		ctx.WriteLLMResult(res)
	}
}

func httpGet(u string) ([]byte, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %s", resp.Status)
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
