import subprocess
import json
import time

def run_cli_command(command):
    # Print the exact command
    # print("Executing command:", " ".join(command))
    
    # Run the command
    result = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    
    if result.returncode != 0:
        # Print error if any
        print("Error:", result.stderr.decode("utf-8"))
        raise Exception("Command failed")
    return result.stdout.decode("utf-8")

def get_unspent_output():
    """Get an unspent output to use as an input for the raw transaction."""
    command = ["bitcoin-cli", "-regtest", "listunspent", "0", "9999999"]
    unspent_outputs = run_cli_command(command)
    
    if unspent_outputs:
        unspent_list = json.loads(unspent_outputs)
        if unspent_list:  # Check if there are any unspent outputs
            # Assuming the first unspent output is good for the transaction
            unspent = unspent_list[0]
            fee = 0.01
            amount = round(unspent['amount'] - fee, 4)
            return unspent['txid'], unspent['vout'], amount  # Return txid, vout, amount
        else:
            print("No unspent outputs available.")
            return None, None, None
    else:
        print("Error: Unable to fetch unspent outputs.")
        return None, None, None

def create_op_return_transaction(txid, vout, address, amount, op_return_data):
    # Format the inputs and outputs correctly
    inputs = [{"txid": txid, "vout": vout}]
    outputs = {address: amount, "data": op_return_data}
    
    # Convert to JSON strings, then enclose in single quotes
    inputs_str = json.dumps(inputs, separators=(",", ":"))  # No spaces
    outputs_str = json.dumps(outputs, separators=(",", ":"))  # No spaces
    
    # Build the command
    command = [
        "bitcoin-cli", 
        "-regtest", 
        "createrawtransaction", 
        inputs_str,
        outputs_str
    ]
    
    # Run the command
    raw_transaction = run_cli_command(command)
    # print("Transaction created with hex:", raw_transaction)
    return raw_transaction

def sign_transaction(transaction_hex):
    command = ["bitcoin-cli", "-regtest", "signrawtransactionwithwallet", transaction_hex]
    return run_cli_command(command)

def send_transaction(signed_hex_response):
    response = json.loads(signed_hex_response)
    signed_hex = response["hex"]
    command = ["bitcoin-cli", "-regtest", "sendrawtransaction", signed_hex]
    return run_cli_command(command)

def wait_for_block(target_block_height):
    """Wait for the specified block height."""
    while True:
        current_height = get_current_block_height()
        if current_height >= target_block_height:
            break
        print(f"Waiting... Current block height: {current_height}")
        time.sleep(5)  # Wait for 5 seconds before checking again

def get_current_block_height():
    """Get the current block height"""
    command = ["bitcoin-cli", "-regtest", "getblockchaininfo"]
    response = json.loads(run_cli_command(command))
    # print(response)
    return response['blocks']

def main():
    """Main application logic."""
    # current_block_height_command = ["bitcoin-cli", "-regtest", "getblockcount"]
    # current_block_height = int(run_cli_command(current_block_height_command))
    current_block_height = get_current_block_height()
    print(f"Current block height: {current_block_height}")

    # Step 2: Ask if they want to stamp the blockchain
    stamp_response = input("Would you like to stamp the blockchain? (yes/no): ").strip().lower()
    if stamp_response != "yes":
        print("Exiting. No stamp will be added.")
        return

    stamp_block_height = int(input("Enter a future block height: "))
    stamp_text = input("Enter the text you want to include in the OP_RETURN: ")
    op_return_data = stamp_text.encode('utf-8').hex()

    txid, vout, amount = get_unspent_output()

    # Create a change output
    address = run_cli_command(["bitcoin-cli", "-regtest", "getnewaddress"]).strip()

    transaction_hex = create_op_return_transaction(txid, vout, address, amount, op_return_data)
    # print("Transaction Hex:", transaction_hex)
    signed_hex_response = sign_transaction(transaction_hex.strip())
    new_txid = send_transaction(signed_hex_response)
    print("Transaction ID:", new_txid)
    print("Waiting for block to be mined...")
    wait_for_block(stamp_block_height)
    print("Congratulations! Your stamp has been added to the blockchain.")

if __name__ == "__main__":
    main()