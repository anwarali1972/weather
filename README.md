# weather
# This is to support Open APIs that can be used by clients to get weather info
# APIs
# myApp/v1/weather/current
# This is to get current weather condition based on latituse/longtitude
# Example
curl --location --request GET 'http://localhost:5555/myApp/v1/weather/current?lat=35.89375&long=-90&appid=78fe259766f8b3433838fb022bf5724a&units=imperial' --header 'Content-Type: application/json' --data-raw ''

# Result
{"condition":"mist","temperature":"moderate","alerts":[{"sender_name":"NWS Memphis (Western Tennessee, Eastern Arkansas and Northern Mississippi)","event":"Wind Advisory","start":1620561600,"end":1620597600,"description":"...WIND ADVISORY REMAINS IN EFFECT FROM 7 AM THIS MORNING TO 5 PM\nCDT THIS AFTERNOON...\n* WHAT...Southwest winds 15 to 25 mph with gusts to 40 mph.\n* WHERE...West Tennessee, East Arkansas and Southeast Missouri.\n* WHEN...From 7 AM this morning to 5 PM CDT this afternoon.\n* IMPACTS...Gusty winds could blow around unsecured objects.\nTree limbs could be blown down and a few power outages may\nresult."}]}

