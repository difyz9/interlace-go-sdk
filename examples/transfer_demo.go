package main

import (
	"context"
	"fmt"
	"log"
	"time"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	fmt.Println("=== 区块链转账管理 API 演示 ===")

	clientID := "your-client-id-here"
	
	// 快速设置客户端
	client, _, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Fatalf("客户端设置失败: %v", err)
	}

	ctx := context.Background()

	// 使用测试钱包ID (实际使用时应该从钱包列表获取)
	walletID := "test-wallet-id"
	fmt.Printf("\n使用测试钱包: %s\n", walletID)

	// 1. 获取转账手续费和额度 (Get Fee and Quota)
	fmt.Println("\n1. 获取转账手续费和额度")
	feeReq := &interlace.FeeAndQuotaRequest{
		WalletID:  walletID,
		Currency:  "USDT",
		Chain:     "TRC20",
		Amount:    "100",
		ToAddress: "TRX7NqM77kjvzP4CLfRhcPBnPeHCkXqbfp",
	}

	feeQuota, err := client.Transfer.GetFeeAndQuota(ctx, feeReq)
	if err != nil {
		log.Printf("获取手续费和额度失败: %v", err)
		fmt.Println("\n注意: 需要有效的钱包ID才能使用转账功能")
	} else {
		fmt.Printf("✓ 手续费信息:\n")
		fmt.Printf("   手续费: %s %s\n", feeQuota.Fee, feeQuota.FeeCurrency)
		fmt.Printf("   最小金额: %s\n", feeQuota.MinAmount)
		fmt.Printf("   最大金额: %s\n", feeQuota.MaxAmount)
		fmt.Printf("   可用额度: %s\n", feeQuota.AvailableQuota)
		if feeQuota.EstimatedArrival != "" {
			fmt.Printf("   预计到账: %s\n", feeQuota.EstimatedArrival)
		}
	}

	// 2. 创建区块链转账 (Create Transfer)
	fmt.Println("\n2. 创建区块链转账")
	transferReq := &interlace.CreateTransferRequest{
		WalletID:       walletID,
		Currency:       "USDT",
		Chain:          "TRC20",
		Amount:         "50",
		ToAddress:      "TRX7NqM77kjvzP4CLfRhcPBnPeHCkXqbfp",
		IdempotencyKey: fmt.Sprintf("transfer-%d", time.Now().Unix()),
	}

	transfer, err := client.Transfer.CreateTransfer(ctx, transferReq)
	if err != nil {
		log.Printf("创建转账失败: %v", err)
		fmt.Println("\n提示: 演示创建转账的请求结构")
		fmt.Printf("   钱包ID: %s\n", transferReq.WalletID)
		fmt.Printf("   币种: %s\n", transferReq.Currency)
		fmt.Printf("   链: %s\n", transferReq.Chain)
		fmt.Printf("   金额: %s\n", transferReq.Amount)
		fmt.Printf("   收款地址: %s\n", transferReq.ToAddress)
		fmt.Printf("   幂等性键: %s\n", transferReq.IdempotencyKey)
	} else {
		fmt.Printf("✓ 转账创建成功\n")
		fmt.Printf("   转账ID: %s\n", transfer.ID)
		fmt.Printf("   状态: %s\n", transfer.Status)
		fmt.Printf("   金额: %s %s\n", transfer.Amount, transfer.Currency)
		fmt.Printf("   手续费: %s\n", transfer.Fee)
		fmt.Printf("   收款地址: %s\n", transfer.ToAddress)
		if transfer.TxHash != "" {
			fmt.Printf("   交易哈希: %s\n", transfer.TxHash)
		}
	}

	// 3. 获取转账列表 (List Transfers)
	fmt.Println("\n3. 获取转账列表")
	listOptions := &interlace.TransferListOptions{
		WalletID: walletID,
		Limit:    10,
		Page:     1,
	}

	transferList, err := client.Transfer.ListTransfers(ctx, listOptions)
	if err != nil {
		log.Printf("获取转账列表失败: %v", err)
	} else {
		fmt.Printf("✓ 获取到 %d 笔转账 (总共 %d 笔)\n", len(transferList.Transfers), transferList.TotalCount)
		
		for i, t := range transferList.Transfers {
			fmt.Printf("   转账 %d:\n", i+1)
			fmt.Printf("     ID: %s\n", t.ID)
			fmt.Printf("     币种/链: %s/%s\n", t.Currency, t.Chain)
			fmt.Printf("     金额: %s\n", t.Amount)
			fmt.Printf("     状态: %s\n", t.Status)
			fmt.Printf("     收款地址: %s\n", t.ToAddress)
			if t.TxHash != "" {
				fmt.Printf("     交易哈希: %s\n", t.TxHash)
			}
			fmt.Printf("     确认数: %d/%d\n", t.Confirmations, t.RequiredConfirms)
			fmt.Printf("     创建时间: %s\n", t.CreatedAt)
			
			// 演示获取单个转账详情
			if i == 0 {
				fmt.Printf("\n   获取转账详情 (ID: %s):\n", t.ID)
				detail, err := client.Transfer.GetTransfer(ctx, t.ID)
				if err != nil {
					log.Printf("     获取转账详情失败: %v", err)
				} else {
					fmt.Printf("     ✓ 状态: %s\n", detail.Status)
					fmt.Printf("     ✓ 更新时间: %s\n", detail.UpdatedAt)
				}
				
				// 演示获取KYT信息
				fmt.Printf("\n   获取转账KYT信息:\n")
				kyt, err := client.Transfer.GetTransferKYT(ctx, t.ID)
				if err != nil {
					log.Printf("     获取KYT信息失败: %v", err)
				} else {
					fmt.Printf("     ✓ 风险评分: %d\n", kyt.RiskScore)
					fmt.Printf("     ✓ 风险等级: %s\n", kyt.RiskLevel)
					fmt.Printf("     ✓ 检查时间: %s\n", kyt.CheckedAt)
					if len(kyt.Alerts) > 0 {
						fmt.Printf("     ⚠️  警告:\n")
						for _, alert := range kyt.Alerts {
							fmt.Printf("        - [%s] %s: %s\n", alert.Severity, alert.Type, alert.Description)
						}
					}
				}
			}
		}
	}

	// 4. 演示按状态筛选转账
	fmt.Println("\n4. 按状态筛选转账")
	statusOptions := &interlace.TransferListOptions{
		WalletID: walletID,
		Status:   "PENDING",
		Limit:    5,
		Page:     1,
	}

	pendingTransfers, err := client.Transfer.ListTransfers(ctx, statusOptions)
	if err != nil {
		log.Printf("获取待处理转账失败: %v", err)
	} else {
		fmt.Printf("✓ 待处理转账: %d 笔\n", len(pendingTransfers.Transfers))
	}

	// 5. 演示按币种和链筛选
	fmt.Println("\n5. 按币种和链筛选转账")
	currencyOptions := &interlace.TransferListOptions{
		WalletID: walletID,
		Currency: "USDT",
		Chain:    "TRC20",
		Limit:    5,
		Page:     1,
	}

	usdtTransfers, err := client.Transfer.ListTransfers(ctx, currencyOptions)
	if err != nil {
		log.Printf("获取USDT转账失败: %v", err)
	} else {
		fmt.Printf("✓ USDT/TRC20 转账: %d 笔\n", len(usdtTransfers.Transfers))
	}

	fmt.Println("\n=== API 功能总结 ===")
	fmt.Println("✅ CreateTransfer() - 创建区块链转账")
	fmt.Println("   - 从钱包转账到区块链地址")
	fmt.Println("   - 需要指定钱包ID、币种、链、金额、收款地址")
	fmt.Println("   - 需要提供幂等性键防止重复提交")
	fmt.Println("")
	fmt.Println("✅ ListTransfers() - 获取转账列表")
	fmt.Println("   - 按钱包ID、币种、链、状态筛选")
	fmt.Println("   - 支持分页 (limit, page)")
	fmt.Println("")
	fmt.Println("✅ GetTransfer() - 获取转账详情")
	fmt.Println("   - 通过转账ID获取完整信息")
	fmt.Println("   - 包含状态、确认数、交易哈希等")
	fmt.Println("")
	fmt.Println("✅ GetTransferKYT() - 获取转账KYT信息")
	fmt.Println("   - Know Your Transaction 风险评估")
	fmt.Println("   - 包含风险评分、风险等级、警告信息")
	fmt.Println("")
	fmt.Println("✅ GetFeeAndQuota() - 获取手续费和额度")
	fmt.Println("   - 预览转账手续费")
	fmt.Println("   - 查看跨链额度限制")
	fmt.Println("   - 获取最小/最大转账金额")
}