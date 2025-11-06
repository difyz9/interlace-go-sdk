package main

import (
	"context"
	"fmt"
	"log"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	fmt.Println("=== 卡片管理功能演示 ===")

	clientID := "your-client-id-here"
	
	// 快速设置客户端
	client, _, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Fatalf("客户端设置失败: %v", err)
	}

	ctx := context.Background()

	// 1. 获取所有卡片列表
	fmt.Println("\n1. 获取所有卡片")
	cards, err := client.Card.ListAllCards(ctx)
	if err != nil {
		log.Printf("获取卡片列表失败: %v", err)
		
		// 如果没有卡片或API调用失败，演示其他功能
		fmt.Println("\n演示卡片管理功能结构...")
		demonstrateCardFeatures()
		return
	}

	fmt.Printf("✓ 找到 %d 张卡片\n", len(cards))

	if len(cards) == 0 {
		fmt.Println("当前没有卡片，演示基本功能...")
		demonstrateCardFeatures()
		return
	}

	// 显示卡片基本信息
	for i, card := range cards {
		fmt.Printf("   卡片 %d:\n", i+1)
		fmt.Printf("     ID: %s\n", card.ID)
		fmt.Printf("     账户 ID: %s\n", card.AccountID)
		fmt.Printf("     类型: %s\n", card.CardType)
		fmt.Printf("     状态: %s\n", card.CardStatus)
		fmt.Printf("     卡号后四位: %s\n", card.Last4Digits)
		fmt.Printf("     到期日期: %s/%s\n", card.ExpiryMonth, card.ExpiryYear)
		fmt.Printf("     是否激活: %v\n", card.IsActive)
		if card.Balance != nil {
			fmt.Printf("     余额: %.2f %s\n", *card.Balance, card.Currency)
		}
		if card.CreditLimit != nil {
			fmt.Printf("     信用额度: %.2f %s\n", *card.CreditLimit, card.Currency)
		}
		fmt.Println()
	}

	// 2. 按状态筛选卡片
	fmt.Println("2. 按状态筛选卡片")
	activeCards, err := client.Card.ListActiveCards(ctx)
	if err != nil {
		log.Printf("获取活跃卡片失败: %v", err)
	} else {
		fmt.Printf("   ✓ 活跃卡片数量: %d\n", len(activeCards))
	}

	// 按不同状态筛选
	statusList := []string{
		interlace.CardStatusActive,
		interlace.CardStatusInactive,
		interlace.CardStatusBlocked,
	}

	for _, status := range statusList {
		statusCards, err := client.Card.ListCardsByStatus(ctx, status)
		if err != nil {
			log.Printf("获取状态为 %s 的卡片失败: %v", status, err)
			continue
		}
		fmt.Printf("   %s 状态卡片: %d 张\n", status, len(statusCards))
	}

	// 3. 按类型筛选卡片
	fmt.Println("\n3. 按类型筛选卡片")
	cardTypes := []string{
		interlace.CardTypeVirtual,
		interlace.CardTypePhysical,
		interlace.CardTypePrepaid,
		interlace.CardTypeCredit,
	}

	for _, cardType := range cardTypes {
		typeCards, err := client.Card.ListCardsByType(ctx, cardType)
		if err != nil {
			log.Printf("获取类型为 %s 的卡片失败: %v", cardType, err)
			continue
		}
		fmt.Printf("   %s 类型卡片: %d 张\n", cardType, len(typeCards))
	}

	// 4. 分页查询演示
	fmt.Println("\n4. 分页查询演示")
	pageResponse, err := client.Card.GetCardsByPage(ctx, 1, 5)
	if err != nil {
		log.Printf("分页查询失败: %v", err)
	} else {
		fmt.Printf("   第1页 (每页5张): 返回 %d 张卡片\n", len(pageResponse.Cards))
		fmt.Printf("   总数: %d, 是否还有更多: %v\n", pageResponse.TotalCount, pageResponse.HasMore)
	}

	// 5. 获取卡片详细信息（包含敏感信息）
	if len(cards) > 0 {
		firstCard := cards[0]
		fmt.Printf("\n5. 获取卡片详细信息 (卡片ID: %s)\n", firstCard.ID)
		
		privateInfo, err := client.Card.GetCardPrivateInfo(ctx, firstCard.ID)
		if err != nil {
			log.Printf("获取卡片详细信息失败: %v", err)
		} else {
			fmt.Printf("   ✓ 卡片详细信息获取成功\n")
			fmt.Printf("   持卡人姓名: %s\n", privateInfo.CardholderName)
			fmt.Printf("   卡号前6位: %s\n", privateInfo.CardBIN)
			fmt.Printf("   卡号后4位: %s\n", privateInfo.Last4Digits)
			fmt.Printf("   状态: %s\n", privateInfo.CardStatus)
			fmt.Printf("   是否激活: %v\n", privateInfo.IsActive)
			fmt.Printf("   注意: 完整卡号和CVV已加密\n")
			if privateInfo.CardNumber != "" {
				fmt.Printf("   加密卡号长度: %d 字符\n", len(privateInfo.CardNumber))
			}
			if privateInfo.CVV != "" {
				fmt.Printf("   加密CVV长度: %d 字符\n", len(privateInfo.CVV))
			}
		}
	}

	// 6. 统计信息演示
	fmt.Println("\n6. 卡片统计信息")
	
	totalCount, err := client.Card.CountCards(ctx)
	if err != nil {
		log.Printf("获取卡片总数失败: %v", err)
	} else {
		fmt.Printf("   总卡片数: %d\n", totalCount)
	}

	activeCount, err := client.Card.CountActiveCards(ctx)
	if err != nil {
		log.Printf("获取活跃卡片数失败: %v", err)
	} else {
		fmt.Printf("   活跃卡片数: %d\n", activeCount)
	}

	// 按账户统计
	if len(cards) > 0 {
		firstAccountID := cards[0].AccountID
		accountCardCount, err := client.Card.CountCardsByAccount(ctx, firstAccountID)
		if err != nil {
			log.Printf("获取账户卡片数失败: %v", err)
		} else {
			fmt.Printf("   账户 %s 的卡片数: %d\n", firstAccountID, accountCardCount)
		}
	}

	// 7. 便捷方法演示
	fmt.Println("\n7. 便捷方法演示")
	
	hasCards, err := client.Card.HasCards(ctx)
	if err != nil {
		log.Printf("检查是否有卡片失败: %v", err)
	} else {
		fmt.Printf("   用户是否有卡片: %v\n", hasCards)
	}

	if len(cards) > 0 {
		firstCard := cards[0]
		
		isActive, err := client.Card.IsCardActive(ctx, firstCard.ID)
		if err != nil {
			log.Printf("检查卡片状态失败: %v", err)
		} else {
			fmt.Printf("   卡片 %s 是否激活: %v\n", firstCard.ID, isActive)
		}

		cardStatus, err := client.Card.GetCardStatus(ctx, firstCard.ID)
		if err != nil {
			log.Printf("获取卡片状态失败: %v", err)
		} else {
			fmt.Printf("   卡片 %s 当前状态: %s\n", firstCard.ID, cardStatus)
		}
	}

	// 8. 卡片删除演示（注意：这是危险操作，生产环境请谨慎使用）
	fmt.Println("\n8. 卡片删除功能说明")
	fmt.Println("   ⚠️  卡片删除是不可逆操作，请谨慎使用")
	fmt.Println("   • 预付卡余额将转移到量子账户")
	fmt.Println("   • 删除后的退款将在T+1基础上转移到量子账户")
	fmt.Println("   • 信用卡退款仍将退回到组织账户")
	
	// 实际删除代码（已注释，避免误删除）
	/*
	if len(cards) > 1 { // 只有在有多张卡片时才演示删除
		lastCard := cards[len(cards)-1]
		fmt.Printf("   如需删除卡片 %s，请使用:\n", lastCard.ID)
		fmt.Printf("   removeResp, err := client.Card.RemoveCard(ctx, \"%s\")\n", lastCard.ID)
		
		// 实际执行删除（取消注释以启用）
		// removeResp, err := client.Card.RemoveCard(ctx, lastCard.ID)
		// if err != nil {
		//     log.Printf("删除卡片失败: %v", err)
		// } else {
		//     fmt.Printf("   ✓ 卡片删除成功\n")
		//     fmt.Printf("   删除时间: %s\n", removeResp.RemovedAt)
		//     if removeResp.BalanceTransfer != nil {
		//         fmt.Printf("   余额转移: %.2f %s -> %s\n", 
		//             removeResp.BalanceTransfer.Amount,
		//             removeResp.BalanceTransfer.Currency,
		//             removeResp.BalanceTransfer.TransferredTo)
		//     }
		// }
	}
	*/
}

func demonstrateCardFeatures() {
	fmt.Println("\n=== 卡片管理功能总览 ===")
	
	fmt.Println("\n📋 卡片列表功能:")
	fmt.Println("   • ListCards() - 基础列表查询（支持筛选和分页）")
	fmt.Println("   • ListAllCards() - 获取所有卡片（自动分页）")
	fmt.Println("   • ListCardsByAccount() - 按账户筛选")
	fmt.Println("   • ListActiveCards() - 获取活跃卡片")
	fmt.Println("   • ListCardsByStatus() - 按状态筛选")
	fmt.Println("   • ListCardsByType() - 按类型筛选")
	
	fmt.Println("\n🔍 卡片查询功能:")
	fmt.Println("   • GetCardPrivateInfo() - 获取敏感信息（加密）")
	fmt.Println("   • IsCardActive() - 检查激活状态")
	fmt.Println("   • GetCardStatus() - 获取卡片状态")
	
	fmt.Println("\n📊 统计功能:")
	fmt.Println("   • CountCards() - 总卡片数")
	fmt.Println("   • CountActiveCards() - 活跃卡片数")
	fmt.Println("   • CountCardsByAccount() - 按账户统计")
	fmt.Println("   • HasCards() - 检查是否有卡片")
	
	fmt.Println("\n📄 分页功能:")
	fmt.Println("   • GetCardsByPage() - 指定页码查询")
	
	fmt.Println("\n🗑️  管理功能:")
	fmt.Println("   • RemoveCard() - 删除卡片（不可逆）")
	
	fmt.Println("\n🏷️  支持的卡片状态:")
	fmt.Printf("   • %s - 激活\n", interlace.CardStatusActive)
	fmt.Printf("   • %s - 未激活\n", interlace.CardStatusInactive)
	fmt.Printf("   • %s - 已阻止\n", interlace.CardStatusBlocked)
	fmt.Printf("   • %s - 已过期\n", interlace.CardStatusExpired)
	fmt.Printf("   • %s - 待处理\n", interlace.CardStatusPending)
	fmt.Printf("   • %s - 已取消\n", interlace.CardStatusCancelled)
	
	fmt.Println("\n🎯 支持的卡片类型:")
	fmt.Printf("   • %s - 虚拟卡\n", interlace.CardTypeVirtual)
	fmt.Printf("   • %s - 实体卡\n", interlace.CardTypePhysical)
	fmt.Printf("   • %s - 预付卡\n", interlace.CardTypePrepaid)
	fmt.Printf("   • %s - 信用卡\n", interlace.CardTypeCredit)
	
	fmt.Println("\n🔐 安全特性:")
	fmt.Println("   • 敏感信息（卡号、CVV）使用AES加密")
	fmt.Println("   • 加密密钥基于clientSecret生成")
	fmt.Println("   • 非敏感信息（BIN、后四位）保持明文")
}