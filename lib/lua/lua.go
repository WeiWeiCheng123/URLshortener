package lua

const IP_script = `
local key = KEYS[1]
local ipLimit = tonumber(ARGV[1])
local period = tonumber(ARGV[2])
local userInfo = redis.call('GET', key)
println("key ",key)
println("iplim ",ipLimit)
println("period ",period)
println("userInfo ",userInfo)
return userInfo
`

/*
if #userInfo == 0 then
    redis.call('SET', key, 1)
	redis.call('EXPIRE',period)
    result = 1
    return result
end
local count = tonumber(userInfo)
if count < ipLimit then
    redis.call('INCR', key)
    result = 1
    return result
else
    result = -1
    return result
end
`
*/
