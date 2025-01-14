package main

import (
	// "bytes"
	// "encoding/hex"
	// "fmt"
	"log"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func transaction() {

	// User inputs
	targetBlockHeight := 800000  // Replace with user-provided block height
	message := "Hello, Bitcoin!" // Replace with user-provided text
	feeRate := int64(10)         // Fee rate in satoshis per byte

	// Create OP_RETURN output
	opReturnScript, err := txscript.NullDataScript([]byte(message))
	if err != nil {
		log.Fatalf("Failed to create OP_RETURN script: %v", err)
	}

	// Select UTXOs
	utxos, err := client.ListUnspent()
	if err != nil {
		log.Fatalf("Failed to list UTXOs: %v", err)
	}

	if len(utxos) == 0 {
		log.Fatalf("No UTXOs available for spending")
	}

	selectedUTXO := utxos[0] // Simple selection: Use the first UTXO
	txHash, err := chainhash.NewHashFromStr(selectedUTXO.TxID)
	if err != nil {
		log.Fatalf("Invalid UTXO TxID: %v", err)
	}

	// Build the transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// Add input
	outPoint := wire.NewOutPoint(txHash, selectedUTXO.Vout)
	txIn := wire.NewTxIn(outPoint, nil, nil)
	tx.AddTxIn(txIn)

	// Add OP_RETURN output
	txOut := wire.NewTxOut(0, opReturnScript)
	tx.AddTxOut(txOut)

	// Add change output if necessary
	changeAmount := int64(selectedUTXO.Amount*1e8) - feeRate*int64(tx.SerializeSize())
	if changeAmount > 0 {
		changeAddress, err := btcutil.DecodeAddress(selectedUTXO.Address, &chaincfg.MainNetParams)
		if err != nil {
			log.Fatalf("Failed to decode change address: %v", err)
		}
		changeScript, err := txscript.PayToAddrScript(changeAddress)
		if err != nil {
			log.Fatalf("Failed to create change script: %v", err)
		}
		tx.AddTxOut(wire.NewTxOut(changeAmount, changeScript))
	}

	// Set nLockTime
	tx.LockTime = uint32(targetBlockHeight)

	log.Printf("Transaction: %+v\n", tx)
	for i, txIn := range tx.TxIn {
		log.Printf("Input %d: %v\n", i, txIn)
	}
	for i, txOut := range tx.TxOut {
		log.Printf("Output %d: %v\n", i, txOut)
	}
	log.Printf("LockTime: %d\n", tx.LockTime)
	// latestBlockHash, err := client.GetBlockHash(blockCount)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
