package main

import (
	"fmt"
	"log"
	"net/http"
	"src/app/logging"
	"strings"

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

		for {
			// Client からのメッセージを読み込む
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}

			// Client からのメッセージを元に返すメッセージを作成し送信する
			err := websocket.Message.Send(ws, fmt.Sprintf("Server: \"%s\" received!", msg))
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func main() {
	// conf := config.Config
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/", "public")
	e.GET("/ws", handleWebSocket)
	e.Logger.Fatal(e.Start(":80"))

	r := mux.NewRouter()
	logger := logging.Logger()
	http.HandleFunc("/hello", sayhelloName)
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	loggerRouter := logging.Middleware(logger)(r)
	err := http.ListenAndServe(fmt.Sprintf(":80"), loggerRouter)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
