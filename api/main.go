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
		msg := ""
		err = websocket.Message.Receive(ws, &msg)
		if err != nil {
			c.Logger().Error(err)
		}

		// Client からのメッセージを元に返すメッセージを作成し送信する
		er := websocket.Message.Send(ws, fmt.Sprintf("Server: \"%s\" received!", msg))
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

type ExchangeRate struct {
	Status bool `json:"status"`
	Result struct {
		Datetime string `json:"datetime"`
		Rate     struct {
			Usdjpy float64 `json:"USDJPY"`
			Eurjpy float64 `json:"EURJPY"`
			Eurusd float64 `json:"EURUSD"`
			Audjpy float64 `json:"AUDJPY"`
			Gbpjpy float64 `json:"GBPJPY"`
			Nzdjpy float64 `json:"NZDJPY"`
			Cadjpy float64 `json:"CADJPY"`
			Chfjpy float64 `json:"CHFJPY"`
			Hkdjpy float64 `json:"HKDJPY"`
			Gbpusd float64 `json:"GBPUSD"`
			Usdchf float64 `json:"USDCHF"`
			Zarjpy float64 `json:"ZARJPY"`
			Audusd float64 `json:"AUDUSD"`
			Nzdusd float64 `json:"NZDUSD"`
			Euraud float64 `json:"EURAUD"`
			Tryjpy float64 `json:"TRYJPY"`
			Cnhjpy float64 `json:"CNHJPY"`
			Nokjpy float64 `json:"NOKJPY"`
			Sekjpy float64 `json:"SEKJPY"`
			Mxnjpy float64 `json:"MXNJPY"`
		} `json:"rate"`
	} `json:"result"`
}

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
	e.GET("/wss", handleWebSocket)
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
