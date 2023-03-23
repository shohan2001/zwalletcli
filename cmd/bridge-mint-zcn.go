package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/gosdk/zcnbridge/wallet"
	"github.com/0chain/gosdk/zcncore"
)

func init() {
	rootCmd.AddCommand(
		createCommandWithBridge(
			"bridge-mint-zcn",
			"mint zcn tokens using the hash of Ethereum burn transaction",
			"mint zcn tokens after burning WZCN tokens in Ethereum chain",
			commandMintZCN,
		))
}

func commandMintZCN(b *zcnbridge.BridgeClient, args ...*Arg) {
	var mintNonce int64
	cb := wallet.NewZCNStatus(&mintNonce)

	cb.Begin()

	err := zcncore.GetMintNonce(cb)
	if err != nil {
		ExitWithError(err)
	}

	if err := cb.Wait(); err != nil {
		ExitWithError(err)
	}

	if !cb.Success {
		ExitWithError(cb.Err)
	}

	burnTickets, err := b.GetNotProcessedWZCNBurnTickets(context.Background(), mintNonce)
	if err != nil {
		ExitWithError(err)
	}

	fmt.Printf("Found %d not processed WZCN burn transactions\n", len(burnTickets))

	for _, burnTicket := range burnTickets {
		fmt.Printf("Query ticket for Ethereum transaction hash: %s\n", burnTicket.TransactionHash)

		payload, err := b.QueryZChainMintPayload(burnTicket.TransactionHash)
		if err != nil {
			ExitWithError(err)
		}

		fmt.Printf("Sending mint transaction to ZCN\n")
		fmt.Printf("Ethereum transaction ID: %s\n", payload.EthereumTxnID)
		fmt.Printf("Payload amount: %d\n", payload.Amount)
		fmt.Printf("Payload nonce: %d\n", payload.Nonce)
		fmt.Printf("Receiving ZCN ClientID: %s\n", payload.ReceivingClientID)

		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*20)
		defer cancelFunc()

		fmt.Println("Starting to mint ZCN")

		txHash, err := b.MintZCN(ctx, payload)
		if err != nil {
			ExitWithError(err)
		}

		fmt.Println("Completed ZCN mint transaction")
		fmt.Printf("Transaction hash: %s\n", txHash)

	}

	fmt.Println("Done.")
}
