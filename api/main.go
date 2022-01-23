package main

import (
	"app/logging"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"bytes"

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

func handleWebSocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		// 初回のメッセージを送信
		err := websocket.Message.Send(ws, "Server: Hello, Client!")
		if err != nil {
			c.Logger().Error(err)
		}

		// Client からのメッセージを読み込む
		var msg string
		err = websocket.Message.Receive(ws, &msg)
		if err != nil {
			c.Logger().Error(err)
		}

		// Client からのメッセージを元に返すメッセージを作成し送信する
		er := websocket.Message.Send(ws, fmt.Sprintf("eceived!"))
		if er != nil {
			c.Logger().Error(er)
		}
		for {
			websocket.Message.Send(ws, fmt.Sprintln("hello"))
			c.Logger().Info("time")
			time.Sleep(3 * time.Second)
		}

	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func getRealtimeETHRate(c echo.Context, coin string) interface{} {
	wsUrl := "wss://api.coin.z.com/ws/public/v1"
	baseUrl := "https://api.coin.z.com"
	sendMsg := (`{
        "command": "subscribe",
        "channel": "ticker",
        "symbol": "ETH"
    }`)

	var rate string

    ws, _ := websocket.Dial(wsUrl, "", baseUrl)
    websocket.Message.Send(ws, sendMsg)
	for {
		err := websocket.Message.Receive(ws, &rate)
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Print(rate)
	}
}
func getRealtimeBTCRate(c echo.Context) error {
	wsUrl := "wss://api.coin.z.com/ws/public/v1"
	baseUrl := "https://api.coin.z.com"
	sendMsg := (`{
        "command": "subscribe",
        "channel": "ticker",
        "symbol": "BTC"
    }`)

	var rate string

    ws, _ := websocket.Dial(wsUrl, "", baseUrl)
    websocket.Message.Send(ws, sendMsg)
	websocket.Handler(func(wss *websocket.Conn) {
	for {
		err := websocket.Message.Receive(wss,&rate)
		var buf bytes.Buffer
        json.Indent(&buf, []byte(rate), "", "  ")

		if err != nil {
			fmt.Println("ERRRRRRRRR")
			fmt.Println(err)
		}
		websocket.JSON.Send(wss,buf.String())
	}

	}).ServeHTTP(c.Response(), c.Request())
	return nil
	// for {
	// 	err := websocket.Message.Receive(ws, &rate)
	// 	if err != nil {
	// 		c.Logger().Error(err)
	// 	}
	// 	fmt.Print(rate)
	// }
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

func main() {
	// conf := config.Config
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/", "public")
	e.GET("/ws", autoWebSocket)
	e.GET("/wss", getRealtimeBTCRate)
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
