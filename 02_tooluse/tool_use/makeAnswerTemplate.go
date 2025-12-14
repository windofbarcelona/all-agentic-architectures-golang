package tooluse

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func DraftCodeTemplate() prompt.DefaultChatTemplate {
	systemTpl := `
	角色：你是一个精于解决用户问题的问题解答师。
	任务：根据用户的需求，回答问题，必要时使用提供给你的工具进行问题回答。
	要求：如果你已知了工具的查询结果，则无需继续调用工具，而应该结果整合进最终的回答中，确保回答的准确性和完整性。
	注意事项：
	1. map_directions可以用来查询出行路线和耗时
	`
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		//schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)
	return *chatTpl
}
