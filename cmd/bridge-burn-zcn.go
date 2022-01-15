package cmd

import (
	"context"
	"fmt"

	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	rootCmd.AddCommand(
		createCommandWithBridge(
			"bridge-burn-zcn",
			"burn zcn tokens",
			"burn zcn tokens that will be minted for WZCN tokens",
			commandBurnZCN,
			amountOption,
		))
}

func commandBurnZCN(b *zcnbridge.BridgeClient, args ...*Arg) {
	amount := GetAmount(args)

	fmt.Printf("Starting burn transaction")
	transaction, err := b.BurnZCN(context.Background(), amount)
	if err == nil {
		fmt.Printf("Submitted burn transaction %s\n", transaction.Hash)
	} else {
		ExitWithError(err)
	}

	fmt.Printf("Starting transaction verification %s\n", transaction.Hash)
	verify(transaction.Hash)
}
