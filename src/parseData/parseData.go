package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

//type Create struct {
//	XMLName xml.Name `xml:"create"`
//	Accoun  Account  `xml:"account"`
//	Symbo   Symbol   `xml:"Symbol"`
//}

type Account struct {
	XMLName xml.Name `xml:"account"`
	ID      int      `xml:"id,attr"`
	Balance float64  `xml:"balance,attr"`
}

type Symbol struct {
	XMLName xml.Name `xml:"symbol"`
	Sym     string   `xml:"sym,attr"`
	Posit   []Posit  `xml:"account"`
}

type Posit struct {
	ID       int     `xml:"id,attr"`
	Position float64 `xml:",chardata"`
}

type Order struct {
	XMLName xml.Name `xml:"order"`
	Sym     string   `xml:"sym,attr"`
	Amount  int      `xml:"amount,attr"`
	Limit   float64  `xml:"limit,attr"`
}

type Query struct {
	XMLName xml.Name `xml:"query"`
	QueryID int      `xml:"id,attr"`
}

type Cancel struct {
	XMLName  xml.Name `xml:"cancel"`
	CancelID xml.Name `xml:"id,attr"`
}

type Transactions struct {
	XMLName xml.Name `xml:"transactions"`
	ID      int      `xml:"account,attr"`
	//Contents []Mixed    `xml:",any"`
}

type Mixed struct {
	Type  string
	Value interface{}
}

func main() {
	data := `
		    <?xml version="1.0" encoding="UTF-8"?>
			<transactions></transactions>
			<create>
			  <account id="123456" balance="1000"/>
			  <symbol sym="SPY">
			    <account id="123456">100000</account>
			  </symbol>
			  <account id="1234567" balance="1000"/>
			</create>
			<transactions account="12345">
		    `
	input := strings.NewReader(data)
	decoder := xml.NewDecoder(input)
	var result, inElement string
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "create" {
				var substr []byte
				for {
					tc, _ := decoder.Token()
					if tc == nil {
						break
					}
					switch sc := tc.(type) {
					case xml.StartElement:
						if sc.Name.Local == "account" {
							var account Account
							decoder.DecodeElement(&account, &sc)
							//fmt.Println(account.ID)
							substr = []byte(string(substr) + "create account" + "\n")
						}
						if sc.Name.Local == "symbol" {
							var symbol Symbol
							decoder.DecodeElement(&symbol, &sc)
							//fmt.Println(symbol.Posit[0].ID, "  ", symbol.Posit[0].Position)
							substr = []byte(string(substr) + "handle symbol" + "\n")
						}
					}
				}
				substr = []byte(xml.Header + "<results>" + " \n" + string(substr) + "</results>" + "\n")
				result = string(substr)
			}
			if inElement == "transactions" {
				var T Transactions
				for _, attr := range se.Attr {
					//attrName := attr.Name.Local
					attrValue := attr.Value
					//fmt.Println(attrName, " ", attrValue)
					ID, _ := strconv.Atoi(attrValue)
					//if err != nil {
					//	return result
					//}
					T.ID = ID
					fmt.Println("T.ID:", T.ID)
				}
				//fmt.Println(T.ID)
				//}
			}
		default:
		}
	}
	fmt.Printf("%s", result)
}
