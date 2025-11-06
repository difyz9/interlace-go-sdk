package main

import (
	"context"
	"fmt"
	"log"
	"time"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	fmt.Println("=== 支付和退款管理 API 演示 ===")

	clientID := "your-client-id-here"
	
	// 快速设置客户端
	client, _, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Fatalf("客户端设置失败: %v", err)
	}

	ctx := context.Background()

	// 1. 创建支付订单 (Create Payment)
	fmt.Println("\n1. 创建支付订单")
	merchantTradeNo := fmt.Sprintf("ORDER-%d", time.Now().Unix())
	
	createPaymentReq := &interlace.CreatePaymentRequest{
		MerchantTradeNo: merchantTradeNo,
		Amount:          "100.00",
		Currency:        "USD",
		Country:         "US",
		Description:     "Test payment order",
		PaymentMethod:   "card",
		ReturnURL:       "https://example.com/return",
		NotifyURL:       "https://example.com/notify",
	}

	payment, err := client.Payment.CreatePayment(ctx, createPaymentReq)
	if err != nil {
		log.Printf("创建支付订单失败: %v", err)
		fmt.Println("\n提示: 演示创建支付订单的请求结构")
		fmt.Printf("   商户订单号: %s\n", createPaymentReq.MerchantTradeNo)
		fmt.Printf("   金额: %s %s\n", createPaymentReq.Amount, createPaymentReq.Currency)
		fmt.Printf("   国家: %s\n", createPaymentReq.Country)
		fmt.Printf("   描述: %s\n", createPaymentReq.Description)
		fmt.Printf("   支付方式: %s\n", createPaymentReq.PaymentMethod)
	} else {
		fmt.Printf("✓ 支付订单创建成功\n")
		fmt.Printf("   订单ID: %s\n", payment.ID)
		fmt.Printf("   商户订单号: %s\n", payment.MerchantTradeNo)
		fmt.Printf("   金额: %s %s\n", payment.Amount, payment.Currency)
		fmt.Printf("   状态: %s\n", payment.Status)
		if payment.CheckoutURL != "" {
			fmt.Printf("   支付链接: %s\n", payment.CheckoutURL)
		}
		fmt.Printf("   创建时间: %s\n", payment.CreatedAt)

		// 2. 查询支付订单 (Query Payment)
		fmt.Printf("\n2. 查询支付订单 (订单号: %s)\n", payment.ID)
		queriedPayment, err := client.Payment.QueryPayment(ctx, payment.ID)
		if err != nil {
			log.Printf("查询支付订单失败: %v", err)
		} else {
			fmt.Printf("✓ 支付订单详情:\n")
			fmt.Printf("   订单ID: %s\n", queriedPayment.ID)
			fmt.Printf("   状态: %s\n", queriedPayment.Status)
			fmt.Printf("   更新时间: %s\n", queriedPayment.UpdatedAt)
			if queriedPayment.PaidAt != "" {
				fmt.Printf("   支付时间: %s\n", queriedPayment.PaidAt)
			}
		}

		// 3. 取消支付订单 (Cancel Payment)
		fmt.Printf("\n3. 取消支付订单\n")
		cancelReq := &interlace.CancelPaymentRequest{
			OrderNo: payment.ID,
		}
		
		cancelledPayment, err := client.Payment.CancelPayment(ctx, cancelReq)
		if err != nil {
			log.Printf("取消支付订单失败: %v", err)
		} else {
			fmt.Printf("✓ 支付订单已取消\n")
			fmt.Printf("   订单ID: %s\n", cancelledPayment.ID)
			fmt.Printf("   状态: %s\n", cancelledPayment.Status)
			fmt.Printf("   更新时间: %s\n", cancelledPayment.UpdatedAt)
		}
	}

	// 4. 创建退款 (Create Refund)
	fmt.Println("\n4. 创建退款订单")
	
	// 假设我们有一个已支付的订单
	sourceMerchantTradeNo := merchantTradeNo // 使用上面创建的支付订单号
	refundMerchantTradeNo := fmt.Sprintf("REFUND-%d", time.Now().Unix())
	
	createRefundReq := &interlace.CreateRefundRequest{
		SourceMerchantTradeNo: sourceMerchantTradeNo,
		MerchantTradeNo:       refundMerchantTradeNo,
		Amount:                "50.00",
		Reason:                "Customer request",
	}

	refund, err := client.Payment.CreateRefund(ctx, createRefundReq)
	if err != nil {
		log.Printf("创建退款失败: %v", err)
		fmt.Println("\n提示: 演示创建退款的请求结构")
		fmt.Printf("   原支付订单号: %s\n", createRefundReq.SourceMerchantTradeNo)
		fmt.Printf("   商户退款单号: %s\n", createRefundReq.MerchantTradeNo)
		fmt.Printf("   退款金额: %s\n", createRefundReq.Amount)
		fmt.Printf("   退款原因: %s\n", createRefundReq.Reason)
	} else {
		fmt.Printf("✓ 退款订单创建成功\n")
		fmt.Printf("   退款ID: %s\n", refund.ID)
		fmt.Printf("   原支付订单ID: %s\n", refund.PaymentID)
		fmt.Printf("   商户退款单号: %s\n", refund.MerchantTradeNo)
		fmt.Printf("   退款金额: %s %s\n", refund.Amount, refund.Currency)
		fmt.Printf("   状态: %s\n", refund.Status)
		fmt.Printf("   创建时间: %s\n", refund.CreatedAt)

		// 5. 查询退款订单 (Query Refund)
		fmt.Printf("\n5. 查询退款订单 (退款号: %s)\n", refund.ID)
		queriedRefund, err := client.Payment.QueryRefund(ctx, refund.ID)
		if err != nil {
			log.Printf("查询退款订单失败: %v", err)
		} else {
			fmt.Printf("✓ 退款订单详情:\n")
			fmt.Printf("   退款ID: %s\n", queriedRefund.ID)
			fmt.Printf("   状态: %s\n", queriedRefund.Status)
			fmt.Printf("   更新时间: %s\n", queriedRefund.UpdatedAt)
			if queriedRefund.ProcessedAt != "" {
				fmt.Printf("   处理时间: %s\n", queriedRefund.ProcessedAt)
			}
		}
	}

	// 6. 批量搜索订单 (Search)
	fmt.Println("\n6. 批量搜索支付和退款订单")
	
	orderNos := []string{merchantTradeNo, refundMerchantTradeNo}
	
	searchResult, err := client.Payment.Search(ctx, orderNos)
	if err != nil {
		log.Printf("搜索订单失败: %v", err)
		fmt.Println("\n提示: Search API 可以批量查询多个订单")
		fmt.Printf("   查询订单号列表: %v\n", orderNos)
	} else {
		fmt.Printf("✓ 搜索结果:\n")
		
		if len(searchResult.Payments) > 0 {
			fmt.Printf("   支付订单: %d 笔\n", len(searchResult.Payments))
			for i, p := range searchResult.Payments {
				fmt.Printf("     %d. 订单ID: %s, 状态: %s, 金额: %s %s\n", 
					i+1, p.ID, p.Status, p.Amount, p.Currency)
			}
		} else {
			fmt.Println("   支付订单: 0 笔")
		}
		
		if len(searchResult.Refunds) > 0 {
			fmt.Printf("   退款订单: %d 笔\n", len(searchResult.Refunds))
			for i, r := range searchResult.Refunds {
				fmt.Printf("     %d. 退款ID: %s, 状态: %s, 金额: %s %s\n", 
					i+1, r.ID, r.Status, r.Amount, r.Currency)
			}
		} else {
			fmt.Println("   退款订单: 0 笔")
		}
	}

	fmt.Println("\n=== API 功能总结 ===")
	fmt.Println("✅ CreatePayment() - 创建支付订单")
	fmt.Println("   - 支持多种支付方式 (卡支付、APM等)")
	fmt.Println("   - 可配置回调URL和通知URL")
	fmt.Println("   - 返回支付链接用于收款")
	fmt.Println("")
	fmt.Println("✅ CancelPayment() - 取消支付订单")
	fmt.Println("   - 用于未完成的支付订单")
	fmt.Println("   - 取消后支付链接将失效")
	fmt.Println("")
	fmt.Println("✅ CreateRefund() - 创建退款")
	fmt.Println("   - 仅对已支付订单有效")
	fmt.Println("   - 支持部分退款")
	fmt.Println("   - 同一笔订单可多次退款，总额不超过原金额")
	fmt.Println("")
	fmt.Println("✅ QueryPayment() - 查询支付订单")
	fmt.Println("   - 支持系统订单号或商户订单号查询")
	fmt.Println("   - 返回订单完整状态信息")
	fmt.Println("")
	fmt.Println("✅ QueryRefund() - 查询退款订单")
	fmt.Println("   - 支持系统退款号或商户退款单号查询")
	fmt.Println("   - 返回退款处理状态")
	fmt.Println("")
	fmt.Println("✅ Search() - 批量搜索订单")
	fmt.Println("   - 同时搜索支付和退款订单")
	fmt.Println("   - 支持多个订单号批量查询")
}