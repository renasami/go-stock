package main

import (
	"app/logging"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //オプションを解析します。デフォルトでは解析しません。
	fmt.Println(r.Form) //このデータはサーバのプリント情報に出力されます。
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //ここでwに入るものがクライアントに出力されます。
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func autoWebSocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		n := 0
		for {
			any := getExchangeRate(c)
			websocket.JSON.Send(ws, any)
			n += 1
			c.Logger().Info("time")
			time.Sleep(2 * time.Second)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

//外為オンラインから来たデータをエンコード
func getExchangeRate(c echo.Context) interface{} {
	// resp, err := http.Get("http://fx.mybluemix.net/")
	resp, err := http.Get("https://www.gaitameonline.com/rateaj/getrate")
	if err != nil {
		c.Logger().Error(err)
		return nil
	}
	defer func() {
		defer resp.Body.Close()
		io.Copy(ioutil.Discard, resp.Body)
	}()
	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		c.Logger().Error("Some Status Error")
		return nil
	}
	var t interface{}
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		c.Logger().Error(err)
		return nil
	}
	return t

}

func SendSubscribeRequest(c echo.Context) error {
	fmt.Print("send subscribe request")
	wsUrl := "wss://api.coin.z.com/ws/public/v1"
	origin := "https://api.coin.z.com"
	btc_message := (`{
		"command": "subscribe",
		"channel": "ticker",
		"symbol": "ETH_JPY"
	}`)
	eth_message := (`{
		"command": "subscribe",
		"channel": "ticker",
		"symbol": "BTC_JPY"
	}`)
	ws, _ := websocket.Dial(wsUrl, "", origin)
	websocket.Message.Send(ws, btc_message)
	fmt.Println("sended1")
	time.Sleep(time.Second * 2)
	websocket.Message.Send(ws, eth_message)
	fmt.Println("sended2")
	return nil
}

func GetRealtimeBtcRate(c echo.Context) error {
	wsUrl := "wss://api.coin.z.com/ws/public/v1"
	origin := "https://api.coin.z.com"
	ws, _ := websocket.Dial(wsUrl, "", origin)
	var receiveMsg string
	websocket.Handler(func(wss *websocket.Conn) {
		for {
			websocket.Message.Receive(ws, &receiveMsg)
			var buf bytes.Buffer
			json.Indent(&buf, []byte(receiveMsg), "", "  ")
			websocket.JSON.Send(wss, buf.String())
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

//maybe i can perform it returning ws construc
func GetRealtimeEthRate(c echo.Context) error {
	wsUrl := "wss://api.coin.z.com/ws/public/v1"
	origin := "https://api.coin.z.com"
	var receiveMsg string
	ws, _ := websocket.Dial(wsUrl, "", origin)
	websocket.Handler(func(wss *websocket.Conn) {
		for {
			websocket.Message.Receive(ws, &receiveMsg)
			fmt.Println(receiveMsg)
			var buf bytes.Buffer
			json.Indent(&buf, []byte(receiveMsg), "", "  ")
			websocket.JSON.Send(wss, buf.String())
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
func main() {
	// conf := config.Config
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/", "public")
	e.Use(middleware.CORS())
	e.GET("/ws", autoWebSocket)
	e.GET("/realtime_btc_rate", GetRealtimeBtcRate)
	e.GET("/realtime_eth_rate", GetRealtimeEthRate)
	e.GET("/first", SendSubscribeRequest)
	e.Logger.Fatal(e.Start(":80"))

	r := mux.NewRouter()
	logger := logging.Logger()
	http.HandleFunc("/hello", sayhelloName)
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	loggerRouter := logging.Middleware(logger)(r)
	err := http.ListenAndServe(":80", loggerRouter)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
