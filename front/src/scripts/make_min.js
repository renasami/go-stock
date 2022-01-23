var request = require('request');

var endPoint = 'https://api.coin.z.com/public';
var path     = '/v1/klines?symbol=BTC&interval=1min&date=20220122';

request(endPoint + path, function (err, response, payload) {
    console.log(JSON.stringify(JSON.parse(payload), null, 2));
});
