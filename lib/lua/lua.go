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
return -1
`
	Save_URL = `
local shortURL = KEYS[1]
local URL = ARGV[1]
local exp = ARGV[2]
redis.call('SET', shortURL, URL)
redis.call('EXPIREAT', shortURL, exp)
`
	Set_NotExist = `
local shortURL = KEYS[1]
redis.call('SET', shortURL, "NotExist")
redis.call('EXPIRE', shortURL, 600)
`
)
