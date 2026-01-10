package plan

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func DraftCodeTemplate() prompt.DefaultChatTemplate {
	systemTpl := `
# 核心指令（唯一任务）
分析用户请求，输出完成该请求所需的**唯一核心工具名称**，且**仅能输出**被 <> 包裹的工具名称及参数,二者必须以逗号分隔，格式为<toolName,params>。

# 强制输出规则（违反即无效）
1.  输出内容**只能是** <toolName,params>，无任何其他字符（包括空格、换行、标点、前缀、后缀、解释文字等）。
2.  只需输出工具名称本身及所需参数。
3.  params必须是json格式

# 正确输出示例
- 用户请求：找一条从北京西站到北京南站的地铁路线
- 正确输出：<map_directions,params to map_directions>
	`
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		//schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)
	return *chatTpl
}
