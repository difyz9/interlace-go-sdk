package main

import (
	"context"
	"fmt"
	"log"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	fmt.Println("=== HTTP Client 封装演示 ===")
	
	// 通过 HTTP Client 封装，现在所有的 HTTP 请求都统一处理
	// 不需要在每个方法中重复设置请求头、处理错误等
	
	clientID := "your-client-id-here"
	
	// 1. 创建客户端（内部自动使用封装的 HTTP Client）
	client := interlace.NewClient(nil)
	ctx := context.Background()
	
	fmt.Println("1. HTTP Client 封装的优势：")
	fmt.Println("   - 统一的请求头管理")
	fmt.Println("   - 自动的认证 token 处理")
	fmt.Println("   - 统一的错误处理")
	fmt.Println("   - 减少重复代码")
	
	// 2. OAuth 操作 - 不需要重复设置请求头
	fmt.Println("\n2. OAuth 授权（使用封装后的 HTTP Client）")
	authData, err := client.OAuth.Authorize(ctx, clientID)
	if err != nil {
		log.Printf("授权失败: %v", err)
		return
	}
	
	fmt.Printf("   ✓ 授权成功，代码: %s\n", authData.Code)
	
	// 3. 获取访问令牌
	fmt.Println("\n3. 获取访问令牌")
	tokenData, err := client.OAuth.GetAccessToken(ctx, authData.Code, clientID)
	if err != nil {
		log.Printf("获取令牌失败: %v", err)
		return
	}
	
	fmt.Printf("   ✓ 令牌获取成功\n")
	fmt.Printf("   Access Token: %s\n", tokenData.AccessToken[:20] + "...")
	
	// 4. 设置访问令牌（自动传播到所有子客户端）
	client.SetAccessToken(tokenData.AccessToken)
	fmt.Println("\n4. 设置访问令牌")
	fmt.Println("   ✓ 令牌已自动设置到所有子客户端")
	
	// 5. 账户注册 - 自动包含认证头
	fmt.Println("\n5. 账户注册（自动包含认证头）")
	registerReq := &interlace.AccountRegisterRequest{
		PhoneNumber:      "15900000042",
		Email:            "15900000042@qq.com",
		Name:             "测试用户",
		PhoneCountryCode: "86",
	}
	
	account, err := client.Account.Register(ctx, registerReq)
	if err != nil {
		log.Printf("账户注册失败: %v", err)
		return
	}
	
	fmt.Printf("   ✓ 账户注册成功\n")
	fmt.Printf("   账户 ID: %s\n", account.ID)
	fmt.Printf("   显示 ID: %s\n", account.DisplayID)
	
	// 6. 账户列表 - 自动包含认证头和查询参数处理
	fmt.Println("\n6. 获取账户列表")
	listOpts := &interlace.AccountListOptions{
		Limit: 10,
		Page:  1,
	}
	
	accountList, err := client.Account.List(ctx, listOpts)
	if err != nil {
		log.Printf("获取账户列表失败: %v", err)
		return
	}
	
	fmt.Printf("   ✓ 账户列表获取成功，共 %s 个账户\n", accountList.Total)
	for i, acc := range accountList.List {
		fmt.Printf("   %d. %s (ID: %s)\n", i+1, acc.VerifiedName, acc.DisplayID)
	}
	
	// 7. 展示 HTTP Client 封装带来的好处
	fmt.Println("\n=== HTTP Client 封装的好处 ===")
	fmt.Println("✅ 代码更简洁：每个 API 调用只需关注业务逻辑")
	fmt.Println("✅ 统一管理：请求头、认证、错误处理都在一个地方")
	fmt.Println("✅ 易于维护：修改 HTTP 逻辑只需要在一个地方")
	fmt.Println("✅ 减少重复：不再需要在每个方法中重复写 HTTP 代码")
	fmt.Println("✅ 类型安全：统一的请求选项和响应处理")
	
	// 8. 对比原来的代码
	fmt.Println("\n=== 代码对比 ===")
	fmt.Println("原来的代码（每个方法都需要）：")
	fmt.Println("  - 设置 Accept 头")
	fmt.Println("  - 设置 Content-Type 头")
	fmt.Println("  - 设置 User-Agent 头")
	fmt.Println("  - 设置认证头")
	fmt.Println("  - 创建 HTTP 请求")
	fmt.Println("  - 执行请求")
	fmt.Println("  - 处理响应")
	fmt.Println("  - 解析 JSON")
	fmt.Println("  - 错误处理")
	fmt.Println("")
	fmt.Println("现在的代码：")
	fmt.Println("  - 只需要调用 httpClient.DoPostRequest() 或类似方法")
	fmt.Println("  - 一行代码完成所有 HTTP 操作")
	
	fmt.Println("\n=== 演示完成 ===")
}