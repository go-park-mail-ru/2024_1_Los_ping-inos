local wrk = require("wrk")

wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Cookie"] = "session_id=4081fb75-538a-4b93-af91-cc50dfc34407" 

request = function()
    math.randomseed(os.time())
    local path = string.format("/api/v1/profile/id=%d", math.random(1, 100))
    return wrk.format(nil, path, nil, nil)
end