package plan

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func DraftCodeTemplate() prompt.DefaultChatTemplate {
	systemTpl := `
	角色:你是一个计划专家
	任务：你的工作是创建一个循序渐进的计划来回答用户的请求。
	要求：
	1. 你每次最多调用一个工具，并且该工具的信息必须是你已知的
	2. 你需要按顺序输出你想要调用的工具名称及调用该工具的参数，
		举例如下：
		{
			{
				"toolName":"工具名1",
				"arg":"参数对应的json字符串"
			},
		}
	`
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		//schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)
	return *chatTpl
}
