from bitcoin.rpc import RawProxy

proxy = RawProxy()

def get_best_block():
    block = proxy.getblock(proxy.getbestblockhash())
    max_block_height = block["height"]
    return max_block_height

def get_block_height(max_block_height, target_time):
    left_block_height, right_block_height = 0, max_block_height

    while left_block_height <= right_block_height:
        mid_block_height = (left_block_height + right_block_height) // 2

        mid_block = proxy.getblock(proxy.getblockhash(mid_block_height))
        left_block = proxy.getblock(proxy.getblockhash(left_block_height))
        right_block = proxy.getblock(proxy.getblockhash(right_block_height))

        mid_block_time = mid_block["time"]
        left_block_time = left_block["time"]
        right_block_time = right_block["time"]

        # Check if target is present at mid
        if mid_block_time == target_time:
            return mid_block_height
        # If target is greater, ignore left half
        elif mid_block_time < target_time:
            left_block_height = mid_block_height + 1
        # If target is smaller, ignore right half
        else:
            right_block_height = mid_block_height - 1

    # If the target is not found in the array
    return mid_block_height

best_block = get_best_block()
print(best_block)