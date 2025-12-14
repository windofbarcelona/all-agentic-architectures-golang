package tooluse

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type state struct {
	Messages []*schema.Message
}

func GetToolUseRunnable() (compose.Runnable[map[string]any, *schema.Message], error) {
	sg := compose.NewGraph[map[string]any, *schema.Message](compose.WithGenLocalState(func(ctx context.Context) *state {
		return &state{Messages: make([]*schema.Message, 0)}
	}))
	ctx := context.Background()
	model, err := GetModel()
	if err != nil {
		return nil, err
	}
	tools := GetBaiDuMapTool(ctx, []string{MapServer})
	toolNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: tools,
	})
	if err != nil {
		return nil, err
	}
	toolsInfo, err := genToolInfos(ctx, tools)
	if err != nil {
		return nil, err
	}
	model, err = model.WithTools(toolsInfo)
	if err != nil {
		return nil, err
	}
	modelPreHandle := func(ctx context.Context, input []*schema.Message, state *state) ([]*schema.Message, error) {
		state.Messages = append(state.Messages, input...)

		return state.Messages, nil
	}
	toolsNodePreHandle := func(ctx context.Context, input *schema.Message, state *state) (*schema.Message, error) {
		if input == nil {
			return state.Messages[len(state.Messages)-1], nil // used for rerun interrupt resume
		}
		state.Messages = append(state.Messages, input)
		return input, nil
	}

	makeAnswerTemplate := DraftCodeTemplate()
	sg.AddChatTemplateNode("MakeAnswerTemplate", &makeAnswerTemplate, compose.WithNodeName("MakeAnswerTemplate"))
	sg.AddChatModelNode("MakeAnswerModel", model, compose.WithNodeName("MakeAnswerModel"), compose.WithStatePreHandler(modelPreHandle))
	sg.AddToolsNode("ToolsNode", toolNode, compose.WithNodeName("ToolsNode"), compose.WithStatePreHandler(toolsNodePreHandle))
	sg.AddChatModelNode("Synthesis", model, compose.WithNodeName("Synthesis"), compose.WithStatePreHandler(modelPreHandle))

	sg.AddEdge(compose.START, "MakeAnswerTemplate")
	sg.AddEdge("MakeAnswerTemplate", "MakeAnswerModel")
	sg.AddEdge("MakeAnswerModel", "ToolsNode")
	sg.AddEdge("ToolsNode", "Synthesis")
	sg.AddEdge("Synthesis", compose.END)
	reflectionRunnable, err := sg.Compile(context.Background())
	return reflectionRunnable, err
}

func genToolInfos(ctx context.Context, tools []tool.BaseTool) ([]*schema.ToolInfo, error) {
	toolInfos := make([]*schema.ToolInfo, 0, len(tools))
	for _, t := range tools {
		tl, err := t.Info(ctx)
		if err != nil {
			return nil, err
		}
		toolInfos = append(toolInfos, tl)
	}

	return toolInfos, nil
}
