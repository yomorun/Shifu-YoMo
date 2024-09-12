package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yomorun/yomo/serverless"
)

// Implement DataTags() to observe data with the given tags
func DataTags() []uint32 {
	return []uint32{0x13}
}

// Implement Init() for state initialization, such as loading LLM Model to GPU memory.
func Init() error {
	return nil
}

// Parameters needed for OpenAI Function Calling
// ref: https://platform.openai.com/docs/guides/function-calling
type Parameter struct {
	Number int `json:"number" jsonschema:"description=The display number on the LED, between 0 and 9999"`
}

// Implement Description() to define the description of OpenAI Function Calling
// ref: https://platform.openai.com/docs/guides/function-calling
func Description() string {
	return "A function that sets the display number of the LED."
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

		_, err = httpPost(
			"http://localhost:30080/deviceshifu-led/number",
			&Req{
				Value: msg.Number,
			},
		)
		if err != nil {
			ch <- "error: " + err.Error()
			return
		}

		ch <- "success"
	}()

	for res := range ch {
		fmt.Println("res:", res)
		ctx.WriteLLMResult(res)
	}
}

type Req struct {
	Value int `json:"value"`
}

func httpPost(url string, req any) ([]byte, error) {
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %s", resp.Status)
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
