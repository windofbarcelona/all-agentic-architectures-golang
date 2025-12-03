package reflectionagent

import (
	"context"

	"code.byted.org/flow/eino/compose"
)

type ReflectionState struct {
	DraftCode   string
	Critique    string
	RefinedCode string
}

func buildReflectionRunnable() (compose.Runnable[string, string], error) {
	sg := compose.NewGraph[string, string](compose.WithGenLocalState(func(ctx context.Context) (state ReflectionState) {
		return ReflectionState{}
	}))
	reflectionRunnable, err := sg.Compile(context.Background())
	return reflectionRunnable, err
}

func GetRunnable() (compose.Runnable[string, string], error) {
	reflectionRunnable, err := buildReflectionRunnable()
	return reflectionRunnable, err
}
