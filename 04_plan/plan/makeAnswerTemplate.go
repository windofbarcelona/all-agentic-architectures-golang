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
	1. 你需要先分解用户请求，总结出需要解决用户问题你需要的前置信息和步骤
	2. 最后需要输出一个json字符串,其中按顺序包括你想要调用的工具名称及调用该工具的参数
	实例：
	1. 用户请求：路人甲在北京西站，路人乙在北京鸟巢，找一个距离两人出行时间接近的商场，要求有火锅店；
	 比如：你首先需要知道路人甲和路人乙距离接近的商场都有哪些，其次你需要知道这些商场是否都有火锅店。
	`
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		//schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)
	return *chatTpl
}
