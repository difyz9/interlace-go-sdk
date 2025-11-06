package main

import (
	"context"
	"fmt"
	"log"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	fmt.Println("=== 卡片管理 API 演示 ===")

	clientID := "your-client-id-here"
	
	// 快速设置客户端
	client, _, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Fatalf("客户端设置失败: %v", err)
	}

	ctx := context.Background()

	// 测试账户ID - 请替换为您自己的账户ID
	accountID := "your-account-id-here"
	fmt.Printf("\n1. 使用测试账户 %s 演示卡片功能\n", accountID)

	// 2. 获取卡片列表 (List Cards)
	fmt.Println("\n2. 获取卡片列表")
	options := &interlace.CardListOptions{
		AccountID: accountID,
		Limit:     10,
		Page:      1,
	}

	cardList, err := client.Card.ListCards(ctx, options)
	if err != nil {
		log.Printf("获取卡片列表失败: %v", err)
	} else {
		fmt.Printf("✓ 获取到 %d 张卡片 (总共 %d 张)\n", len(cardList.Cards), cardList.TotalCount)
		
		for i, card := range cardList.Cards {
			fmt.Printf("   卡片 %d:\n", i+1)
			fmt.Printf("     ID: %s\n", card.ID)
			fmt.Printf("     账户ID: %s\n", card.AccountID)
			fmt.Printf("     类型: %s\n", card.CardType)
			fmt.Printf("     状态: %s\n", card.CardStatus)
			fmt.Printf("     是否激活: %v\n", card.IsActive)
			fmt.Printf("     卡号后四位: %s\n", card.Last4Digits)
			fmt.Printf("     持卡人姓名: %s\n", card.CardholderName)
			
			// 3. 获取卡片私密信息 (Get Card Private Info)
			fmt.Printf("\n   获取卡片 %s 的私密信息:\n", card.ID)
			privateInfo, err := client.Card.GetCardPrivateInfo(ctx, card.ID)
			if err != nil {
				fmt.Printf("     ❌ 获取私密信息失败: %v\n", err)
			} else {
				fmt.Printf("     ✓ 卡号: %s (加密)\n", privateInfo.CardNumber)
				fmt.Printf("     ✓ CVV: %s (加密)\n", privateInfo.CVV)
				fmt.Printf("     ✓ 过期日期: %s/%s\n", privateInfo.ExpiryMonth, privateInfo.ExpiryYear)
				fmt.Printf("     ✓ BIN: %s (明文)\n", privateInfo.CardBIN)
				fmt.Printf("     ✓ 后四位: %s (明文)\n", privateInfo.Last4Digits)
			}
			
			// 注意：实际环境中不要删除卡片，这里只做演示
			if i == 0 { // 只对第一张卡片做删除演示（但不实际执行）
				fmt.Printf("\n   删除卡片演示 (仅展示API调用，不实际执行):\n")
				fmt.Printf("     ⚠️  删除操作不可恢复\n")
				fmt.Printf("     ⚠️  预付卡余额将转入Quantum账户\n")
				fmt.Printf("     ⚠️  删除后的退款将在T+1基础上转入Quantum账户\n")
				fmt.Printf("     API调用: client.Card.RemoveCard(ctx, \"%s\")\n", card.ID)
			}
		}
	}

	// 4. 演示带筛选条件的卡片列表
	fmt.Println("\n4. 演示筛选功能")
	
	// 筛选激活的卡片
	activeFilter := true
	activeOptions := &interlace.CardListOptions{
		AccountID: accountID,
		IsActive:  &activeFilter,
		Limit:     5,
		Page:      1,
	}
	
	activeCards, err := client.Card.ListCards(ctx, activeOptions)
	if err != nil {
		log.Printf("获取激活卡片失败: %v", err)
	} else {
		fmt.Printf("✓ 激活的卡片: %d 张\n", len(activeCards.Cards))
	}

	// 按卡片类型筛选 (如果有虚拟卡)
	virtualOptions := &interlace.CardListOptions{
		AccountID: accountID,
		CardType:  "VIRTUAL",
		Limit:     5,
		Page:      1,
	}
	
	virtualCards, err := client.Card.ListCards(ctx, virtualOptions)
	if err != nil {
		log.Printf("获取虚拟卡失败: %v", err)
	} else {
		fmt.Printf("✓ 虚拟卡: %d 张\n", len(virtualCards.Cards))
	}

	// 按状态筛选
	statusOptions := &interlace.CardListOptions{
		AccountID:  accountID,
		CardStatus: "ACTIVE",
		Limit:      5,
		Page:       1,
	}
	
	statusCards, err := client.Card.ListCards(ctx, statusOptions)
	if err != nil {
		log.Printf("按状态筛选失败: %v", err)
	} else {
		fmt.Printf("✓ 状态为ACTIVE的卡片: %d 张\n", len(statusCards.Cards))
	}

	fmt.Println("\n=== API 功能总结 ===")
	fmt.Println("✅ ListCards() - 获取卡片列表，支持筛选和分页")
	fmt.Println("   - 按accountId筛选")
	fmt.Println("   - 按cardStatus筛选 (ACTIVE, INACTIVE, BLOCKED等)")
	fmt.Println("   - 按cardType筛选 (VIRTUAL, PHYSICAL等)")
	fmt.Println("   - 按isActive筛选")
	fmt.Println("   - 分页参数 (limit, page)")
	fmt.Println("")
	fmt.Println("✅ GetCardPrivateInfo() - 获取卡片敏感信息")
	fmt.Println("   - 卡号和CVV使用AES加密 (密钥为clientSecret)")
	fmt.Println("   - BIN和后四位保持明文")
	fmt.Println("")
	fmt.Println("✅ RemoveCard() - 删除卡片 (不可恢复)")
	fmt.Println("   - 预付卡余额自动转入Quantum账户")
	fmt.Println("   - 删除后退款T+1转入Quantum账户")
	fmt.Println("   - 信用卡退款仍返回原组织")
}