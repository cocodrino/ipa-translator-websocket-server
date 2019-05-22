package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	pronunciator2 "traductorWS2/pronunciator"

	"github.com/bregydoc/gtranslate"
)

type WebsocketResponse struct {
	Original   string `json:"original"`
	Ipa        string `json:"ipa"`
	Simplified string `json:"simplified"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var pronunciator *pronunciator2.Pronunciator

func init() {
	pronunciator, _ = pronunciator2.NewPronunciator()
}

//run go run $(ls -1 *.go | grep -v _test.go)

//fmt.Println(Pronounce(Ipa, "what! are you doing today let me know now."))
func Handler(w http.ResponseWriter, request *http.Request) {
	conn, _ := upgrader.Upgrade(w, request, nil)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		fmt.Println(string(msg))
		//resp := Pronounce(Ipa, string(msg))
		translatedText, err := gtranslate.TranslateWithFromTo(string(msg), gtranslate.FromTo{From: "auto", To: "en"})
		if err != nil {
			fmt.Println("error translating ", err)
		}

		jsonObj := &WebsocketResponse{Original: translatedText, Ipa: pronunciator.Pronounce(pronunciator2.Ipa, translatedText), Simplified: pronunciator.Pronounce(pronunciator2.Simplified, translatedText)}
		fmt.Println(jsonObj)
		resp, err := json.Marshal(jsonObj)
		fmt.Println("json", string(resp))

		if err != nil {
			fmt.Println("error json ", err)
		}

		if err = conn.WriteMessage(msgType, resp); err != nil {

		}
	}
}
