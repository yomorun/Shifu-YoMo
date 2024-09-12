package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/yomorun/yomo/serverless"
)

var apiKey string

// Implement DataTags() to observe data with the given tags
func DataTags() []uint32 {
	return []uint32{0x11}
}

// Implement Init() for state initialization, such as loading LLM Model to GPU memory.
func Init() error {
	if v, ok := os.LookupEnv("VIVGRID_TOKEN_WITHOUT_TOOLS"); ok {
		apiKey = v
	}
	return nil
}

// Parameters needed for OpenAI Function Calling
// ref: https://platform.openai.com/docs/guides/function-calling
type Parameter struct {
}

// Implement Description() to define the description of OpenAI Function Calling
// ref: https://platform.openai.com/docs/guides/function-calling
func Description() string {
	return "A function that gets current status of the LED display number and PLC state from a virtual capture camera."
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
		body, err := httpGet("http://localhost:30080/deviceshifu-camera/capture")
		if err != nil {
			ch <- "error: " + err.Error()
			return
		}

		image := base64.StdEncoding.EncodeToString(body)

		config := openai.DefaultConfig(apiKey)
		config.BaseURL = "https://openai.vivgrid.com/v1"
		client := openai.NewClientWithConfig(config)

		response, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Messages: []openai.ChatCompletionMessage{
					{
						Role: "user",
						MultiContent: []openai.ChatMessagePart{
							{
								Type: "text",
								Text: "Thanks! Can you tell me what is in the image? Specifically what is the display number on the LED and PLC state(whether it has 4 output lights on)? Please return in json format, like {\"led_display_number\":2929,\"plc_state\":true}.",
							},
							{
								Type: "image_url",
								ImageURL: &openai.ChatMessageImageURL{
									URL: "data:image/jpeg;base64," + image,
								},
							},
						},
					},
				},
			},
		)
		if err != nil {
			ch <- "error: " + err.Error()
			return
		}

		ch <- response.Choices[0].Message.Content

		// ch <- "{\"led_display_number\":2929,\"plc_state\":true}"
	}()

	for res := range ch {
		fmt.Println("res:", res)
		ctx.WriteLLMResult(res)
	}
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %s", resp.Status)
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
