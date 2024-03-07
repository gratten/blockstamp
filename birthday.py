from bitcoin.rpc import RawProxy
import datetime

proxy = RawProxy()

def binary_search(max_block_height, target_time):
    # left, right = 0, len(arr) - 1
    # print(max_block_height)
    left_block_height, right_block_height = 0, max_block_height

    while left_block_height <= right_block_height:
        mid_block_height = (left_block_height + right_block_height) // 2

        # print(mid_block_height)
        mid_block = proxy.getblock(proxy.getblockhash(mid_block_height))
        left_block = proxy.getblock(proxy.getblockhash(left_block_height))
        right_block = proxy.getblock(proxy.getblockhash(right_block_height))

        mid_block_time = mid_block["time"]
        left_block_time = left_block["time"]
        right_block_time = right_block["time"]

        # print (mid_block_time)
        # print (left_block_time)
        # print (right_block_time)

        # print()

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

# get max block height
block = proxy.getblock(proxy.getbestblockhash())
max_block_height = block["height"]

# get target time
year = int(input('Enter year:\n'))
month = int(input('Enter month:\n'))
day = int(input('Enter day:\n'))
hour = int(input('Enter hour:\n'))
minute = int(input('Enter minute:\n'))
sec = int(input('Enter second:\n'))

# given_datetime = datetime.datetime(2020, 2, 18, 15, 38, 7)
given_datetime = datetime.datetime(year, month, day, hour, minute, sec)

# 2022, 6, 3, 15, 7, 0

target_time = int(given_datetime.timestamp())
# print(type(target_time))
# print(target_time)
# target_time = 1582056227
# target_time = 1654283220

# find result
print("Finding blockheight...")
result = binary_search(max_block_height, target_time)
print(f'The blockheight at this date and time was: {result}')

# if result != -1:
#     print(f"Element {target} is present at index {result}")
# else:
#     print(f"Element {target} is not present in the array")