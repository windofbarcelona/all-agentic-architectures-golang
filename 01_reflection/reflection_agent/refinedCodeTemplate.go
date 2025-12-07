package reflectionagent

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func RefinedCodeTemplate() prompt.DefaultChatTemplate {
	systemTpl := `
	角色：你是一个精通golang编程语言的高级软件工程师，同时熟悉计算机编程的所有范式和最佳实践。
	任务：根据初版代码及修改意见对代码进行修改，并返回修改后的代码。
	数据：
		1. 原始代码: {DraftCode}
		2. 修改意见: {Critique}
	要求：
	1. 严格遵照修改意见进行代码修改。
	2. 如果没有修改意见，请返回原始代码。
	`
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		//schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)
	return *chatTpl
}
