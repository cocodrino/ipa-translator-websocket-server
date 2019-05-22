package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"

	"github.com/bregydoc/gtranslate"
)

var arphabetToIPA = map[string]string{
	"AA": "ɑ",
	"AE": "æ",
	"AH": "ʌ",
	"AO": "ɔ",
	"AW": "aʊ",
	"AY": "aɪ",
	"B":  "b",
	"CH": "tʃ",
	"D":  "d",
	"DH": "ð",
	"EH": "ɛ",
	"ER": "ɝ",
	"EY": "eɪ",
	"F":  "ɾ",
	"G":  "ɡ",
	"HH": "h",
	"IH": "ɪ",
	"IY": "i",
	"JH": "dʒ",
	"K":  "k",
	"L":  "l",
	"M":  "m",
	"N":  "n",
	"NG": "ŋ",
	"OW": "oʊ",
	"OY": "ɔɪ",
	"P":  "p",
	"R":  "ɹ",
	"S":  "s",
	"SH": "ʃ",
	"T":  "t",
	"TH": "θ",
	"UW": "u",
	"UH": "ʊ",
	"V":  "v",
	"W":  "w",
	"Y":  "j",
	"Z":  "z",
	"ZH": "ʒ",
}

var simplifySounds = map[string]string{
	"AA": "A",
	"AE": "a",
	"AH": "ʌ",
	"AO": "ao",
	"AW": "Au",
	"AY": "ai",
	"B":  "b",
	"CH": "ch",
	"D":  "d",
	"DH": "D",
	"EH": "E",
	"ER": "or",
	"EY": "ei",
	"F":  "f",
	"G":  "ɡ",
	"HH": "J",
	"IH": "i",
	"IY": "ii",
	"JH": "y",
	"K":  "k",
	"L":  "l",
	"M":  "m",
	"N":  "n",
	"NG": "n",
	"OW": "Ou",
	"OY": "Oy",
	"P":  "p",
	"R":  "r",
	"S":  "s",
	"SH": "ʃ",
	"T":  "t",
	"TH": "tt",
	"UW": "U",
	"UH": "Uu",
	"V":  "v",
	"W":  "w",
	"Y":  "y",
	"Z":  "z",
	"ZH": "zz",
}

type WebsocketResponse struct {
	Original   string `json:"original"`
	Ipa        string `json:"ipa"`
	Simplified string `json:"simplified"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


type Format int

const (
	Simplified Format = iota
	Ipa
)

var pronunciator *Pronunciator



//run go run $(ls -1 *.go | grep -v _test.go)
func main() {
	pronunciator, _ = NewPronunciator()
	//dictionary := make(map[string]string)

	//fmt.Println(Pronounce(Ipa, "what! are you doing today let me know now."))
	http.HandleFunc("/translate", func(w http.ResponseWriter, request *http.Request) {
		conn, _ := upgrader.Upgrade(w, request, nil)

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			fmt.Println(string(msg))
			//resp := Pronounce(Ipa, string(msg))
			translatedText,err := gtranslate.TranslateWithFromTo(string(msg),gtranslate.FromTo{From:"auto",To:"en"})
			if err != nil{
				fmt.Println("error translating ",err)
			}

			jsonObj := &WebsocketResponse{Original: translatedText, Ipa:  pronunciator.Pronounce(Ipa, translatedText), Simplified: pronunciator.Pronounce(Simplified, translatedText)}
			fmt.Println(jsonObj)
			resp, err := json.Marshal(jsonObj)
			fmt.Println("json", string(resp))

			if err != nil {
				fmt.Println("error json ", err)
			}

			if err = conn.WriteMessage(msgType, resp); err != nil {

			}
		}
	})

	fmt.Println("server running in port :8080, you can start to translating using the endpoint /translate")
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Println("error running the server >\n", err)
	}

}
