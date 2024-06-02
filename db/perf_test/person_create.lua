local wrk = require("wrk")

wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

request = function()
    math.randomseed(os.time())

    local name = generateRandomString()
    local email = generateRandomEmail()
    local password = generateRandomPassword()

    local body = string.format('{"name":"%s", "birthday":"%d", "email":"%s", "password":"%s", "gender":"%s"}', 
    name, 20010101, email, password, "male")
    return wrk.format(nil, "/api/v1/registration", nil, body)
end

-- Function to generate a random string
function generateRandomString()
    local length = 10 -- Length of the random string
    local characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" -- Characters to choose from
    local randomString = ""
    
    for i = 1, length do
        local randomIndex = math.random(1, #characters)
        randomString = randomString .. string.sub(characters, randomIndex, randomIndex)
    end
    
    return randomString
end

-- Function to generate a random email address
function generateRandomEmail()
    local emailLength = 10 -- Length of the email address
    local domain = "@example.com" -- Domain name for the email address
    local characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" -- Characters to choose from
    local randomEmail = ""

    for i = 1, emailLength do
        local randomIndex = math.random(1, #characters)
        randomEmail = randomEmail .. string.sub(characters, randomIndex, randomIndex)
    end

    return randomEmail .. domain
end

-- Function to generate a random password
function generateRandomPassword()
    local passwordLength = 8 -- Length of the password
    local characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" -- Characters to choose from
    local randomPassword = ""

    for i = 1, passwordLength do
        local randomIndex = math.random(1, #characters)
        randomPassword = randomPassword .. string.sub(characters, randomIndex, randomIndex)
    end

    return randomPassword
end