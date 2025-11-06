package main

import (
	"context"
	"fmt"
	"log"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	fmt.Println("=== 卡 BIN 管理 API 演示 ===")

	clientID := "your-client-id-here"
	
	// 快速设置客户端
	client, _, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Fatalf("客户端设置失败: %v", err)
	}

	ctx := context.Background()

	// 测试账户ID
	accountID := "your-account-id-here"

	// 1. 列出所有可用的卡 BIN (List Card BINs)
	fmt.Println("\n1. 列出所有可用的卡 BIN")
	
	cardBinsResp, err := client.CardBin.ListCardBins(ctx, accountID)
	if err != nil {
		log.Printf("列出卡 BIN 失败: %v", err)
		fmt.Println("\n提示: 该接口用于获取所有可用的卡 BIN 信息")
	} else {
		cardBins := cardBinsResp.List
		fmt.Printf("✓ 可用卡 BIN 数量: %d\n", len(cardBins))
		if len(cardBins) > 0 {
			fmt.Println("\n卡 BIN 详情:")
			for i, bin := range cardBins {
				fmt.Printf("   %d. ID: %s\n", i+1, bin.ID)
				fmt.Printf("      BIN: %s\n", bin.Bin)
				fmt.Printf("      品牌: %s\n", bin.Brand)
				fmt.Printf("      类型: %s\n", bin.Type)
				fmt.Printf("      币种: %s\n", bin.Currency)
				fmt.Printf("      状态: %s\n", bin.Status)
				if bin.Description != "" {
					fmt.Printf("      描述: %s\n", bin.Description)
				}
				fmt.Printf("      创建时间: %s\n", bin.CreatedAt)
				if i < len(cardBins)-1 {
					fmt.Println()
				}
			}
		} else {
			fmt.Println("   暂无可用的卡 BIN")
		}
	}

	// 2. 列出维护中的卡 BIN (List Card BINs Under Maintenance)
	fmt.Println("\n2. 列出维护中的卡 BIN")
	
	maintainBinsResp, err := client.CardBin.ListCardBinsMaintain(ctx, accountID)
	if err != nil {
		log.Printf("列出维护中的卡 BIN 失败: %v", err)
		fmt.Println("\n提示: 该接口用于获取正在维护中的卡 BIN")
		fmt.Println("   维护中的 BIN 可能会影响交易")
	} else {
		maintainBins := maintainBinsResp.List
		fmt.Printf("✓ 维护中的卡 BIN 数量: %d\n", len(maintainBins))
		if len(maintainBins) > 0 {
			fmt.Println("\n维护中的卡 BIN 详情:")
			for i, bin := range maintainBins {
				fmt.Printf("   %d. ID: %s\n", i+1, bin.ID)
				fmt.Printf("      BIN: %s\n", bin.Bin)
				fmt.Printf("      品牌: %s\n", bin.Brand)
				fmt.Printf("      类型: %s\n", bin.Type)
				fmt.Printf("      状态: %s\n", bin.Status)
				if bin.Description != "" {
					fmt.Printf("      描述: %s\n", bin.Description)
				}
				fmt.Printf("      更新时间: %s\n", bin.UpdatedAt)
				if i < len(maintainBins)-1 {
					fmt.Println()
				}
			}
		} else {
			fmt.Println("   暂无维护中的卡 BIN")
		}
	}

	fmt.Println("\n=== API 功能总结 ===")
	fmt.Println("✅ ListCardBins() - 列出所有可用的卡 BIN")
	fmt.Println("   - 获取所有可用的卡 BIN 信息")
	fmt.Println("   - 包含 BIN 号、品牌、类型、币种等信息")
	fmt.Println("   - 用于创建持卡人或卡片时选择 BIN")
	fmt.Println("")
	fmt.Println("✅ ListCardBinsMaintain() - 列出维护中的卡 BIN")
	fmt.Println("   - 获取正在维护中的卡 BIN 列表")
	fmt.Println("   - 维护中的 BIN 可能影响交易")
	fmt.Println("   - 建议在创建卡片前检查维护状态")
	fmt.Println("")
	fmt.Println("💡 使用场景:")
	fmt.Println("   - 创建持卡人前，先查询可用的 BIN ID")
	fmt.Println("   - 检查特定 BIN 是否在维护中")
	fmt.Println("   - 了解支持的卡品牌和类型")
}
