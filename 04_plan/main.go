package main

import (
	"context"
	"fmt"
	"log"
	"planAgent/plan"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	// Init eino devops server
	//err := devops.Init(ctx)

	if runnable, err := plan.GetPlanRunnable(); err != nil {
		panic(err)
	} else {
		if res, err := runnable.Invoke(ctx, map[string]any{
			"user_query": "找一条从北京西站到北京南站的地铁路线,并告诉我北京南站附近的火锅店都有哪些",
		}, compose.WithCallbacks(&loggerCallbacks{})); err != nil {
			panic(err)
		} else {
			fmt.Println("res:", res.Content)
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
