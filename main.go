package main

import "bytes"
import "flag"
import "fmt"
import "io/ioutil"
import "log"
import "net/http"
import "os"
import "encoding/xml"
import "time"

// download weather data for an airport and determine flight rules

type stationid string

func str2stationid(id string) (stationid, error) {
	return stationid(id), error(nil)
}

/*
type ADDSdata struct {
	XMLName  xml.Name `xml:"response"`
	Key      string   `xml:"key,attr"`
	Value    string   `xml:"value,attr"`
}
*/

// Note feature request:
// https://code.google.com/p/go/issues/detail?id=3688
type ADDSdata_source struct {
	XMLName  xml.Name
	Name	string	`xml:"name,attr"`
}

type ADDSMETAR struct {
	XMLName  xml.Name
	raw_text string
	station_id string
	observation_time string
	latitude string
	longitude string
	temp_c int
	dewpoint_c int
	wind_dir_degrees int
	wind_speed_kt int
	visibility_statute_mi float32
	altim_in_hg float32
	quality_control_flags string
	//sky_condition sky_cover="OVC" cloud_base_ft_agl="1900" />
	flight_category string
	metar_type string
	elevation_m float32
}

type ADDSdata struct {
	XMLName  xml.Name `xml:"data"`
	Num_results    int   `xml:"num_results,attr"`
	Metars []ADDSMETAR `xml:"METARS"`
}


type ADDSresponse struct {
	XMLName  xml.Name `xml:"response"`
	Request_index      int   `xml:"request_index"`
	Data_source      ADDSdata_source   `xml:"data_source"`
	Data      ADDSdata   `xml:"data"`
}


func fetchADDS(sym stationid) ([]byte, error) {

	//url := "http://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=metars&requestType=retrieve&format=xml&stationString=KSQL&hoursBeforeNow=1"
	url := fmt.Sprintf("http://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=metars&requestType=retrieve&format=xml&stationString=%v&hoursBeforeNow=1", sym)

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

func writeADDS(sym stationid, body []byte) (error) {
	ssym := string(sym)
	err := os.MkdirAll(ssym, 0777)
	if err != nil { return err }

	err = os.Chdir(ssym)
	if err != nil { return err }

	fname := fmt.Sprintf("%v.%v.adds", sym, time.Now().Format(time.RFC3339))

	err = ioutil.WriteFile(fname, body, 0666)
	if err != nil { return err }

	return nil
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

func usage() {
	fmt.Println("Subcommands are:")
	fmt.Println("  fetch <code>     Pull METAR for airport with code")
	fmt.Println("  parse <code>     Parse downloaded data for airport code")
}


func main() {
	fmt.Println("Starting")

	flag.Parse();

	cmd := flag.Args()

	fmt.Println(cmd)

	if( len(cmd) < 2 ) {
		usage()
		os.Exit(1)
	}

	var e error

	stationid, e := str2stationid(cmd[1])

	if( e != nil ) { log.Fatal(e); os.Exit(1) }

	if( len(cmd) > 0 && cmd[0] == "fetch" ) {
		fmt.Printf("Fetching for airport %v\n", stationid)
		body, e := fetchADDS(stationid)
		if( e != nil ) {
			log.Fatal(e)
			os.Exit(1)
		}
		writeADDS(stationid, body)
		os.Exit(0)
	}

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



