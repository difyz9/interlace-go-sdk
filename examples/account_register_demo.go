package main

import (
	"context"
	"fmt"
	"log"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	fmt.Println("=== 新增账户注册函数演示 ===")

	clientID := "your-client-id-here"
	
	// 快速设置客户端
	client, _, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Fatalf("客户端设置失败: %v", err)
	}

	ctx := context.Background()

	fmt.Println("\n1. 使用原有的 Register 函数")
	req := &interlace.AccountRegisterRequest{
		PhoneCountryCode: "86",
		PhoneNumber:      "15900000031",
		Email:            "15900000031@qq.com",
		Name:             "golang_test",
	}

	account1, err := client.Account.Register(ctx, req)
	if err != nil {
		log.Printf("注册失败: %v", err)
	} else {
		fmt.Printf("✓ 账户注册成功 - ID: %s, 显示ID: %s\n", account1.ID, account1.DisplayID)
	}

	fmt.Println("\n2. 使用新增的 RegisterWithDetails 函数")
	account2, err := client.Account.RegisterWithDetails(ctx, "86", "15900000032", "15900000032@qq.com", "golang_test_2")
	if err != nil {
		log.Printf("注册失败: %v", err)
	} else {
		fmt.Printf("✓ 账户注册成功 - ID: %s, 显示ID: %s\n", account2.ID, account2.DisplayID)
	}

	fmt.Println("\n3. 使用便捷的 RegisterGolangTest 函数")
	account3, err := client.Account.RegisterGolangTest(ctx)
	if err != nil {
		log.Printf("注册失败: %v", err)
	} else {
		fmt.Printf("✓ 测试账户注册成功 - ID: %s, 显示ID: %s\n", account3.ID, account3.DisplayID)
	}

	fmt.Println("\n=== 函数对比 ===")
	fmt.Println("1. Register(ctx, req): 通用注册函数，需要构造 AccountRegisterRequest")
	fmt.Println("2. RegisterWithDetails(ctx, phone, number, email, name): 直接传参，更简洁")
	fmt.Println("3. RegisterGolangTest(ctx): 预设测试数据，一键创建测试账户")

	fmt.Printf("\n这些函数都会发送相同的 POST 请求到: /open-api/v3/accounts/register\n")
	fmt.Printf("请求体格式与你的 curl 命令完全一致:\n")
	fmt.Printf(`{
  "phoneCountryCode": "86",
  "phoneNumber": "15900000031", 
  "email": "15900000031@qq.com",
  "name": "golang_test"
}`)
}