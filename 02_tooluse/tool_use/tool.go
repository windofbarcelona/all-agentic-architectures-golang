package tooluse

import (
	"context"
	"log"

	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetBaiDuMapTool(ctx context.Context, input []string) []tool.BaseTool {
	toolList := []tool.BaseTool{}
	for _, in := range input {
		cli, err := client.NewSSEMCPClient(in)
		if err != nil {
			log.Fatal(err)
		}
		err = cli.Start(ctx)
		if err != nil {
			log.Fatal(err)
		}

		initRequest := mcp.InitializeRequest{}

		_, err = cli.Initialize(ctx, initRequest)
		if err != nil {
			log.Fatal(err)
		}

		tools, err := mcpp.GetTools(ctx, &mcpp.Config{Cli: cli})
		if err != nil {
			log.Fatal(err)
		}
		toolList = append(toolList, tools...)
	}
	return toolList
}
