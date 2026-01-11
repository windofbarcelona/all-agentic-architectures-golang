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
	Query    string
}

func GetPlanRunnable() (compose.Runnable[map[string]any, *schema.Message], error) {
	// init
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
	toolsNodePostHandle := func(ctx context.Context, input []*schema.Message, state *state) ([]*schema.Message, error) {
		state.Messages = append(state.Messages, input...)
		return input, nil
	}

	makeAnswerTemplate := DraftCodeTemplate()
	synthesisTemplate := SummaryTemplate()

	// add all nodes
	sg.AddChatTemplateNode("MakeAnswerTemplate", &makeAnswerTemplate, compose.WithNodeName("MakeAnswerTemplate"), compose.WithStatePreHandler(func(ctx context.Context, in map[string]any, state *state) (map[string]any, error) {
		state.Query = in["user_query"].(string)
		return in, nil
	}))
	sg.AddChatModelNode("MakeAnswerModel", model, compose.WithNodeName("MakeAnswerModel"), compose.WithStatePreHandler(modelPreHandle))
	sg.AddToolsNode("ToolsNode", toolNode, compose.WithNodeName("ToolsNode"), compose.WithStatePreHandler(toolsNodePreHandle), compose.WithStatePostHandler(toolsNodePostHandle))
	sg.AddLambdaNode("parseToolNode", compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) (output map[string]any, err error) {
		output = make(map[string]any, 1)
		return output, nil
	}), compose.WithStatePostHandler(func(ctx context.Context, output map[string]any, state *state) (map[string]any, error) {
		output["user_query"] = state.Query
		return output, nil
	}))
	sg.AddChatTemplateNode("MakeSynthesisTemplate", &synthesisTemplate, compose.WithNodeName("MakeSynthesisTemplate"))
	sg.AddChatModelNode("Synthesis", model, compose.WithNodeName("Synthesis"), compose.WithStatePreHandler(modelPreHandle))

	// add edges
	sg.AddEdge(compose.START, "MakeAnswerTemplate")
	sg.AddEdge("MakeAnswerTemplate", "MakeAnswerModel")
	sg.AddEdge("MakeAnswerModel", "ToolsNode")
	sg.AddEdge("ToolsNode", "parseToolNode")
	sg.AddEdge("parseToolNode", "MakeSynthesisTemplate")
	sg.AddEdge("MakeSynthesisTemplate", "Synthesis")
	sg.AddEdge("Synthesis", compose.END)

	//compileOpts := []compose.GraphCompileOption{compose.WithMaxRunSteps(20)}
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
