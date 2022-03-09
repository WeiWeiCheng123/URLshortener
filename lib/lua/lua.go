package lua

const (
	IP_script = `
local ip = KEYS[1]
local ipLimit = tonumber(ARGV[1])
local period = tonumber(ARGV[2])
local count = redis.call('GET', ip)
if (not count) then
    redis.call('SET', ip, 1)
    redis.call('EXPIRE', ip, period)
    return 1
end
if (tonumber(count) < ipLimit) then
    redis.call('INCR', ip)
    redis.call('EXPIRE', ip, period)
    return 1
end
redis.call('EXPIRE', ip, period)
return -1
`
)
