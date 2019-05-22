package pronunciator

import (
	"bufio"
	"fmt"
	. "github.com/yhirose/go-peg"
	"log"
	"os"
	"regexp"
	"strings"
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


type ElementType int

const (
	World ElementType = iota
	Symbol
	Space
)

// en dictionary se van a cargar todas las palabras que contiene el archivo cmudict
var dictionary = make(map[string]string)

type Format int

const (
	Simplified Format = iota
	Ipa
)

type Element struct {
	elementType ElementType
	value string
}

type Pronunciator struct{
	parser *Parser
}


func NewPronunciator() (myparser *Pronunciator,err error){
	parser, err := NewParser(`
		# Grammar for simple calculator...
		DOCUMENT <- TEXT*
		TEXT <- WORD / SPACE / SYMBOL
		SYMBOL <- [.,?!#$%&/()='¿¡{};:] 
		WORD <- LETTER+
		SPACE <- ' '
		LETTER <- [a-zA-Z0-9]
	`)

	if err != nil{
		fmt.Println(err)
		return
	}

	g := parser.Grammar

	g["DOCUMENT"].Action = func(v *Values, d Any) (any Any, e error) {
		fmt.Println("running document")
		return v.Vs,nil
	}

	g["SYMBOL"].Action = func(v *Values, d Any) (any Any, e error) {
		return &Element{elementType:Symbol,value:v.S},nil
	}
	g["WORD"].Action = func(v *Values, d Any) (any Any, e error) {
		return &Element{elementType:World, value: v.S},nil
	}
	g["SPACE"].Action = func(v *Values, d Any) (any Any, e error) {
		return &Element{elementType:Space,value:" "},nil
	}

	myparser = &Pronunciator{parser: parser}

return
}

func (prs *Pronunciator) getElements(txt string)[]Element{
	values, _ := prs.parser.ParseAndGetValue(txt,nil)

	var results []Element
	for _,v := range values.([]Any){
		elemento := (v).(*Element)
		results = append(results,*elemento)
	}

	return results

}

// Pronounce : get some format: IPA or simple---and an english text and return
// the correct pronounciation
func (prs *Pronunciator) Pronounce(format Format, text string) string {
	var str strings.Builder
	elements := prs.getElements(text)
	for _, element := range elements {
		switch element.elementType {
		case Symbol:
			str.WriteString(element.value)
		case Space:
			str.WriteString(" ")
		case World:
			arphabetSound, wasFound := getArphabetPhonetic(strings.ToLower(element.value))
			if !wasFound {
				str.WriteString(element.value)
				continue
			}
			word := arphabetTo(format, arphabetSound)
			str.WriteString(word)
		}
	}
	return str.String()
}

//receive format: ipa or simplified and some cmu string extracted from the dictionary, return the same world but
//displaying the ipa or simplified pronunciation for instance WHAT soundCMU is W AH1 T
func arphabetTo(format Format, soundCMU string) string {
	requiredMap := simplifySounds
	if format == Ipa {
		requiredMap = arphabetToIPA
	}
	soundsCMU := strings.Split(soundCMU, " ") //convert "W AH1 T" in ["W","AH1","T"]
	soundWords := make([]string, len(soundsCMU))
	notDigitRg := regexp.MustCompile(`\d`)
	for _, sound := range soundsCMU {
		sound = notDigitRg.ReplaceAllString(sound, "") //remove the digits in the CMU sound
		if soundWord, ok := requiredMap[sound]; ok {
			soundWords = append(soundWords, soundWord) //for cmu AH return f.i ha
		} else {
			soundWords = append(soundWords, sound)
		}
	}
	return strings.Join(soundWords, "")
}
// check the cmu dictionary cmudict for the world and returns the cmu pronunciation
func getArphabetPhonetic(word string) (string, bool) {
	if val, ok := dictionary[word]; ok {
		return val, true
	}
	return word, false
}

func init(){
	file, err := os.Open("cmudict-0.7b.txt")
	if err != nil {
		log.Fatal("error opening file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	regxValidLine := regexp.MustCompile(`^[A-Z]+\s`)
	for scanner.Scan() {
		text := scanner.Text()
		if !regxValidLine.MatchString(text) {
			continue
		}

		words := strings.Split(text, " ")
		if len(words) > 1 {
			//for instance "A" cmu sound is AH0
			//so we take the A as key for our dictionary and join the rest of words AH0 and other sounds
			// hello would ve dictionary["hello"] "HH AH0 L OW1"
			dictionary[strings.ToLower(words[0])] = strings.TrimSpace(strings.Join(words[1:], " "))
		}
	}

	//fmt.Println("imprimiendo hola")
	//fmt.Println(dictionary["hello"])

	if err := scanner.Err(); err != nil {
		log.Fatal("error with scanner ", err)
	}
}


