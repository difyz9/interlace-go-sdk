package main

import (
	"context"
	"fmt"
	"log"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	fmt.Println("=== CDD (Customer Due Diligence) 详情功能演示 ===")

	clientID := "your-client-id-here"
	
	// 快速设置客户端
	client, _, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Fatalf("客户端设置失败: %v", err)
	}

	ctx := context.Background()

	// 使用测试账户 ID（实际使用中应该是真实的账户 ID）
	accountID := "your-account-id-here"

	// 1. 获取完整的 CDD 详情
	fmt.Println("\n1. 获取 CDD 详情")
	cddDetail, err := client.KYC.GetCDDDetail(ctx, accountID)
	if err != nil {
		log.Printf("获取 CDD 详情失败: %v", err)
		
		// 如果 API 调用失败，我们继续演示其他功能的代码结构
		fmt.Println("\n演示 CDD 数据结构和功能...")
		demonstrateCDDStructure()
		return
	}

	fmt.Printf("✓ CDD 详情获取成功\n")
	fmt.Printf("   账户 ID: %s\n", cddDetail.AccountID)
	fmt.Printf("   总体状态: %s\n", cddDetail.OverallStatus)
	fmt.Printf("   最后更新: %s\n", cddDetail.LastUpdated)

	// 2. 显示 KYC 验证详情
	if cddDetail.KYCVerification != nil {
		fmt.Println("\n2. KYC 验证详情")
		kyc := cddDetail.KYCVerification
		fmt.Printf("   申请 ID: %s\n", kyc.ApplicationID)
		fmt.Printf("   状态: %s\n", kyc.Status)
		fmt.Printf("   提交时间: %s\n", kyc.SubmittedTime)
		if kyc.ReviewedTime != "" {
			fmt.Printf("   审核时间: %s\n", kyc.ReviewedTime)
		}
		if kyc.RejectionReason != "" {
			fmt.Printf("   拒绝原因: %s\n", kyc.RejectionReason)
		}

		// 显示个人信息
		if kyc.PersonalInfo != nil {
			fmt.Println("\n   个人信息:")
			fmt.Printf("     姓名: %s %s\n", kyc.PersonalInfo.FirstName, kyc.PersonalInfo.LastName)
			fmt.Printf("     出生日期: %s\n", kyc.PersonalInfo.DateOfBirth)
			fmt.Printf("     国籍: %s\n", kyc.PersonalInfo.Nationality)
			fmt.Printf("     地址: %s, %s\n", kyc.PersonalInfo.Address, kyc.PersonalInfo.City)
		}

		// 显示文档信息
		if kyc.DocumentInfo != nil {
			fmt.Println("\n   文档验证:")
			fmt.Printf("     ID 类型: %s\n", kyc.DocumentInfo.IDType)
			fmt.Printf("     ID 号码: %s\n", kyc.DocumentInfo.IDNumber)
			fmt.Printf("     正面照片状态: %s\n", kyc.DocumentInfo.IDFrontImageStatus)
			fmt.Printf("     自拍照片状态: %s\n", kyc.DocumentInfo.SelfieImageStatus)
			fmt.Printf("     文档匹配: %v\n", kyc.DocumentInfo.DocumentMatch)
			fmt.Printf("     人脸匹配: %v\n", kyc.DocumentInfo.FaceMatch)
		}

		// 显示验证检查
		if kyc.VerificationChecks != nil {
			fmt.Println("\n   验证检查:")
			displayCheckResult("身份验证", kyc.VerificationChecks.IdentityVerification)
			displayCheckResult("文档验证", kyc.VerificationChecks.DocumentVerification)
			displayCheckResult("生物识别验证", kyc.VerificationChecks.BiometricVerification)
			displayCheckResult("地址验证", kyc.VerificationChecks.AddressVerification)
			displayCheckResult("观察名单筛查", kyc.VerificationChecks.WatchlistScreening)
			displayCheckResult("制裁筛查", kyc.VerificationChecks.SanctionsScreening)
			displayCheckResult("PEP 筛查", kyc.VerificationChecks.PEPScreening)
		}
	}

	// 3. 显示 KYB 验证详情
	if cddDetail.KYBVerification != nil {
		fmt.Println("\n3. KYB 验证详情")
		kyb := cddDetail.KYBVerification
		fmt.Printf("   申请 ID: %s\n", kyb.ApplicationID)
		fmt.Printf("   状态: %s\n", kyb.Status)
		fmt.Printf("   提交时间: %s\n", kyb.SubmittedTime)

		// 显示企业信息
		if kyb.BusinessInfo != nil {
			fmt.Println("\n   企业信息:")
			fmt.Printf("     公司名称: %s\n", kyb.BusinessInfo.CompanyName)
			fmt.Printf("     业务类型: %s\n", kyb.BusinessInfo.BusinessType)
			fmt.Printf("     注册号: %s\n", kyb.BusinessInfo.RegistrationNumber)
			fmt.Printf("     注册国家: %s\n", kyb.BusinessInfo.RegistrationCountry)
			fmt.Printf("     行业: %s\n", kyb.BusinessInfo.Industry)
		}

		// 显示合规检查
		if kyb.ComplianceChecks != nil {
			fmt.Println("\n   合规检查:")
			displayCheckResult("企业注册", kyb.ComplianceChecks.BusinessRegistration)
			displayCheckResult("董事筛查", kyb.ComplianceChecks.DirectorsScreening)
			displayCheckResult("股东筛查", kyb.ComplianceChecks.ShareholdersScreening)
			displayCheckResult("UBO 验证", kyb.ComplianceChecks.UBOVerification)
			displayCheckResult("许可证验证", kyb.ComplianceChecks.LicenseVerification)
		}
	}

	// 4. 显示风险评估
	displayRiskAssessment(cddDetail)

	// 5. 演示便捷方法
	demonstrateConvenienceMethods(client, ctx, accountID)
}

func displayCheckResult(name string, check *interlace.CheckResult) {
	if check != nil {
		status := check.Status
		emoji := getStatusEmoji(status)
		fmt.Printf("     %s %s: %s", emoji, name, status)
		if check.Details != "" {
			fmt.Printf(" (%s)", check.Details)
		}
		if check.Score != nil {
			fmt.Printf(" [评分: %d]", *check.Score)
		}
		fmt.Printf("\n")
	}
}

func getStatusEmoji(status string) string {
	switch status {
	case interlace.CheckStatusPass:
		return "✅"
	case interlace.CheckStatusFail:
		return "❌"
	case interlace.CheckStatusWarning:
		return "⚠️"
	case interlace.CheckStatusPending:
		return "⏳"
	default:
		return "❓"
	}
}

func displayRiskAssessment(cdd *interlace.CDDDetailData) {
	fmt.Println("\n4. 风险评估")
	
	var riskAssessment *interlace.RiskAssessment
	if cdd.KYCVerification != nil && cdd.KYCVerification.RiskAssessment != nil {
		riskAssessment = cdd.KYCVerification.RiskAssessment
		fmt.Println("   来源: KYC 验证")
	} else if cdd.KYBVerification != nil && cdd.KYBVerification.RiskAssessment != nil {
		riskAssessment = cdd.KYBVerification.RiskAssessment
		fmt.Println("   来源: KYB 验证")
	}

	if riskAssessment != nil {
		riskEmoji := getRiskEmoji(riskAssessment.RiskLevel)
		fmt.Printf("   %s 风险等级: %s\n", riskEmoji, riskAssessment.RiskLevel)
		fmt.Printf("   风险评分: %d\n", riskAssessment.RiskScore)
		fmt.Printf("   最后更新: %s\n", riskAssessment.LastUpdated)
		
		if len(riskAssessment.Factors) > 0 {
			fmt.Println("   风险因素:")
			for i, factor := range riskAssessment.Factors {
				fmt.Printf("     %d. %s\n", i+1, factor)
			}
		}
	} else {
		fmt.Println("   暂无风险评估数据")
	}
}

func getRiskEmoji(riskLevel string) string {
	switch riskLevel {
	case interlace.RiskLevelLow:
		return "🟢"
	case interlace.RiskLevelMedium:
		return "🟡"
	case interlace.RiskLevelHigh:
		return "🔴"
	default:
		return "❓"
	}
}

func demonstrateConvenienceMethods(client *interlace.Client, ctx context.Context, accountID string) {
	fmt.Println("\n5. 便捷方法演示")

	// 检查风险等级
	isHighRisk, err := client.KYC.IsHighRisk(ctx, accountID)
	if err != nil {
		log.Printf("检查风险等级失败: %v", err)
	} else {
		fmt.Printf("   高风险账户: %v\n", isHighRisk)
	}

	// 检查所有验证是否通过
	allPassed, failedChecks, err := client.KYC.HasPassedAllChecks(ctx, accountID)
	if err != nil {
		log.Printf("检查验证状态失败: %v", err)
	} else {
		fmt.Printf("   所有检查通过: %v\n", allPassed)
		if !allPassed && len(failedChecks) > 0 {
			fmt.Println("   未通过的检查:")
			for _, check := range failedChecks {
				fmt.Printf("     - %s\n", check)
			}
		}
	}

	fmt.Println("\n=== CDD 功能总结 ===")
	fmt.Println("✅ GetCDDDetail() - 获取完整 CDD 详情")
	fmt.Println("✅ GetKYCVerificationDetail() - 获取 KYC 验证详情")
	fmt.Println("✅ GetKYBVerificationDetail() - 获取 KYB 验证详情")
	fmt.Println("✅ GetRiskAssessment() - 获取风险评估")
	fmt.Println("✅ IsHighRisk() - 检查是否高风险")
	fmt.Println("✅ GetVerificationChecks() - 获取验证检查")
	fmt.Println("✅ GetComplianceChecks() - 获取合规检查")
	fmt.Println("✅ HasPassedAllChecks() - 检查所有验证是否通过")
}

func demonstrateCDDStructure() {
	fmt.Println("\n=== CDD 数据结构说明 ===")
	fmt.Println("CDD (Customer Due Diligence) 包含:")
	fmt.Println("• KYC 验证 (个人客户)")
	fmt.Println("  - 个人信息验证")
	fmt.Println("  - 身份文档验证")
	fmt.Println("  - 生物识别验证")
	fmt.Println("  - 反洗钱筛查")
	fmt.Println("• KYB 验证 (企业客户)")
	fmt.Println("  - 企业注册验证")
	fmt.Println("  - 董事/股东筛查")
	fmt.Println("  - UBO (实际受益人) 验证")
	fmt.Println("  - 合规检查")
	fmt.Println("• 风险评估")
	fmt.Println("  - 低风险 (LOW)")
	fmt.Println("  - 中风险 (MEDIUM)")
	fmt.Println("  - 高风险 (HIGH)")
}