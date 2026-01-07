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
	// 初始化模型
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:   os.Getenv("OPENAI_API_KEY"),
		Model:    os.Getenv("OPENAI_MODEL_NAME"),
		Timeout:  &timeout,
		Thinking: &arkModel.Thinking{Type: arkModel.ThinkingTypeEnabled},
	})
	if err != nil {
		panic(err)
	}
	return model, nil
}
