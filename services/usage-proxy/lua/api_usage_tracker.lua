-- api_usage_tracker.lua: Track API usage events
local cjson = require "cjson"
local http = require "resty.http"

-- Get API usage variables
local user_id = ngx.var.usage_user_id or "anonymous"
local project_id = ngx.var.usage_project_id or ""
local bytes_sent = tonumber(ngx.var.body_bytes_sent) or 0
local response_time = tonumber(ngx.var.request_time) or 0
local request_method = ngx.var.request_method
local request_uri = ngx.var.request_uri

-- Determine API endpoint for more granular tracking
local endpoint = request_uri:match("^/api/([^/]+)")
if not endpoint then
    endpoint = "unknown"
end

-- Create API usage event
local api_event = {
    event_type = _G.USAGE_EVENTS.API_CALL,
    user_id = user_id,
    project_id = project_id,
    timestamp = ngx.time(),
    metadata = {
        endpoint = endpoint,
        method = request_method,
        uri = request_uri,
        bytes_sent = bytes_sent,
        response_time = response_time,
        status_code = ngx.status,
        user_agent = ngx.var.http_user_agent,
        remote_addr = ngx.var.remote_addr
    }
}

-- Update API-specific counters
local usage_cache = ngx.shared.usage_cache
usage_cache:incr("api_calls_total", 1, 0)
usage_cache:incr("api_bytes_total", bytes_sent, 0)
usage_cache:incr("api_" .. endpoint .. "_calls", 1, 0)

-- Send to meter service
local httpc = http.new()
httpc:set_timeout(1000)

local res, err = httpc:request_uri(_G.METER_SVC_URL .. "/v1/events", {
    method = "POST",
    headers = {
        ["Content-Type"] = "application/json"
    },
    body = cjson.encode(api_event)
})

if not res or res.status ~= 200 then
    ngx.log(ngx.ERR, "Failed to send API usage event: " .. (err or "unknown error"))
end

httpc:close()
