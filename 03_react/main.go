package react

import (
	"context"
	"fmt"
	"log"
	"reactAgent/react"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	// Init eino devops server

	if runnable, err := react.GetToolUseRunnable(); err != nil {
		panic(err)
	} else {
		if res, err := runnable.Invoke(ctx, map[string]any{
			"user_query": "规划一条从鸟巢到北京西站，骑自行车的路线",
		}, compose.WithCallbacks(&loggerCallbacks{})); err != nil {
			panic(err)
		} else {
			fmt.Println("res:", res)
		}
	}
}

type loggerCallbacks struct{}

func (l *loggerCallbacks) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	log.Printf("OnStart name: %v, type: %v, component: %v, input: %v\n\n", info.Name, info.Type, info.Component, input)
	return ctx
}

func (l *loggerCallbacks) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	log.Printf("OnEnd name: %v, type: %v, component: %v, output: %v\n\n", info.Name, info.Type, info.Component, output)
	return ctx
}

func (l *loggerCallbacks) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	log.Printf("name: %v, type: %v, component: %v, error: %v", info.Name, info.Type, info.Component, err)
	return ctx
}

func (l *loggerCallbacks) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo, input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	return ctx
}

func (l *loggerCallbacks) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
	return ctx
}
