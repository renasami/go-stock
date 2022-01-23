const WebSocket = require("ws");
const ws = new WebSocket("wss://api.coin.z.com/ws/public/v1");
console.log("start program")
ws.on("open", () => {
    console.log("open")
    const message = JSON.stringify(
    {
        "command": "subscribe",
        "channel": "ticker",
        "symbol": "BTC_JPY"
    });
    ws.send(message);
});

ws.on("message", (data) => {
    console.log("WebSocket message: ", data);
});