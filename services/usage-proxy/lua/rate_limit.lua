-- rate_limit.lua: Implement rate limiting
local rate_limit = ngx.shared.rate_limit

-- Get user ID (set by auth_check.lua)
local user_id = ngx.var.usage_user_id or ngx.var.remote_addr

-- Rate limit keys
local minute_key = "rate:" .. user_id .. ":" .. math.floor(ngx.time() / 60)
local hour_key = "rate:" .. user_id .. ":" .. math.floor(ngx.time() / 3600)
local day_key = "rate:" .. user_id .. ":" .. math.floor(ngx.time() / 86400)

-- Get current counts
local minute_count = rate_limit:get(minute_key) or 0
local hour_count = rate_limit:get(hour_key) or 0
local day_count = rate_limit:get(day_key) or 0

-- Check limits
if minute_count >= _G.RATE_LIMITS.DEFAULT_PER_MINUTE then
    ngx.log(ngx.WARN, "Rate limit exceeded for user " .. user_id .. " (per minute)")
    ngx.status = 429
    ngx.header["Retry-After"] = "60"
    ngx.say("Rate limit exceeded")
    ngx.exit(429)
end

if hour_count >= _G.RATE_LIMITS.DEFAULT_PER_HOUR then
    ngx.log(ngx.WARN, "Rate limit exceeded for user " .. user_id .. " (per hour)")
    ngx.status = 429
    ngx.header["Retry-After"] = "3600"
    ngx.say("Rate limit exceeded")
    ngx.exit(429)
end

if day_count >= _G.RATE_LIMITS.DEFAULT_PER_DAY then
    ngx.log(ngx.WARN, "Rate limit exceeded for user " .. user_id .. " (per day)")
    ngx.status = 429
    ngx.header["Retry-After"] = "86400"
    ngx.say("Rate limit exceeded")
    ngx.exit(429)
end

-- Increment counters
rate_limit:incr(minute_key, 1, 0, 60) -- Expire after 1 minute
rate_limit:incr(hour_key, 1, 0, 3600) -- Expire after 1 hour
rate_limit:incr(day_key, 1, 0, 86400) -- Expire after 1 day
