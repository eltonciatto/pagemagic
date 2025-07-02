-- metrics.lua: Expose Prometheus metrics
local usage_cache = ngx.shared.usage_cache
local rate_limit = ngx.shared.rate_limit

-- Get metrics from shared dict
local total_requests = usage_cache:get("total_requests") or 0
local total_bytes = usage_cache:get("total_bytes") or 0
local api_calls_total = usage_cache:get("api_calls_total") or 0
local api_bytes_total = usage_cache:get("api_bytes_total") or 0

-- Generate Prometheus format
local metrics = {}

-- Request metrics
table.insert(metrics, "# HELP usage_proxy_requests_total Total number of requests")
table.insert(metrics, "# TYPE usage_proxy_requests_total counter")
table.insert(metrics, "usage_proxy_requests_total " .. total_requests)

-- Bytes metrics
table.insert(metrics, "# HELP usage_proxy_bytes_total Total bytes transferred")
table.insert(metrics, "# TYPE usage_proxy_bytes_total counter")
table.insert(metrics, "usage_proxy_bytes_total " .. total_bytes)

-- API metrics
table.insert(metrics, "# HELP usage_proxy_api_calls_total Total API calls")
table.insert(metrics, "# TYPE usage_proxy_api_calls_total counter")
table.insert(metrics, "usage_proxy_api_calls_total " .. api_calls_total)

table.insert(metrics, "# HELP usage_proxy_api_bytes_total Total API bytes")
table.insert(metrics, "# TYPE usage_proxy_api_bytes_total counter")
table.insert(metrics, "usage_proxy_api_bytes_total " .. api_bytes_total)

-- Worker info
local worker_id = ngx.worker.id()
table.insert(metrics, "# HELP usage_proxy_worker_info Worker information")
table.insert(metrics, "# TYPE usage_proxy_worker_info gauge")
table.insert(metrics, "usage_proxy_worker_info{worker_id=\"" .. worker_id .. "\"} 1")

-- Output metrics
ngx.header.content_type = "text/plain"
ngx.say(table.concat(metrics, "\n"))
