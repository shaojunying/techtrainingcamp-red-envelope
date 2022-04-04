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
    return wrk.format(nil, path, nil, body)
end