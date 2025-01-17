package main

import (
	// "encoding/hex"
	// "fmt"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func transaction(blockheight int, stamp string) {

	log.Println(blockheight)
	log.Println(stamp)

	// User inputs
	targetBlockHeight := blockheight // Replace with user-provided block height
	message := stamp                 // Replace with user-provided text
	feeRate := int64(10)             // Fee rate in satoshis per byte

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

	// log.Printf("Transaction: %+v\n", tx)
	// for i, txIn := range tx.TxIn {
	// 	log.Printf("Input %d: %v\n", i, txIn)
	// }
	// for i, txOut := range tx.TxOut {
	// 	log.Printf("Output %d: %v\n", i, txOut)
	// }
	// log.Printf("LockTime: %d\n", tx.LockTime)
	// latestBlockHash, err := client.GetBlockHash(blockCount)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	log.Println("priv key: ", generatePrivateKey())

	// Sign the transaction
	privateKeyWIF := generatePrivateKey() // Replace with your private key
	wif, err := btcutil.DecodeWIF(privateKeyWIF)
	if err != nil {
		log.Fatalf("Failed to decode WIF: %v", err)
	}
	// Convert ScriptPubKey from string to []byte
	// scriptPubKeyBytes := []byte(selectedUTXO.ScriptPubKey)
	scriptPubKeyBytes, err := hex.DecodeString(selectedUTXO.ScriptPubKey)
	if err != nil {
		log.Fatalf("Failed to decode ScriptPubKey: %v", err)
	}
	log.Printf("ScriptPubKey (hex): %x", scriptPubKeyBytes)
	// Create a PrevOutputFetcher to provide UTXO details
	prevOutFetcher := txscript.NewCannedPrevOutputFetcher(
		scriptPubKeyBytes,          // The ScriptPubKey as a []byte
		int64(selectedUTXO.Amount), // The value in satoshis of the UTXO
	)
	// Prepare the transaction's hash cache for SegWit
	hashCache := txscript.NewTxSigHashes(tx, prevOutFetcher)

	// Generate the witness signature
	witness, err := txscript.WitnessSignature(
		tx,                         // Transaction to sign
		hashCache,                  // Transaction signature hashes
		0,                          // Index of the input being signed
		int64(selectedUTXO.Amount), // Value of the UTXO in satoshis
		scriptPubKeyBytes,          // The UTXO's ScriptPubKey
		txscript.SigHashAll,        // Hash type
		wif.PrivKey,                // Private key
		true,                       // Compress public key
	)
	if err != nil {
		log.Fatalf("Failed to create witness signature: %v", err)
	}

	// Assign the witness stack to the input
	tx.TxIn[0].Witness = witness

	// Log the generated witness
	log.Printf("Generated witness: %x", witness)
	// sigScript, err := txscript.SignatureScript(tx, 0, scriptPubKeyBytes, txscript.SigHashAll, wif.PrivKey, true)
	// log.Printf("ScriptPubKey (hex): %s", selectedUTXO.ScriptPubKey)
	// log.Printf("ScriptPubKey (bytes): %x", scriptPubKeyBytes)
	// log.Printf("Generated signature: %x", sigScript)

	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}
	// txIn.SignatureScript = sigScript

	// Serialize and broadcast the transaction
	var buf bytes.Buffer
	err = tx.Serialize(&buf)
	if err != nil {
		log.Fatalf("Failed to serialize transaction: %v", err)
	}

	// txHex := hex.EncodeToString(buf.Bytes())
	txID, err := client.SendRawTransaction(tx, false)
	if err != nil {
		log.Fatalf("Failed to broadcast transaction: %v", err)
	}

	fmt.Printf("Transaction broadcasted! TXID: %s\n", txID.String())
}

func generatePrivateKey() string {
	// Generate a new private key using btcec
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		log.Fatalf("Error generating private key: %v", err)
	}

	// Convert the private key to WIF format
	wif, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, true)
	if err != nil {
		log.Fatalf("Error converting private key to WIF: %v", err)
	}

	return wif.String() // Return the private key in WIF format
}

// func generatePrivateKey() []byte {
// 	// Generate a new private key using btcec
// 	privKey, err := btcec.NewPrivateKey()
// 	if err != nil {
// 		log.Fatalf("Error generating private key: %v", err)
// 	}

// 	// Return the serialized private key as a byte slice
// 	return privKey.Serialize()
// }
