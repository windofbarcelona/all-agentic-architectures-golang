package plan

import (
	"context"
	"io"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type state struct {
	Messages []*schema.Message
	Plans    []string
}

func GetPlanRunnable() (compose.Runnable[map[string]any, *schema.Message], error) {
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

	// toolsNodePostHandle := func(ctx context.Context, input *schema.Message, state *state) (*schema.Message, error) {
	// 	state.Messages = append(state.Messages, input)
	// 	return input, nil
	// }

	// modelPostBranchCondition := func(ctx context.Context, sr *schema.StreamReader[*schema.Message]) (endNode string, err error) {
	// 	if isToolCall, err := firstChunkStreamToolCallChecker(ctx, sr); err != nil {
	// 		return "", err
	// 	} else if isToolCall {
	// 		return "ToolsNode", nil
	// 	}
	// 	return compose.END, nil
	// }

	makeAnswerTemplate := DraftCodeTemplate()
	sg.AddChatTemplateNode("MakeAnswerTemplate", &makeAnswerTemplate, compose.WithNodeName("MakeAnswerTemplate"))
	sg.AddChatModelNode("MakeAnswerModel", model, compose.WithNodeName("MakeAnswerModel"), compose.WithStatePreHandler(modelPreHandle))
	sg.AddToolsNode("ToolsNode", toolNode, compose.WithNodeName("ToolsNode"), compose.WithStatePreHandler(toolsNodePreHandle))
	sg.AddChatModelNode("Synthesis", model, compose.WithNodeName("Synthesis"), compose.WithStatePreHandler(modelPreHandle))

	sg.AddEdge(compose.START, "MakeAnswerTemplate")
	sg.AddEdge("MakeAnswerTemplate", "MakeAnswerModel")
	sg.AddEdge("MakeAnswerModel", compose.END)
	//sg.AddEdge("MakeAnswerModel", "ToolsNode")
	// if err = sg.AddBranch("MakeAnswerModel", compose.NewStreamGraphBranch(modelPostBranchCondition, map[string]bool{"ToolsNode": true, compose.END: true})); err != nil {
	// 	return nil, err
	// }
	// sg.AddEdge("ToolsNode", "MakeAnswerModel")
	compileOpts := []compose.GraphCompileOption{compose.WithMaxRunSteps(20)}
	reflectionRunnable, err := sg.Compile(context.Background(), compileOpts...)
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

func firstChunkStreamToolCallChecker(_ context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
	defer sr.Close()

	for {
		msg, err := sr.Recv()
		if err == io.EOF {
			return false, nil
		}
		if err != nil {
			return false, err
		}

		if len(msg.ToolCalls) > 0 {
			return true, nil
		}

		if len(msg.Content) == 0 { // skip empty chunks at the front
			continue
		}

		return false, nil
	}
}
