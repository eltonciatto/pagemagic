-- init.lua: Initialize global variables and modules
local cjson = require "cjson"
local http = require "resty.http"

-- Global configuration
_G.METER_SVC_URL = "http://meter-svc:8080"
_G.AUTH_SVC_URL = "http://auth-svc:8080"

-- Usage event types
_G.USAGE_EVENTS = {
    PAGE_VIEW = "page_view",
    API_CALL = "api_call",
    BYTES_TRANSFERRED = "bytes_transferred",
    CONTAINER_TIME = "container_time"
}

-- Rate limit defaults
_G.RATE_LIMITS = {
    DEFAULT_PER_MINUTE = 60,
    DEFAULT_PER_HOUR = 1000,
    DEFAULT_PER_DAY = 10000
}

ngx.log(ngx.INFO, "Usage proxy initialized")
