import subprocess
import json
from dotenv import load_dotenv
import os

def get_block(block_hash):
    """
    Fetch block information using the block hash.
    """
    try:
        block_info = subprocess.check_output(
            ["bitcoin-cli", "getblock", block_hash],
            universal_newlines=True
        )
        return json.loads(block_info)
    except subprocess.CalledProcessError as e:
        print(f"Error fetching block: {e}")
        return None

def get_raw_transaction(tx_id, block_hash):
    """
    Fetch raw transaction details by transaction ID and block hash.
    """
    try:
        raw_tx = subprocess.check_output(
            ["bitcoin-cli", "getrawtransaction", tx_id, "true", block_hash],
            universal_newlines=True
        )
        return json.loads(raw_tx)
    except subprocess.CalledProcessError as e:
        print(f"Error fetching transaction: {e}")
        return None

def extract_op_return(transaction):
    """
    Extract OP_RETURN data from the transaction outputs.
    """
    for output in transaction.get("vout", []):
        script = output.get("scriptPubKey", {})
        if script.get("type") == "nulldata":
            return script.get("asm").replace("OP_RETURN ", "")
    return None

def main():
    # Prompt user for block hash and transaction ID
    block_hash = input("Enter the block hash: ").strip()
    tx_id = input("Enter the transaction ID: ").strip()

    def get_block(block_hash):
        try:
            result = subprocess.check_output(['bitcoin-cli', 'getblock', block_hash, '2'])
            return json.loads(result)
        except subprocess.CalledProcessError as e:
            print(f"Error fetching block: {e}")
            return None
    block = get_block(block_hash)
    print(block)

    # # Fetch block and transaction details
    # block = get_block(block_hash)
    # if not block:
    #     print("Invalid block hash or unable to fetch block data.")
    #     return
    # print(block)

    # transaction = get_raw_transaction(tx_id, block_hash)
    # if not transaction:
    #     print("Invalid transaction ID or unable to fetch transaction data.")
    #     return

    # # Extract OP_RETURN data
    # op_return_data = extract_op_return(transaction)
    # if op_return_data:
    #     print(f"OP_RETURN data: {op_return_data}")
    # else:
    #     print("No OP_RETURN data found in the transaction.")

if __name__ == "__main__":
    load_dotenv()
    # Example input and Bitcoin node RPC credentials
    # txid = input("Enter the TXID of the broadcasted transaction: ")
    rpc_url = f'http://{os.getenv("HOST")}'  # Adjust the URL to your Bitcoin node
    rpc_user = os.getenv("BTCUSER")
    rpc_password = os.getenv("PASSWORD")
    main()
