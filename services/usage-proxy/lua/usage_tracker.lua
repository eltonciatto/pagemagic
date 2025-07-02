-- usage_tracker.lua: Track usage events for sites
local cjson = require "cjson"
local http = require "resty.http"

-- Get usage variables
local user_id = ngx.var.usage_user_id or "anonymous"
local project_id = ngx.var.usage_project_id or ""
local site_id = ngx.var.usage_site_id or ""
local bytes_sent = tonumber(ngx.var.body_bytes_sent) or 0
local response_time = tonumber(ngx.var.request_time) or 0

-- Create usage event
local usage_event = {
    event_type = _G.USAGE_EVENTS.PAGE_VIEW,
    user_id = user_id,
    project_id = project_id,
    site_id = site_id,
    timestamp = ngx.time(),
    metadata = {
        bytes_sent = bytes_sent,
        response_time = response_time,
        status_code = ngx.status,
        user_agent = ngx.var.http_user_agent,
        referer = ngx.var.http_referer,
        remote_addr = ngx.var.remote_addr
    }
}

-- Update local counters
local usage_cache = ngx.shared.usage_cache
usage_cache:incr("total_requests", 1, 0)
usage_cache:incr("total_bytes", bytes_sent, 0)

-- Batch events to avoid overwhelming meter service
local batch_key = "batch:" .. ngx.time()
local current_batch = usage_cache:get(batch_key)
local events = {}

if current_batch then
    events = cjson.decode(current_batch)
end

table.insert(events, usage_event)

-- Send batch every 10 events or every 60 seconds
if #events >= 10 or (ngx.time() % 60) == 0 then
    -- Send to meter service asynchronously
    local httpc = http.new()
    httpc:set_timeout(1000) -- 1 second timeout
    
    local res, err = httpc:request_uri(_G.METER_SVC_URL .. "/v1/events/batch", {
        method = "POST",
        headers = {
            ["Content-Type"] = "application/json"
        },
        body = cjson.encode({
            events = events
        })
    })
    
    if res and res.status == 200 then
        -- Clear batch on success
        usage_cache:delete(batch_key)
    else
        ngx.log(ngx.ERR, "Failed to send usage events: " .. (err or "unknown error"))
    end
    
    httpc:close()
else
    -- Store updated batch
    usage_cache:set(batch_key, cjson.encode(events), 300) -- 5 minute expiry
end
