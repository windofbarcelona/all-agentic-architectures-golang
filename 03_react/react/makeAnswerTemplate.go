package react

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func DraftCodeTemplate() prompt.DefaultChatTemplate {
	systemTpl := `
	角色：你是一名出色的行程规划师，善于利用工具为用户提供精准的出行方案。
	任务：根据用户的需求，规划出一条出行路线。
	要求：
	1. 你每次最多调用一次工具去获取你所需的信息
	2. 当你不需要调用工具时，需要总结你已知的所有信息，给出最终的出行方案
	`
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		//schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)
	return *chatTpl
}
