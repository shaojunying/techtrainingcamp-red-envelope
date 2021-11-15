local json = require "json"

wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
local uid = 1 --- 根据数据库情况设置，
local path = "/redenvelope/snatch"
local body = '{"uid": 1}'
local data = ""

function request()
    body = '{"uid":%s}'
    uid = math.random(10000000)
    body = string.format(body, uid)
    path = "/redenvelope/snatch"
    return wrk.format(nil, path, nil, body)
end

function response(status, headers, body)
    local res = json.decode(body)
    data = res["data"]
end