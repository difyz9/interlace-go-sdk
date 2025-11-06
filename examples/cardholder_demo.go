package main

import (
	"context"
	"fmt"
	"log"
	"time"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	fmt.Println("=== 持卡人管理 API 演示 ===")

	clientID := "your-client-id-here"
	
	// 快速设置客户端
	client, _, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Fatalf("客户端设置失败: %v", err)
	}

	ctx := context.Background()
	
	// 测试账户ID（需要使用实际的账户ID）
	accountID := "your-account-id-here"

	// 1. 创建持卡人 (Create Cardholder)
	fmt.Println("\n1. 创建持卡人")
	
	createReq := &interlace.CreateCardholderRequest{
		AccountID:        accountID,
		FirstName:        "John",
		LastName:         "Doe",
		Email:            fmt.Sprintf("john.doe.%d@example.com", time.Now().Unix()),
		PhoneNumber:      "15900000001",
		PhoneCountryCode: "86",
		DateOfBirth:      "1990-01-01",
		Nationality:      "CN",
		Gender:          "M",
		Address: &interlace.CardholderAddress{
			Country:    "CN",
			State:      "Beijing",
			City:       "Beijing",
			PostalCode: "100000",
			Line1:      "123 Main Street",
		},
		IdentityDocument: &interlace.IdentityDocument{
			Type:           "CN-RIC",
			Number:         "110101199001011234",
			IssuingCountry: "CN",
		},
		IdempotencyKey: fmt.Sprintf("cardholder_%d", time.Now().Unix()),
	}

	cardholder, err := client.Cardholder.CreateCardholder(ctx, createReq)
	if err != nil {
		log.Printf("创建持卡人失败: %v", err)
		fmt.Println("\n提示: 演示创建持卡人的请求结构")
		fmt.Printf("   账户ID: %s\n", createReq.AccountID)
		fmt.Printf("   姓名: %s %s\n", createReq.FirstName, createReq.LastName)
		fmt.Printf("   邮箱: %s\n", createReq.Email)
		fmt.Printf("   电话: +%s %s\n", createReq.PhoneCountryCode, createReq.PhoneNumber)
		fmt.Printf("   出生日期: %s\n", createReq.DateOfBirth)
		fmt.Printf("   国籍: %s\n", createReq.Nationality)
	} else {
		fmt.Printf("✓ 持卡人创建成功\n")
		fmt.Printf("   持卡人ID: %s\n", cardholder.ID)
		fmt.Printf("   账户ID: %s\n", cardholder.AccountID)
		fmt.Printf("   姓名: %s %s\n", cardholder.FirstName, cardholder.LastName)
		fmt.Printf("   邮箱: %s\n", cardholder.Email)
		fmt.Printf("   电话: +%s %s\n", cardholder.PhoneCountryCode, cardholder.PhoneNumber)
		fmt.Printf("   状态: %s\n", cardholder.Status)
		fmt.Printf("   创建时间: %s\n", cardholder.CreatedAt)

		// 2. 获取持卡人详情 (Get Cardholder)
		fmt.Printf("\n2. 获取持卡人详情 (ID: %s)\n", cardholder.ID)
		retrievedCardholder, err := client.Cardholder.GetCardholder(ctx, cardholder.ID)
		if err != nil {
			log.Printf("获取持卡人详情失败: %v", err)
		} else {
			fmt.Printf("✓ 持卡人详情:\n")
			fmt.Printf("   ID: %s\n", retrievedCardholder.ID)
			fmt.Printf("   姓名: %s %s\n", retrievedCardholder.FirstName, retrievedCardholder.LastName)
			fmt.Printf("   邮箱: %s\n", retrievedCardholder.Email)
			fmt.Printf("   状态: %s\n", retrievedCardholder.Status)
			fmt.Printf("   更新时间: %s\n", retrievedCardholder.UpdatedAt)
		}

		// 3. 更新持卡人信息 (Update Cardholder)
		fmt.Printf("\n3. 更新持卡人信息\n")
		updateReq := &interlace.UpdateCardholderRequest{
			Email: fmt.Sprintf("john.doe.updated.%d@example.com", time.Now().Unix()),
		}

		updatedCardholder, err := client.Cardholder.UpdateCardholder(ctx, cardholder.ID, updateReq)
		if err != nil {
			log.Printf("更新持卡人信息失败: %v", err)
		} else {
			fmt.Printf("✓ 持卡人信息已更新\n")
			fmt.Printf("   ID: %s\n", updatedCardholder.ID)
			fmt.Printf("   新邮箱: %s\n", updatedCardholder.Email)
			fmt.Printf("   更新时间: %s\n", updatedCardholder.UpdatedAt)
		}
	}

	// 4. 列出持卡人 (List Cardholders)
	fmt.Println("\n4. 列出持卡人")
	
	listOpts := &interlace.CardholderListOptions{
		AccountID: accountID,
		Page:      1,
		Limit:     10,
	}

	cardholderList, err := client.Cardholder.ListCardholders(ctx, listOpts)
	if err != nil {
		log.Printf("列出持卡人失败: %v", err)
		fmt.Println("\n提示: 演示列出持卡人的请求参数")
		fmt.Printf("   账户ID: %s\n", listOpts.AccountID)
		fmt.Printf("   页码: %d\n", listOpts.Page)
		fmt.Printf("   每页数量: %d\n", listOpts.Limit)
	} else {
		fmt.Printf("✓ 持卡人列表 (总数: %s):\n", cardholderList.Total)
		if len(cardholderList.List) > 0 {
			for i, ch := range cardholderList.List {
				fmt.Printf("   %d. ID: %s, 姓名: %s %s, 状态: %s\n", 
					i+1, ch.ID, ch.FirstName, ch.LastName, ch.Status)
			}
		} else {
			fmt.Println("   暂无持卡人数据")
		}
	}

	fmt.Println("\n=== API 功能总结 ===")
	fmt.Println("✅ CreateCardholder() - 创建持卡人")
	fmt.Println("   - 需要提供完整的个人信息")
	fmt.Println("   - 包括姓名、联系方式、身份证件等")
	fmt.Println("   - 支持地址和身份文档信息")
	fmt.Println("")
	fmt.Println("✅ GetCardholder() - 获取持卡人详情")
	fmt.Println("   - 通过持卡人ID查询")
	fmt.Println("   - 返回完整的持卡人信息")
	fmt.Println("")
	fmt.Println("✅ UpdateCardholder() - 更新持卡人信息")
	fmt.Println("   - 可更新邮箱、电话、地址等信息")
	fmt.Println("   - 部分字段可选更新")
	fmt.Println("")
	fmt.Println("✅ ListCardholders() - 列出持卡人")
	fmt.Println("   - 支持按账户ID筛选")
	fmt.Println("   - 支持分页查询")
	fmt.Println("   - 返回持卡人列表和总数")
}
