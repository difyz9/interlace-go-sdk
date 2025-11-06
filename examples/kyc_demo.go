package main

import (
	"context"
	"fmt"
	"log"
	"time"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	fmt.Println("=== KYC 认证功能演示 ===")

	clientID := "your-client-id-here"
	
	// 快速设置客户端
	client, _, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Fatalf("客户端设置失败: %v", err)
	}

	ctx := context.Background()

	// 首先注册一个测试账户
	fmt.Println("\n1. 注册测试账户")
	account, err := client.Account.RegisterGolangTest(ctx)
	if err != nil {
		log.Printf("账户注册失败: %v", err)
		return
	}
	fmt.Printf("✓ 账户注册成功 - ID: %s\n", account.ID)

	// 2. 模拟文件上传（实际场景中需要上传真实文件）
	fmt.Println("\n2. 准备 KYC 文档")
	fmt.Println("注意：在实际使用中，您需要：")
	fmt.Println("   a) 上传身份证正面照片")
	fmt.Println("   b) 上传身份证背面照片（可选）")
	fmt.Println("   c) 上传自拍照片")
	
	// 模拟文件 ID（实际中通过 client.File.UploadFile 获得）
	mockIDFrontFileID := "mock-id-front-file-id"
	mockIDBackFileID := "mock-id-back-file-id"
	mockSelfieFileID := "mock-selfie-file-id"

	// 3. 使用 Builder 模式构建 KYC 请求
	fmt.Println("\n3. 构建 KYC 请求")
	kycBuilder := interlace.NewKYCBuilder().
		SetPersonalInfo("张", "三", "1990-01-15", "M").
		SetNationality("CN", "CN").
		SetAddress("北京市朝阳区某某街道123号", "北京", "100000", "CN").
		SetIDInfo(interlace.IDTypeCNRIC, "110101199001011234").
		SetIDExpiryDate("2030-01-15").
		SetOccupationInfo("01", "Employment"). // 01 是职业代码，需要查看文档
		SetAnnualIncome("100000").
		SetAccountPurpose("Personal Banking").
		SetExpectedTxnVolume("Medium").
		SetDocumentFiles(mockIDFrontFileID, mockSelfieFileID).
		SetIDBackFile(mockIDBackFileID)

	// 验证请求
	if err := kycBuilder.Validate(); err != nil {
		log.Printf("KYC 请求验证失败: %v", err)
		return
	}

	kycRequest := kycBuilder.Build()
	fmt.Printf("✓ KYC 请求构建完成\n")

	// 4. 提交 KYC 申请
	fmt.Println("\n4. 提交 KYC 申请")
	kycSubmitData, err := client.KYC.SubmitKYC(ctx, account.ID, kycRequest)
	if err != nil {
		log.Printf("KYC 提交失败: %v", err)
		
		// 演示如何继续使用已有的 KYC 申请ID（如果有的话）
		fmt.Println("\n继续演示其他 KYC 功能...")
		mockKYCAppID := "mock-kyc-application-id"
		demonstrateKYCStatus(client, ctx, account.ID, mockKYCAppID)
		return
	}

	fmt.Printf("✓ KYC 申请提交成功\n")
	fmt.Printf("   申请 ID: %s\n", kycSubmitData.KYCApplicationID)
	fmt.Printf("   状态: %s\n", kycSubmitData.Status)
	fmt.Printf("   提交时间: %s\n", kycSubmitData.SubmittedTime)

	// 5. 查询 KYC 状态
	demonstrateKYCStatus(client, ctx, account.ID, kycSubmitData.KYCApplicationID)

	// 6. 演示便捷方法
	demonstrateKYCConvenienceMethods(client, ctx, account.ID)
}

func demonstrateKYCStatus(client *interlace.Client, ctx context.Context, accountID, kycAppID string) {
	fmt.Println("\n5. 查询 KYC 状态")
	
	kycStatus, err := client.KYC.GetKYCStatus(ctx, accountID)
	if err != nil {
		log.Printf("查询 KYC 状态失败: %v", err)
		return
	}

	fmt.Printf("✓ KYC 状态查询成功\n")
	fmt.Printf("   账户 ID: %s\n", kycStatus.AccountID)
	fmt.Printf("   申请 ID: %s\n", kycStatus.KYCApplicationID)
	fmt.Printf("   状态: %s\n", kycStatus.Status)
	
	if kycStatus.SubmittedTime != "" {
		fmt.Printf("   提交时间: %s\n", kycStatus.SubmittedTime)
	}
	if kycStatus.ReviewedTime != "" {
		fmt.Printf("   审核时间: %s\n", kycStatus.ReviewedTime)
	}
	if kycStatus.RejectionReason != "" {
		fmt.Printf("   拒绝原因: %s\n", kycStatus.RejectionReason)
	}
	if kycStatus.ExpiryTime != "" {
		fmt.Printf("   过期时间: %s\n", kycStatus.ExpiryTime)
	}
}

func demonstrateKYCConvenienceMethods(client *interlace.Client, ctx context.Context, accountID string) {
	fmt.Println("\n6. KYC 便捷方法演示")

	// 检查各种状态
	isPending, err := client.KYC.IsKYCPending(ctx, accountID)
	if err != nil {
		log.Printf("检查待审核状态失败: %v", err)
	} else {
		fmt.Printf("   KYC 待审核: %v\n", isPending)
	}

	isApproved, err := client.KYC.IsKYCApproved(ctx, accountID)
	if err != nil {
		log.Printf("检查已批准状态失败: %v", err)
	} else {
		fmt.Printf("   KYC 已批准: %v\n", isApproved)
	}

	isRejected, err := client.KYC.IsKYCRejected(ctx, accountID)
	if err != nil {
		log.Printf("检查已拒绝状态失败: %v", err)
	} else {
		fmt.Printf("   KYC 已拒绝: %v\n", isRejected)
	}

	// 演示轮询等待（仅在测试环境中使用）
	fmt.Println("\n7. 等待 KYC 审批演示（实际环境请谨慎使用）")
	fmt.Println("   注意：这将尝试轮询 KYC 状态，最多 3 次")
	
	finalStatus, err := client.KYC.WaitForKYCApproval(ctx, accountID, 3)
	if err != nil {
		log.Printf("等待 KYC 审批失败: %v", err)
	} else {
		fmt.Printf("✓ KYC 最终状态: %s\n", finalStatus.Status)
	}

	fmt.Println("\n=== KYC 功能总结 ===")
	fmt.Println("✅ SubmitKYC() - 提交 KYC 申请")
	fmt.Println("✅ GetKYCStatus() - 查询 KYC 状态")
	fmt.Println("✅ IsKYCApproved() - 检查是否已批准")
	fmt.Println("✅ IsKYCPending() - 检查是否待审核")
	fmt.Println("✅ IsKYCRejected() - 检查是否已拒绝")
	fmt.Println("✅ WaitForKYCApproval() - 轮询等待审批")
	fmt.Println("✅ KYCBuilder - 流式构建器模式")

	fmt.Println("\n=== 支持的 ID 类型 ===")
	fmt.Printf("• %s - 中国居民身份证\n", interlace.IDTypeCNRIC)
	fmt.Printf("• %s - 香港身份证\n", interlace.IDTypeHKHKID)
	fmt.Printf("• %s - 护照\n", interlace.IDTypePassport)
	fmt.Printf("• %s - 驾驶证\n", interlace.IDTypeDLN)
	fmt.Printf("• %s - 政府颁发身份证\n", interlace.IDTypeGovernmentID)

	fmt.Println("\n=== KYC 状态说明 ===")
	fmt.Printf("• %s - 待审核\n", interlace.KYCStatusPending)
	fmt.Printf("• %s - 已批准\n", interlace.KYCStatusApproved)
	fmt.Printf("• %s - 已拒绝\n", interlace.KYCStatusRejected)
	fmt.Printf("• %s - 已过期\n", interlace.KYCStatusExpired)
}