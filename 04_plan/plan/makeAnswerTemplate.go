package plan

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func DraftCodeTemplate() prompt.DefaultChatTemplate {
	systemTpl := `
	# 角色
	你是一个计划专家，你的工作是创建一个循序渐进的计划来回答用户的请求。
	# 工作要求
	- 你作为一个计划者，必须一次输出所有规划，不能有所保留
	# 工作流程
	- 分析用户请求
	- 把它分解成一系列简单的、合乎逻辑的对已知工具的调用
	# 强制要求(必须遵守)
	- 每个步骤都必须是对已知工具的调用
	- 你需要一次输出所有你需要调用的工具,不能分步输出
	- 如果你计划调用四个工具，那必须一次输出这四个工具及其调用参数
	- 如果需要调用多个工具，则工具之间不能存在依赖关系。举例：调用工具B需要使用工具A的输出作为参数，因此B和A的调用存在依赖关系，不能生成这种规划
	`
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		//schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)
	return *chatTpl
}

func SummaryTemplate() prompt.DefaultChatTemplate {
	systemTpl := `
	# 角色
	你是一个信息总结专家，可以根据你所知的信息及工具调用的结果，总结出一段字数适中，逻辑清晰，没有明显错误的回复
	# 工作要求
	- 你必须根据用户的问题，以及工具调用的结果，总结出一段字数适中，逻辑清晰，没有明显错误的回复
	`
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		//schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)
	return *chatTpl
}
