package reflectionagent

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func DraftCodeTemplate() prompt.DefaultChatTemplate {
	systemTpl := `
	角色：你是一个精通golang编程语言的高级软件工程师，同时熟悉计算机编程的所有范式和最佳实践。
	任务：根据用户的需求，编写高质量、可维护且高效的golang代码。
	要求：
	1. 理解用户需求，确保代码满足其功能要求。
	2. 遵循golang的最佳实践和编码规范。
	3. 提供清晰的注释和文档，便于他人理解和维护代码。
	4. 考虑代码的性能和可扩展性，确保其在不同场景下表现良好。
	`
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		//schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)
	return *chatTpl
}
