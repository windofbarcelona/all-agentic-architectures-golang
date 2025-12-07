package reflectionagent

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func CritiqueTemplate() prompt.DefaultChatTemplate {
	systemTpl := `
	角色:你是一个精通golang编程语言的高级测试工程师,擅长寻找代码中的bug并给出修复意见。
	任务:根据用户给出的代码,查找其中的bug及不够健壮的部分,给出修改建议。
	数据：
		1. 原始代码: {DraftCode}
	举例：
		1. 比如未能兼容某些参数的边界情况，比如入参可能为0、负数、超大值、空字符串等。
		2. 比如未能处理可能出现的错误情况，比如网络请求失败、文件读写异常等。
		3. 比如存在潜在的性能问题，比如算法复杂度过高、频繁的内存分配等。
	`
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		//schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("用户的原始需求是：{user_query}"),
	)
	return *chatTpl
}
