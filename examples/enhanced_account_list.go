package main

import (
	"context"
	"fmt"
	"log"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	fmt.Println("=== 增强的账户列表功能演示 ===")

	clientID := "your-client-id-here"
	
	// 快速设置客户端
	client, _, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Fatalf("客户端设置失败: %v", err)
	}

	ctx := context.Background()

	// 1. 基本列表功能
	fmt.Println("\n1. 基本账户列表")
	basicOpts := &interlace.AccountListOptions{
		Limit: 5,
		Page:  1,
	}

	accountList, err := client.Account.List(ctx, basicOpts)
	if err != nil {
		log.Printf("获取账户列表失败: %v", err)
	} else {
		fmt.Printf("✓ 找到 %s 个账户，显示前 5 个:\n", accountList.Total)
		for i, account := range accountList.List {
			fmt.Printf("  %d. %s (ID: %s, 状态: %s, 类型: %d)\n", 
				i+1, account.VerifiedName, account.DisplayID, account.Status, account.Type)
		}
	}

	// 2. 获取所有活跃账户
	fmt.Println("\n2. 所有活跃账户")
	activeAccounts, err := client.Account.ListActiveAccounts(ctx)
	if err != nil {
		log.Printf("获取活跃账户失败: %v", err)
	} else {
		fmt.Printf("✓ 找到 %d 个活跃账户:\n", len(activeAccounts))
		for i, account := range activeAccounts {
			fmt.Printf("  %d. %s (状态: %s)\n", i+1, account.VerifiedName, account.Status)
		}
	}

	// 3. 按状态过滤
	fmt.Println("\n3. 按状态过滤账户")
	statusOpts := &interlace.AccountListOptions{
		Status: interlace.AccountStatusActive,
		Limit:  10,
		Page:   1,
	}

	statusFiltered, err := client.Account.List(ctx, statusOpts)
	if err != nil {
		log.Printf("按状态过滤失败: %v", err)
	} else {
		fmt.Printf("✓ 找到 %d 个 ACTIVE 状态的账户\n", len(statusFiltered.List))
	}

	// 4. 按类型过滤
	fmt.Println("\n4. 按类型过滤账户")
	typeOpts := &interlace.AccountListOptions{
		Type:  interlace.AccountTypePersonal,
		Limit: 10,
		Page:  1,
	}

	typeFiltered, err := client.Account.List(ctx, typeOpts)
	if err != nil {
		log.Printf("按类型过滤失败: %v", err)
	} else {
		fmt.Printf("✓ 找到 %d 个个人账户 (类型=1)\n", len(typeFiltered.List))
	}

	// 5. 获取账户总数
	fmt.Println("\n5. 账户总数统计")
	totalCount, err := client.Account.Count(ctx)
	if err != nil {
		log.Printf("获取账户总数失败: %v", err)
	} else {
		fmt.Printf("✓ 账户总数: %d\n", totalCount)
	}

	// 6. 分页获取
	fmt.Println("\n6. 分页获取账户")
	pageData, err := client.Account.GetAccountsByPage(ctx, 1, 3)
	if err != nil {
		log.Printf("分页获取失败: %v", err)
	} else {
		fmt.Printf("✓ 第 1 页，每页 3 条，总计 %s 条:\n", pageData.Total)
		for i, account := range pageData.List {
			fmt.Printf("  %d. %s\n", i+1, account.VerifiedName)
		}
	}

	// 7. 获取所有账户（自动处理分页）
	fmt.Println("\n7. 获取所有账户（自动分页）")
	allAccounts, err := client.Account.ListAll(ctx)
	if err != nil {
		log.Printf("获取所有账户失败: %v", err)
	} else {
		fmt.Printf("✓ 成功获取所有 %d 个账户\n", len(allAccounts))
	}

	fmt.Println("\n=== 功能总结 ===")
	fmt.Println("✅ List() - 基础列表功能，支持分页和过滤")
	fmt.Println("✅ ListAll() - 获取所有账户，自动处理分页")
	fmt.Println("✅ ListActiveAccounts() - 获取活跃账户")
	fmt.Println("✅ ListInactiveAccounts() - 获取非活跃账户")
	fmt.Println("✅ ListByStatus() - 按状态过滤")
	fmt.Println("✅ ListByType() - 按类型过滤")
	fmt.Println("✅ Count() - 获取账户总数")
	fmt.Println("✅ GetAccountsByPage() - 分页获取")
	fmt.Println("✅ Get() - 获取特定账户")

	fmt.Println("\n=== 支持的过滤参数 ===")
	fmt.Println("• accountId - 特定账户ID")
	fmt.Println("• limit - 每页数量 (1-100)")
	fmt.Println("• page - 页码")
	fmt.Println("• status - 账户状态 (ACTIVE, INACTIVE, PENDING, SUSPENDED)")
	fmt.Println("• type - 账户类型 (1=个人, 2=企业, 3=子账户)")
}