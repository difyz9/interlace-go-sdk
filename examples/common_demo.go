package main

import (
	"context"
	"fmt"
	"log"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	fmt.Println("=== 通用 API 演示 - 交易场景 ===")

	clientID := "your-client-id-here"
	
	// 快速设置客户端
	client, _, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Fatalf("客户端设置失败: %v", err)
	}

	ctx := context.Background()

	// 测试账户ID
	accountID := "your-account-id-here"

	// 列出所有交易场景 (List Consumption Scenarios)
	fmt.Println("\n列出所有 Infinity Card 交易场景")
	
	scenariosResp, err := client.Common.ListConsumptionScenarios(ctx, accountID)
	if err != nil {
		log.Printf("列出交易场景失败: %v", err)
		fmt.Println("\n提示: 该接口用于获取所有可用的卡交易场景")
		fmt.Println("   交易场景用于定义卡片的使用限制和规则")
	} else {
		scenarios := scenariosResp.List
		fmt.Printf("✓ 交易场景数量: %d\n", len(scenarios))
		if len(scenarios) > 0 {
			fmt.Println("\n交易场景详情:")
			for i, scenario := range scenarios {
				fmt.Printf("   %d. 代码: %s\n", i+1, scenario.Code)
				fmt.Printf("      名称: %s\n", scenario.Name)
				if scenario.Description != "" {
					fmt.Printf("      描述: %s\n", scenario.Description)
				}
				if scenario.Category != "" {
					fmt.Printf("      分类: %s\n", scenario.Category)
				}
				if scenario.Status != "" {
					fmt.Printf("      状态: %s\n", scenario.Status)
				}
				if i < len(scenarios)-1 {
					fmt.Println()
				}
			}
		} else {
			fmt.Println("   暂无交易场景数据")
		}
	}

	fmt.Println("\n=== API 功能总结 ===")
	fmt.Println("✅ ListConsumptionScenarios() - 列出所有交易场景")
	fmt.Println("   - 获取所有可用的卡交易场景")
	fmt.Println("   - 包含场景代码、名称、描述等信息")
	fmt.Println("   - 用于了解卡片支持的交易类型")
	fmt.Println("")
	fmt.Println("💡 使用场景:")
	fmt.Println("   - 创建卡片时了解可用的交易场景")
	fmt.Println("   - 设置卡片交易限制")
	fmt.Println("   - 配置卡片的使用规则")
}
