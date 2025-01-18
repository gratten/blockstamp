import json
import requests
from bitcoinlib.transactions import Transaction
import binascii
import os
from dotenv import load_dotenv

def get_transaction_details(txid, rpc_url, rpc_user, rpc_password):
    # Bitcoin node RPC URL and credentials
    url = rpc_url
    headers = {'Content-Type': 'application/json'}
    data = {
        "jsonrpc": "1.0",
        "id": "1",
        "method": "gettransaction",
        "params": [txid]
    }
    
    # Send the RPC request
    response = requests.post(url, data=json.dumps(data), headers=headers, auth=(rpc_user, rpc_password))
    
    # Check if the response is valid
    if response.status_code == 200:
        result = response.json()
        if 'error' in result and result['error']:
            print("Error:", result['error'])
            return None
        return result['result']
    else:
        print(f"Error {response.status_code}: Unable to fetch transaction details")
        return None

# Extract OP_RETURN data from the transaction outputs
def extract_op_return(tx_data):
    for output in tx_data['outputs']:
        if output['script_type'] == 'nulldata':  # This is the OP_RETURN type
            op_return_hex = output['script']
            # Decode the hex string (remove the "6a" prefix and decode the remaining hex)
            op_return_data = binascii.unhexlify(op_return_hex[2:]).decode('utf-8')
            return op_return_data
    return None


def analyze_transaction(txid, rpc_url, rpc_user, rpc_password):
    # Get transaction details from Bitcoin node
    tx_details = get_transaction_details(txid, rpc_url, rpc_user, rpc_password)
    # print(tx_details)
    
    if not tx_details:
        return
    
    # Print basic transaction information
    print(f"Transaction ID: {txid}")
    print(f"Amount: {tx_details['amount']}")
    print(f"Fee: {tx_details['fee']}")
    print(f"Confirmations: {tx_details['confirmations']}")
    # print(f"Hex: {tx_details['hex']}")

    transaction_data = Transaction.parse(tx_details['hex']).as_dict()

    op_return_data = extract_op_return(transaction_data)

    # decoded_hex = decode_hex(tx_details['hex'])

    print(f"OP_RETURN: {op_return_data}")


if __name__ == '__main__':
    load_dotenv()
    # Example input and Bitcoin node RPC credentials
    txid = input("Enter the TXID of the broadcasted transaction: ")
    rpc_url = f'http://{os.getenv("HOST")}'  # Adjust the URL to your Bitcoin node
    rpc_user = os.getenv("BTCUSER")
    rpc_password = os.getenv("PASSWORD")
    
    # Analyze the transaction
    analyze_transaction(txid, rpc_url, rpc_user, rpc_password)
