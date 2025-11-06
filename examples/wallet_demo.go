package main

import (
	"context"
	"fmt"
	"log"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	fmt.Println("=== 加密钱包管理 API 演示 ===")

	clientID := "your-client-id-here"
	
	// 快速设置客户端
	client, _, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Fatalf("客户端设置失败: %v", err)
	}

	ctx := context.Background()

	// 使用测试账户ID
	accountID := "your-account-id-here"
	fmt.Printf("\n使用测试账户: %s\n", accountID)

	// 1. 创建钱包 (Create Wallet)
	fmt.Println("\n1. 创建加密钱包")
	createReq := &interlace.CreateWalletRequest{
		AccountID:      accountID,
		Nickname:       "My Crypto Wallet",
		IdempotencyKey: fmt.Sprintf("wallet-%d", 1699876543), // 使用时间戳作为幂等性键
	}

	wallet, err := client.Wallet.CreateWallet(ctx, createReq)
	if err != nil {
		log.Printf("创建钱包失败: %v", err)
		// 如果创建失败，我们尝试列出现有的钱包
		fmt.Println("\n尝试列出现有钱包...")
	} else {
		fmt.Printf("✓ 钱包创建成功\n")
		fmt.Printf("   钱包ID: %s\n", wallet.ID)
		fmt.Printf("   昵称: %s\n", wallet.Nickname)
		fmt.Printf("   状态: %s\n", wallet.Status)
		fmt.Printf("   创建时间: %s\n", wallet.CreatedAt)
	}

	// 2. 获取钱包列表 (List Wallets)
	fmt.Println("\n2. 获取钱包列表")
	listOptions := &interlace.WalletListOptions{
		AccountID: accountID,
		Limit:     10,
		Page:      1,
	}

	walletList, err := client.Wallet.ListWallets(ctx, listOptions)
	if err != nil {
		log.Printf("获取钱包列表失败: %v", err)
		return
	}

	fmt.Printf("✓ 获取到 %d 个钱包 (总共 %d 个)\n", len(walletList.Wallets), walletList.TotalCount)
	
	var walletID string
	for i, w := range walletList.Wallets {
		fmt.Printf("   钱包 %d:\n", i+1)
		fmt.Printf("     ID: %s\n", w.ID)
		fmt.Printf("     昵称: %s\n", w.Nickname)
		fmt.Printf("     状态: %s\n", w.Status)
		fmt.Printf("     创建时间: %s\n", w.CreatedAt)
		
		if i == 0 {
			walletID = w.ID
		}
	}

	if walletID == "" {
		fmt.Println("\n暂无钱包可用于演示后续功能")
		return
	}

	// 3. 获取特定钱包 (Get Wallet)
	fmt.Printf("\n3. 获取钱包详情 (ID: %s)\n", walletID)
	walletDetail, err := client.Wallet.GetWallet(ctx, walletID)
	if err != nil {
		log.Printf("获取钱包详情失败: %v", err)
	} else {
		fmt.Printf("✓ 钱包详情:\n")
		fmt.Printf("   ID: %s\n", walletDetail.ID)
		fmt.Printf("   账户ID: %s\n", walletDetail.AccountID)
		fmt.Printf("   昵称: %s\n", walletDetail.Nickname)
		fmt.Printf("   状态: %s\n", walletDetail.Status)
		fmt.Printf("   创建时间: %s\n", walletDetail.CreatedAt)
		fmt.Printf("   更新时间: %s\n", walletDetail.UpdatedAt)
	}

	// 4. 更新钱包 (Update Wallet)
	fmt.Printf("\n4. 更新钱包昵称\n")
	updateReq := &interlace.UpdateWalletRequest{
		Nickname: "Updated Crypto Wallet",
	}

	updatedWallet, err := client.Wallet.UpdateWallet(ctx, walletID, updateReq)
	if err != nil {
		log.Printf("更新钱包失败: %v", err)
	} else {
		fmt.Printf("✓ 钱包更新成功\n")
		fmt.Printf("   新昵称: %s\n", updatedWallet.Nickname)
		fmt.Printf("   更新时间: %s\n", updatedWallet.UpdatedAt)
	}

	// 5. 创建区块链地址 (Create Blockchain Address)
	fmt.Printf("\n5. 为钱包创建区块链地址\n")
	addressReq := &interlace.CreateAddressRequest{
		Currency: "USDT",
		Chain:    "TRC20",
	}

	address, err := client.Wallet.CreateWalletAddress(ctx, walletID, addressReq)
	if err != nil {
		log.Printf("创建区块链地址失败: %v", err)
	} else {
		fmt.Printf("✓ 区块链地址创建成功\n")
		fmt.Printf("   地址ID: %s\n", address.ID)
		fmt.Printf("   币种: %s\n", address.Currency)
		fmt.Printf("   链: %s\n", address.Chain)
		fmt.Printf("   地址: %s\n", address.Address)
		if address.Tag != "" {
			fmt.Printf("   标签: %s\n", address.Tag)
		}
		fmt.Printf("   创建时间: %s\n", address.CreatedAt)
	}

	fmt.Println("\n=== API 功能总结 ===")
	fmt.Println("✅ CreateWallet() - 创建加密钱包")
	fmt.Println("   - 为指定账户创建子钱包")
	fmt.Println("   - 可选设置钱包昵称")
	fmt.Println("")
	fmt.Println("✅ ListWallets() - 获取钱包列表")
	fmt.Println("   - 按accountId筛选")
	fmt.Println("   - 支持分页 (limit, page)")
	fmt.Println("")
	fmt.Println("✅ GetWallet() - 获取钱包详情")
	fmt.Println("   - 通过钱包ID获取完整信息")
	fmt.Println("")
	fmt.Println("✅ UpdateWallet() - 更新钱包昵称")
	fmt.Println("   - 修改钱包的显示名称")
	fmt.Println("")
	fmt.Println("✅ CreateWalletAddress() - 创建区块链地址")
	fmt.Println("   - 为钱包生成特定币种/链的地址")
	fmt.Println("   - 支持多种币种和链 (USDT/TRC20, BTC, ETH等)")
}