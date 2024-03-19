from bitcoin.rpc import RawProxy
# import datetime

proxy = RawProxy()

def get_best_block():
    block = proxy.getblock(proxy.getbestblockhash())
    max_block_height = block["height"]
    return max_block_height

best_block = get_best_block()
print(best_block)