package reflectionagent

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type ReflectionState struct {
	DraftCode   string
	Critique    string
	RefinedCode string
	query       string
}

func buildReflectionRunnable() (compose.Runnable[map[string]any, *schema.Message], error) {
	sg := compose.NewGraph[map[string]any, *schema.Message](compose.WithGenLocalState(func(ctx context.Context) (state ReflectionState) {
		return ReflectionState{}
	}))
	temp := DraftCodeTemplate()
	model, err := GetModel()
	if err != nil {
		return nil, err
	}
	sg.AddChatTemplateNode("DraftCodeTemplate", &temp, compose.WithStatePreHandler(func(ctx context.Context, in map[string]any, state ReflectionState) (map[string]any, error) {
		//fmt.Println("DraftCodeTemplate input:", in)
		state.query = in["user_query"].(string)
		return in, nil
	}))
	sg.AddChatModelNode("DraftCodeModel", model, compose.WithStatePostHandler(func(ctx context.Context, out *schema.Message, state ReflectionState) (*schema.Message, error) {
		state.DraftCode = out.Content
		fmt.Println("DraftCodeNode output:", out)
		return out, nil
	}))

	sg.AddLambdaNode("parseDraftCodeTemplate", compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (output map[string]any, err error) {
		// some logic
		output = make(map[string]any, 2)
		output["DraftCode"] = input.Content
		return output, nil
	}), compose.WithStatePostHandler(func(ctx context.Context, output map[string]any, state ReflectionState) (map[string]any, error) {
		output["user_query"] = state.query
		return output, nil
	}))

	temp_2 := CritiqueTemplate()
	model2, err := GetModel()
	sg.AddChatTemplateNode("CritiqueNodeTemplate", &temp_2)
	sg.AddChatModelNode("CritiqueNodeModel", model2, compose.WithStatePostHandler(func(ctx context.Context, out *schema.Message, state ReflectionState) (*schema.Message, error) {
		state.Critique = out.Content
		fmt.Println("CritiqueNodeModel output:", out)
		return out, nil
	}))

	sg.AddLambdaNode("parseCritiqueNodeTemplate", compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (output map[string]any, err error) {
		// some logic
		output = make(map[string]any, 3)
		output["Critique"] = input.Content
		return output, nil
	}), compose.WithStatePostHandler(func(ctx context.Context, output map[string]any, state ReflectionState) (map[string]any, error) {
		output["user_query"] = state.query
		output["DraftCode"] = state.DraftCode
		return output, nil
	}))

	temp_3 := RefinedCodeTemplate()
	model3, err := GetModel()
	sg.AddChatTemplateNode("RefinedCodeTemplate", &temp_3)
	sg.AddChatModelNode("RefinedCodeModel", model3, compose.WithStatePostHandler(func(ctx context.Context, out *schema.Message, state ReflectionState) (*schema.Message, error) {
		state.DraftCode = out.Content
		fmt.Println("RefinedCodeModel output:", out)
		return out, nil
	}))

	sg.AddEdge(compose.START, "DraftCodeTemplate")
	sg.AddEdge("DraftCodeTemplate", "DraftCodeModel")
	//sg.AddEdge("DraftCodeModel", compose.END)
	sg.AddEdge("DraftCodeModel", "parseDraftCodeTemplate")
	sg.AddEdge("parseDraftCodeTemplate", "CritiqueNodeTemplate")
	sg.AddEdge("CritiqueNodeTemplate", "CritiqueNodeModel")
	sg.AddEdge("CritiqueNodeModel", "parseCritiqueNodeTemplate")
	sg.AddEdge("parseCritiqueNodeTemplate", "RefinedCodeTemplate")
	sg.AddEdge("RefinedCodeTemplate", "RefinedCodeModel")
	sg.AddEdge("RefinedCodeModel", compose.END)

	reflectionRunnable, err := sg.Compile(context.Background())
	return reflectionRunnable, err
}

func GetRunnable() (compose.Runnable[map[string]any, *schema.Message], error) {
	reflectionRunnable, err := buildReflectionRunnable()
	return reflectionRunnable, err
}
