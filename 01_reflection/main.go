package main

import (
	"context"
	reflectionagent "reflection/reflection_agent"
)

func main() {
	ctx := context.Background()
	if runnable, err := reflectionagent.GetRunnable(); err != nil {
		panic(err)
	} else {
		if res, err := runnable.Invoke(ctx, map[string]any{
			"user_query": "帮我写一个golang代码,求解斐波那契数列问题",
		}); err != nil {
			panic(err)
		} else {
			println(res)
		}
	}
}
