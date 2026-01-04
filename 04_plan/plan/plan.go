package plan

import (
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type state struct {
	Messages []*schema.Message
}

func GetPlanRunnable() (compose.Runnable[map[string]any, *schema.Message], error) {
	return nil, nil
}
