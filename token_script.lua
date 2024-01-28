-- KEY[1] - bucket key
-- ARGV[1] - rate (tokens added per interval)
-- ARGV[2] - interval (milliseconds)
-- ARGV[3] - max tokens

local bucket = redis.call('HMGET', KEYS[1], 'tokens', 'last_time')
local tokens = tonumber(bucket[1])
local last_time = tonumber(bucket[2])

if tokens == nil then
    tokens = ARGV[3]
    last_time = redis.call('TIME')[1] * 1000 + redis.call('TIME')[2] / 1000
end

local current_time = redis.call('TIME')[1] * 1000 + redis.call('TIME')[2] / 1000
local delta = math.max(current_time - last_time, 0)
local fill_tokens = math.floor(delta / ARGV[2]) * ARGV[1]
tokens = math.min(tokens + fill_tokens, ARGV[3])
last_time = current_time

if tokens > 0 then
    tokens = tokens - 1
    redis.call('HMSET', KEYS[1], 'tokens', tokens, 'last_time', last_time)
    return 1
else
    return 0
end
