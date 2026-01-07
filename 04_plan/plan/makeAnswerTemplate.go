package plan

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func DraftCodeTemplate() prompt.DefaultChatTemplate {
	systemTpl := `
	# 角色
	你是一个计划专家，你的工作是创建一个循序渐进的计划来回答用户的请求。
	# 工作流程
	- 分析用户请求
	- 把它分解成一系列简单的、合乎逻辑的对已知工具的调用
	- 输出你思考后的结论
	`
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		//schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)
	return *chatTpl
}
