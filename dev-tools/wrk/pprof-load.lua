wrk.method = "POST"
wrk.headers["Content-Type"] = "text/plain"

function request()
    local r = math.random()

    if r > 0.8 then
        return createShortURL()
    else
        return followOriginalURL()
    end
end

function followOriginalURL()
    return wrk.format("GET", "/tEsT")
end

function createShortURL()
    local random_id = math.random(10000)
    local body = string.format('https://practicum.yandex.ru/%d', random_id)
    return wrk.format("POST", nil, nil, body)
end
