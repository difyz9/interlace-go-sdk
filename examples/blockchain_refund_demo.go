package main

import (
	"context"
	"fmt"
	"log"
	"time"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	// Initialize client
	config := interlace.DefaultConfig()
	client := interlace.NewClient(config)

	// Set your access token
	accessToken := "your-access-token-here"
	client.SetAccessToken(accessToken)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Example IDs
	walletID := "your-wallet-id"
	transferID := "your-transfer-id"

	fmt.Println("=== Blockchain Refund Demo ===\n")

	// 1. Get Refund Gas Fee Estimate
	fmt.Println("1. Estimating Gas Fee for Refund...")
	gasFeeRequest := &interlace.GetRefundGasFeeRequest{
		Chain:    "ETHEREUM",
		Currency: "USDT",
		Amount:   100.00,
	}

	gasFee, err := client.BlockchainRefund.GetRefundGasFee(ctx, gasFeeRequest)
	if err != nil {
		log.Printf("Error getting gas fee: %v\n", err)
	} else {
		fmt.Printf("Gas fee estimate:\n")
		fmt.Printf("  Chain: %s\n", gasFee.Chain)
		fmt.Printf("  Currency: %s\n", gasFee.Currency)
		fmt.Printf("  Estimated Gas Fee: %.6f\n", gasFee.EstimatedGasFee)
		fmt.Printf("  Gas Price: %s\n", gasFee.GasPrice)
		fmt.Printf("  Gas Limit: %d\n", gasFee.GasLimit)
		fmt.Printf("  Valid Until: %s\n", gasFee.ValidUntil)
		fmt.Println("\n  ⚠️  Note: Gas fees are estimates and may vary")
	}

	// 2. Create Blockchain Refund
	fmt.Println("\n2. Creating Blockchain Refund...")
	refundRequest := &interlace.CreateBlockchainRefundRequest{
		WalletID:         walletID,
		TransferID:       transferID,
		Chain:            "ETHEREUM",
		Currency:         "USDT",
		Amount:           100.00,
		ToAddress:        "0x1234567890abcdef1234567890abcdef12345678",
		MerchantRefundNo: fmt.Sprintf("REFUND-%d", time.Now().Unix()),
		Reason:           "Customer requested refund",
	}

	refund, err := client.BlockchainRefund.CreateBlockchainRefund(ctx, refundRequest)
	if err != nil {
		log.Printf("Error creating refund: %v\n", err)
	} else {
		fmt.Printf("Refund created successfully!\n")
		fmt.Printf("  Refund ID: %s\n", refund.RefundID)
		fmt.Printf("  Wallet ID: %s\n", refund.WalletID)
		fmt.Printf("  Transfer ID: %s\n", refund.TransferID)
		fmt.Printf("  Chain: %s\n", refund.Chain)
		fmt.Printf("  Currency: %s\n", refund.Currency)
		fmt.Printf("  Amount: %.2f %s\n", refund.Amount, refund.Currency)
		fmt.Printf("  From Address: %s\n", refund.FromAddress)
		fmt.Printf("  To Address: %s\n", refund.ToAddress)
		fmt.Printf("  Gas Fee: %.6f\n", refund.GasFee)
		fmt.Printf("  Status: %s\n", refund.Status)
		fmt.Printf("  Merchant Refund No: %s\n", refund.MerchantRefundNo)
		fmt.Printf("  Reason: %s\n", refund.Reason)
		fmt.Printf("  Created At: %s\n", refund.CreatedAt)
		if refund.TxHash != "" {
			fmt.Printf("  Transaction Hash: %s\n", refund.TxHash)
		}
	}

	// 3. Get Specific Refund Details
	fmt.Println("\n3. Getting Refund Details...")
	if refund != nil {
		refundDetails, err := client.BlockchainRefund.GetBlockchainRefund(ctx, refund.RefundID)
		if err != nil {
			log.Printf("Error getting refund details: %v\n", err)
		} else {
			fmt.Printf("Refund details retrieved:\n")
			fmt.Printf("  Refund ID: %s\n", refundDetails.RefundID)
			fmt.Printf("  Status: %s\n", refundDetails.Status)
			if refundDetails.TxHash != "" {
				fmt.Printf("  Transaction Hash: %s\n", refundDetails.TxHash)
				fmt.Printf("  🔗 View on Explorer: https://etherscan.io/tx/%s\n", refundDetails.TxHash)
			}
			if refundDetails.CompletedAt != "" {
				fmt.Printf("  Completed At: %s\n", refundDetails.CompletedAt)
			}
			if refundDetails.FailedReason != "" {
				fmt.Printf("  ⚠️  Failed Reason: %s\n", refundDetails.FailedReason)
			}
		}
	}

	// 4. List All Refunds for Wallet
	fmt.Println("\n4. Listing All Refunds for Wallet...")
	listOptions := &interlace.ListBlockchainRefundsOptions{
		WalletID: walletID,
		Limit:    10,
		Page:     1,
	}

	refundList, err := client.BlockchainRefund.ListBlockchainRefunds(ctx, listOptions)
	if err != nil {
		log.Printf("Error listing refunds: %v\n", err)
	} else {
		fmt.Printf("Found %d refunds (showing %d):\n", refundList.TotalCount, len(refundList.Refunds))
		for i, r := range refundList.Refunds {
			fmt.Printf("  %d. Refund %s\n", i+1, r.RefundID[:16]+"...")
			fmt.Printf("     Amount: %.2f %s | Chain: %s | Status: %s\n",
				r.Amount, r.Currency, r.Chain, r.Status)
			fmt.Printf("     Created: %s\n", r.CreatedAt)
			if r.TxHash != "" {
				fmt.Printf("     TX: %s...%s\n", r.TxHash[:10], r.TxHash[len(r.TxHash)-8:])
			}
		}
		fmt.Printf("  Page: %d, Limit: %d\n", refundList.Page, refundList.Limit)
	}

	// 5. Filter by Status
	fmt.Println("\n5. Filtering Refunds by Status (COMPLETED)...")
	completedOptions := &interlace.ListBlockchainRefundsOptions{
		WalletID: walletID,
		Status:   "COMPLETED",
		Limit:    5,
	}

	completedRefunds, err := client.BlockchainRefund.ListBlockchainRefunds(ctx, completedOptions)
	if err != nil {
		log.Printf("Error listing completed refunds: %v\n", err)
	} else {
		fmt.Printf("Found %d completed refunds:\n", completedRefunds.TotalCount)
		for _, r := range completedRefunds.Refunds {
			fmt.Printf("  - %s: %.2f %s (Completed: %s)\n",
				r.RefundID[:16]+"...", r.Amount, r.Currency, r.CompletedAt)
		}
	}

	// 6. Filter by Chain
	fmt.Println("\n6. Filtering Refunds by Chain (ETHEREUM)...")
	chainOptions := &interlace.ListBlockchainRefundsOptions{
		WalletID: walletID,
		Chain:    "ETHEREUM",
		Limit:    10,
	}

	ethereumRefunds, err := client.BlockchainRefund.ListBlockchainRefunds(ctx, chainOptions)
	if err != nil {
		log.Printf("Error listing Ethereum refunds: %v\n", err)
	} else {
		fmt.Printf("Found %d Ethereum refunds\n", ethereumRefunds.TotalCount)
	}

	// 7. Filter by Date Range
	fmt.Println("\n7. Filtering Refunds by Date Range...")
	dateOptions := &interlace.ListBlockchainRefundsOptions{
		WalletID:  walletID,
		StartTime: "2024-01-01T00:00:00Z",
		EndTime:   "2024-12-31T23:59:59Z",
		Limit:     20,
	}

	dateFilteredRefunds, err := client.BlockchainRefund.ListBlockchainRefunds(ctx, dateOptions)
	if err != nil {
		log.Printf("Error listing refunds by date: %v\n", err)
	} else {
		fmt.Printf("Found %d refunds in 2024\n", dateFilteredRefunds.TotalCount)
	}

	// 8. Filter by Transfer ID
	fmt.Println("\n8. Finding Refunds for Specific Transfer...")
	transferOptions := &interlace.ListBlockchainRefundsOptions{
		TransferID: transferID,
		Limit:      5,
	}

	transferRefunds, err := client.BlockchainRefund.ListBlockchainRefunds(ctx, transferOptions)
	if err != nil {
		log.Printf("Error listing transfer refunds: %v\n", err)
	} else {
		fmt.Printf("Found %d refunds for transfer %s\n",
			transferRefunds.TotalCount, transferID[:16]+"...")
		for _, r := range transferRefunds.Refunds {
			fmt.Printf("  - Refund %s: %.2f %s (%s)\n",
				r.RefundID[:16]+"...", r.Amount, r.Currency, r.Status)
		}
	}

	fmt.Println("\n=== Blockchain Refund Workflow Complete ===")

	fmt.Println("\nTypical Refund Flow:")
	fmt.Println("1. User requests refund for a blockchain transfer")
	fmt.Println("2. Check gas fee estimate for the refund")
	fmt.Println("3. Create refund transaction")
	fmt.Println("4. Monitor refund status (PENDING → PROCESSING → COMPLETED)")
	fmt.Println("5. Verify transaction on blockchain explorer")
	fmt.Println("6. Track all refunds in system")

	fmt.Println("\n=== Refund Status Flow ===")
	fmt.Println("PENDING     → Refund created, waiting for processing")
	fmt.Println("PROCESSING  → Refund being processed on blockchain")
	fmt.Println("COMPLETED   → Refund successfully sent to recipient")
	fmt.Println("FAILED      → Refund failed (check failedReason)")

	fmt.Println("\n=== Best Practices ===")
	fmt.Println("• Always check gas fee estimate before creating refund")
	fmt.Println("• Store merchantRefundNo for tracking")
	fmt.Println("• Monitor refund status regularly")
	fmt.Println("• Verify transaction hash on blockchain explorer")
	fmt.Println("• Handle failed refunds appropriately")
	fmt.Println("• Consider gas fee in refund amount calculations")
	fmt.Println("• Validate recipient address before refund")
	fmt.Println("• Keep audit trail of all refund requests")

	fmt.Println("\n=== Supported Chains ===")
	fmt.Println("• ETHEREUM (ETH, USDT, USDC, DAI, etc.)")
	fmt.Println("• TRON (TRX, USDT-TRC20, etc.)")
	fmt.Println("• BITCOIN (BTC)")
	fmt.Println("• POLYGON (MATIC, USDT, USDC, etc.)")
	fmt.Println("• BSC (BNB, USDT-BEP20, BUSD, etc.)")
	fmt.Println("• And more...")

	fmt.Println("\n=== Important Notes ===")
	fmt.Println("⚠️  Gas fees are deducted from wallet balance")
	fmt.Println("⚠️  Refunds are irreversible once confirmed on-chain")
	fmt.Println("⚠️  Double-check recipient address before confirming")
	fmt.Println("⚠️  Gas prices fluctuate - estimate may vary")
	fmt.Println("⚠️  Some chains may have minimum refund amounts")
}
