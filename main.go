package main

import "bytes"
import "fmt"
import "io/ioutil"
import "log"
import "net/http"
import "os"
import "encoding/xml"

// download weather data for an airport and determine flight rules

type stationid string


type ADDSdata struct {
	XMLName  xml.Name `xml:"response"`
	Key      string   `xml:"key,attr"`
	Value    string   `xml:"value,attr"`
}

type ADDSresponse struct {
	XMLName  xml.Name `xml:"response"`
	Request_index      int   `xml:"request_index"`
	Data_source      string   `xml:"data_source>name,attr"`
}


func fetchADDS(sym stationid) ([]byte, error) {

	//url := "http://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=metars&requestType=retrieve&format=xml&stationString=KSQL&hoursBeforeNow=1"
	url := "http://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=metars&requestType=retrieve&format=xml&stationString=KSQL&hoursBeforeNow=1"

	resp, err := http.Get(url)
	if err != nil { return nil, err }

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	} else {
		fmt.Println(string(body))
	}

	return body, nil
}

func ADDS2xml(body []byte) (*ADDSresponse, error) {
	breader := bytes.NewReader(body)
	adds := &ADDSresponse{}
	decoder := xml.NewDecoder(breader)

	err := decoder.Decode(adds)
	if( err != nil ) {
		return nil, err
	}


	return adds, nil
}


func main() {
	fmt.Println("Starting")

	//_, e := fetchADDS("KSQL")

	var e error

	addstxt, e := ioutil.ReadFile("sample_adds.xml")

	if( e != nil ) {
		log.Fatal(e)
		os.Exit(1)
	}

	fmt.Println(string(addstxt[:]))

	adds, xmle := ADDS2xml(addstxt)
	if( xmle != nil ) {
		log.Fatal(xmle)
		os.Exit(1)
	}

	log.Print("Parsed xml")

	fmt.Println(adds.XMLName)
	fmt.Println(adds.Request_index)
	fmt.Println(adds.Data_source)
}



