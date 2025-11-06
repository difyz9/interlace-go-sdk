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

	// Example wallet ID
	walletID := "your-wallet-id"

	fmt.Println("=== Cryptocurrency Conversion Demo ===\n")

	// 1. Get Currency Pairs
	fmt.Println("1. Getting Available Currency Pairs...")
	pairs, err := client.Convert.GetCurrencyPairs(ctx)
	if err != nil {
		log.Printf("Error getting currency pairs: %v\n", err)
	} else {
		fmt.Printf("Found %d available trading pairs:\n", len(pairs))
		for i, pair := range pairs {
			if i < 10 { // Show first 10
				fmt.Printf("  - %s → %s (Min: %.2f, Max: %.2f) %s\n",
					pair.FromCurrency, pair.ToCurrency,
					pair.MinAmount, pair.MaxAmount,
					map[bool]string{true: "✓ Available", false: "✗ Unavailable"}[pair.Available])
			}
		}
		if len(pairs) > 10 {
			fmt.Printf("  ... and %d more pairs\n", len(pairs)-10)
		}
	}

	// 2. Get Conversion Quote (from amount)
	fmt.Println("\n2. Getting Conversion Quote (from 1000 USDT to BTC)...")
	quoteRequest1 := &interlace.GetConvertQuoteRequest{
		FromCurrency: "USDT",
		ToCurrency:   "BTC",
		FromAmount:   1000.00,
	}

	quote1, err := client.Convert.GetConvertQuote(ctx, quoteRequest1)
	if err != nil {
		log.Printf("Error getting quote: %v\n", err)
	} else {
		fmt.Printf("Quote received:\n")
		fmt.Printf("  Quote ID: %s\n", quote1.QuoteID)
		fmt.Printf("  From: %.2f %s\n", quote1.FromAmount, quote1.FromCurrency)
		fmt.Printf("  To: %.8f %s\n", quote1.ToAmount, quote1.ToCurrency)
		fmt.Printf("  Exchange Rate: %.8f\n", quote1.ExchangeRate)
		fmt.Printf("  Fee: %.2f %s\n", quote1.Fee, quote1.FromCurrency)
		fmt.Printf("  Valid Until: %s\n", quote1.ValidUntil)
		fmt.Printf("  Created At: %s\n", quote1.CreatedAt)
	}

	// 3. Get Conversion Quote (to amount)
	fmt.Println("\n3. Getting Conversion Quote (to 0.01 BTC from USDT)...")
	quoteRequest2 := &interlace.GetConvertQuoteRequest{
		FromCurrency: "USDT",
		ToCurrency:   "BTC",
		ToAmount:     0.01,
	}

	quote2, err := client.Convert.GetConvertQuote(ctx, quoteRequest2)
	if err != nil {
		log.Printf("Error getting quote: %v\n", err)
	} else {
		fmt.Printf("Quote received:\n")
		fmt.Printf("  From: %.2f %s (to get 0.01 BTC)\n", quote2.FromAmount, quote2.FromCurrency)
		fmt.Printf("  Exchange Rate: %.8f\n", quote2.ExchangeRate)
		fmt.Printf("  Fee: %.2f %s\n", quote2.Fee, quote2.FromCurrency)
	}

	// 4. Create Conversion Trade (without quote)
	fmt.Println("\n4. Creating Conversion Trade (Direct)...")
	tradeRequest := &interlace.CreateConvertTradeRequest{
		WalletID:        walletID,
		FromCurrency:    "USDT",
		ToCurrency:      "BTC",
		FromAmount:      500.00,
		MerchantTradeNo: fmt.Sprintf("CONVERT-%d", time.Now().Unix()),
	}

	trade, err := client.Convert.CreateConvertTrade(ctx, tradeRequest)
	if err != nil {
		log.Printf("Error creating trade: %v\n", err)
	} else {
		fmt.Printf("Trade created successfully!\n")
		fmt.Printf("  Trade ID: %s\n", trade.TradeID)
		fmt.Printf("  Wallet ID: %s\n", trade.WalletID)
		fmt.Printf("  From: %.2f %s\n", trade.FromAmount, trade.FromCurrency)
		fmt.Printf("  To: %.8f %s\n", trade.ToAmount, trade.ToCurrency)
		fmt.Printf("  Exchange Rate: %.8f\n", trade.ExchangeRate)
		fmt.Printf("  Fee: %.2f\n", trade.Fee)
		fmt.Printf("  Status: %s\n", trade.Status)
		fmt.Printf("  Merchant Trade No: %s\n", trade.MerchantTradeNo)
		fmt.Printf("  Created At: %s\n", trade.CreatedAt)
		if trade.CompletedAt != "" {
			fmt.Printf("  Completed At: %s\n", trade.CompletedAt)
		}
	}

	// 5. Create Conversion Trade (with quote)
	fmt.Println("\n5. Creating Conversion Trade (Using Quote)...")
	if quote1 != nil {
		tradeWithQuoteRequest := &interlace.CreateConvertTradeRequest{
			WalletID:        walletID,
			FromCurrency:    "USDT",
			ToCurrency:      "BTC",
			FromAmount:      1000.00,
			QuoteID:         quote1.QuoteID, // Use the quote to lock in the rate
			MerchantTradeNo: fmt.Sprintf("CONVERT-QUOTE-%d", time.Now().Unix()),
		}

		tradeWithQuote, err := client.Convert.CreateConvertTrade(ctx, tradeWithQuoteRequest)
		if err != nil {
			log.Printf("Error creating trade with quote: %v\n", err)
		} else {
			fmt.Printf("Trade created with quoted rate!\n")
			fmt.Printf("  Trade ID: %s\n", tradeWithQuote.TradeID)
			fmt.Printf("  Exchange Rate: %.8f (locked)\n", tradeWithQuote.ExchangeRate)
			fmt.Printf("  Status: %s\n", tradeWithQuote.Status)
		}
	}

	// 6. List Conversion Trades
	fmt.Println("\n6. Listing Recent Conversion Trades...")
	listOptions := &interlace.ListConvertTradesOptions{
		WalletID:     walletID,
		FromCurrency: "USDT",
		ToCurrency:   "BTC",
		Status:       "COMPLETED",
		Limit:        10,
		Page:         1,
	}

	tradeList, err := client.Convert.ListConvertTrades(ctx, listOptions)
	if err != nil {
		log.Printf("Error listing trades: %v\n", err)
	} else {
		fmt.Printf("Found %d trades (showing %d):\n", tradeList.TotalCount, len(tradeList.Trades))
		for i, t := range tradeList.Trades {
			fmt.Printf("  %d. Trade %s: %.2f %s → %.8f %s (Rate: %.8f, Status: %s)\n",
				i+1, t.TradeID[:8]+"...",
				t.FromAmount, t.FromCurrency,
				t.ToAmount, t.ToCurrency,
				t.ExchangeRate, t.Status)
		}
		fmt.Printf("  Page: %d, Limit: %d\n", tradeList.Page, tradeList.Limit)
	}

	// 7. Filter by Status
	fmt.Println("\n7. Listing Pending Trades...")
	pendingOptions := &interlace.ListConvertTradesOptions{
		WalletID: walletID,
		Status:   "PENDING",
		Limit:    5,
	}

	pendingTrades, err := client.Convert.ListConvertTrades(ctx, pendingOptions)
	if err != nil {
		log.Printf("Error listing pending trades: %v\n", err)
	} else {
		fmt.Printf("Found %d pending trades:\n", pendingTrades.TotalCount)
		for _, t := range pendingTrades.Trades {
			fmt.Printf("  - %s: %.2f %s → %s (Created: %s)\n",
				t.TradeID[:8]+"...", t.FromAmount, t.FromCurrency, t.ToCurrency, t.CreatedAt)
		}
	}

	// 8. Filter by Date Range
	fmt.Println("\n8. Listing Trades by Date Range...")
	dateOptions := &interlace.ListConvertTradesOptions{
		WalletID:  walletID,
		StartTime: "2024-01-01T00:00:00Z",
		EndTime:   "2024-12-31T23:59:59Z",
		Limit:     20,
	}

	dateFilteredTrades, err := client.Convert.ListConvertTrades(ctx, dateOptions)
	if err != nil {
		log.Printf("Error listing trades by date: %v\n", err)
	} else {
		fmt.Printf("Found %d trades in 2024:\n", dateFilteredTrades.TotalCount)
		fmt.Printf("  Showing %d trades\n", len(dateFilteredTrades.Trades))
	}

	fmt.Println("\n=== Conversion Workflow Complete ===")
	fmt.Println("\nTypical Conversion Flow:")
	fmt.Println("1. Check available currency pairs")
	fmt.Println("2. Get a quote to see the exchange rate")
	fmt.Println("3. Create trade (optionally using quote ID to lock rate)")
	fmt.Println("4. Monitor trade status")
	fmt.Println("5. Check transaction history")

	fmt.Println("\n=== Best Practices ===")
	fmt.Println("• Use quotes to lock favorable exchange rates")
	fmt.Println("• Check quote validity before creating trade")
	fmt.Println("• Monitor trade status for completion")
	fmt.Println("• Set up merchant trade numbers for tracking")
	fmt.Println("• Consider fees when calculating amounts")
	fmt.Println("• Be aware of min/max amounts for each pair")
}
