-- init_worker.lua: Initialize worker-specific resources
local cjson = require "cjson"

-- Initialize metrics counters
local usage_cache = ngx.shared.usage_cache
local rate_limit = ngx.shared.rate_limit

-- Set initial values
usage_cache:set("total_requests", 0)
usage_cache:set("total_bytes", 0)
usage_cache:set("active_sites", 0)

-- Worker ID for logging
local worker_id = ngx.worker.id()
ngx.log(ngx.INFO, "Worker " .. worker_id .. " initialized")
