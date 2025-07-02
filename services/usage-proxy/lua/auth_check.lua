-- auth_check.lua: Validate JWT tokens and extract user information
local cjson = require "cjson"
local jwt = require "resty.jwt"
local http = require "resty.http"

-- Extract JWT from Authorization header
local auth_header = ngx.var.http_authorization
if not auth_header then
    -- No auth for public sites, allow through
    return
end

local token = auth_header:match("Bearer%s+(.+)")
if not token then
    ngx.log(ngx.ERR, "Invalid authorization header format")
    ngx.status = 401
    ngx.say("Unauthorized")
    ngx.exit(401)
end

-- Verify JWT token with auth service
local httpc = http.new()
httpc:set_timeout(5000) -- 5 second timeout

local res, err = httpc:request_uri(_G.AUTH_SVC_URL .. "/v1/auth/verify", {
    method = "POST",
    headers = {
        ["Content-Type"] = "application/json",
        ["Authorization"] = "Bearer " .. token
    }
})

if not res then
    ngx.log(ngx.ERR, "Failed to verify token: " .. err)
    ngx.status = 500
    ngx.say("Internal Server Error")
    ngx.exit(500)
end

if res.status ~= 200 then
    ngx.log(ngx.WARN, "Token verification failed: " .. res.status)
    ngx.status = 401
    ngx.say("Unauthorized")
    ngx.exit(401)
end

-- Parse response and extract user info
local user_info = cjson.decode(res.body)
if user_info and user_info.user_id then
    -- Set variables for logging
    ngx.var.usage_user_id = user_info.user_id
    ngx.var.usage_project_id = user_info.project_id or ""
    
    -- Store in shared dict for usage tracking
    local usage_cache = ngx.shared.usage_cache
    usage_cache:set("user:" .. user_info.user_id .. ":last_seen", ngx.time())
end

httpc:close()
