package plan

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/joho/godotenv"
	arkModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

func GetModel() (model.ToolCallingChatModel, error) {
	err := godotenv.Load("../model.env")
	if err != nil {
		fmt.Println("Error loading.env file:", err)
		return nil, err
	}
	//openAIBaseURL := os.Getenv("OPENAI_BASE_URL")
	ctx := context.Background()

	timeout := 300 * time.Second
	type OutputFormat struct {
		Params   string `json:"tool_params"`
		ToolName string `json:"tool_name"`
	}
	//outputFormatSchema, err := openapi3gen.NewSchemaRefForValue(&OutputFormat{}, nil)
	// 初始化模型
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:   os.Getenv("OPENAI_API_KEY"),
		Model:    os.Getenv("OPENAI_MODEL_NAME"),
		Timeout:  &timeout,
		Thinking: &arkModel.Thinking{Type: arkModel.ThinkingTypeEnabled},
		// ResponseFormat: &ark.ResponseFormat{
		// 	Type: arkModel.ResponseFormatJsonObject,
		// 	JSONSchema: &arkModel.ResponseFormatJSONSchemaJSONSchemaParam{
		// 		Name:        "工具调用输出",
		// 		Description: "调用工具所需的名称和参数",
		// 		Schema:      outputFormatSchema,
		// 	},
		// },
	})
	if err != nil {
		panic(err)
	}
	return model, nil
}
