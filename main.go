package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

type agencyResponseData struct {
	Status int
	Result agency
}

type agency struct {
	Id string
	Name string
	ServiceProviderCode string
	HostId string
	LogLevel int
	Enabled bool
	DisplayName string
	IsForDemo bool
	State string
	Country string
	IsDefaultAppActive bool
	DatabaseType int
}

type agencyConf struct {
	Agencies []map[string]string `yaml:agencies`
}

var filePath *string
var spc *string
var target *string

func main(){

	filePath := flag.String("path", "C:\\Dev\\Go\\src\\github.com\\rydavidson\\AZMigrate\\agencies.yaml", "Path to the agencies.yml file")
	spc := flag.String("agency", "", "Agency to move from current host to target")
	target := flag.String("target", "", "Target hostid to move the agency to")
	flag.Parse()

	var aConf agencyConf

	var wg sync.WaitGroup

	if *target == "" && *spc == "" {
		aConf.readAgencies(*filePath)
		for _, v0 := range aConf.Agencies {
			for k1, v1 := range v0 {
				if k1 == "agency_name" {
					log.Printf("Prcessing agency %v", v1)
					agencyData := agencyResponseData{}
					wg.Add(1)
					go updateAgency(v1,&agencyData, &wg)
				}
			}
		}
	}

	wg.Wait()
	log.Print("Done processing all agencies")
}

func processAgency(a string, targetRes *agencyResponseData, wg *sync.WaitGroup){

	defer wg.Done()

	subsystemAccessKey := os.Getenv("CONSTRUCT_ACCESS_KEY")
	url := "https://admin.accela.com/apis/v4/agencies/" + a
	client := &http.Client{}
	req,_ := http.NewRequest("GET",url,nil)
	req.Header.Set("x-accela-subsystem-accesskey",subsystemAccessKey)

	res, err := client.Do(req)

	if err != nil {
		log.Printf("Error getting agency data: $%v", err)
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(targetRes)

	log.Print(targetRes.Result)
}

func updateAgency(a string, targetRes *agencyResponseData, wg *sync.WaitGroup){

	defer wg.Done()

	subsystemAccessKey := os.Getenv("CONSTRUCT_ACCESS_KEY")
	url := "https://admin.accela.com/apis/v4/agencies/" + a
	client := &http.Client{}
	req,_ := http.NewRequest("GET",url,nil)
	req.Header.Set("x-accela-subsystem-accesskey",subsystemAccessKey)

	res, err := client.Do(req)

	if err != nil {
		log.Printf("Error getting agency data: $%v", err)
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(targetRes)

	agencyData := targetRes.Result

	agencyData.updateForAzureMT()

	//log.Print(targetRes.Result)

}

func (agencyData *agency) updateForAzureMT(){

	subsystemAccessKey := os.Getenv("CONSTRUCT_ACCESS_KEY")

	agencyData.DatabaseType = 2
	agencyData.HostId = "52f80df3-26db-462a-8c4c-125fdb29ce97"

	url := "https://admin.accela.com/apis/v4/agencies/" + agencyData.Id
	client := &http.Client{}
	body, jsonErr := json.Marshal(agencyData)
	if jsonErr != nil {
		log.Fatalf("Error marshaling json: #%v",jsonErr)
	}
	req,_ := http.NewRequest("PUT",url,bytes.NewBuffer(body))
	req.Header.Set("x-accela-subsystem-accesskey",subsystemAccessKey)
	req.Header.Set("Content-Type","application/json")

	log.Println(string(body))

	res, err := client.Do(req)

	if err != nil {
		log.Printf("Error updating agency data: $%v", err)
	}

	defer res.Body.Close()

	resbody,err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(resbody))

}


func (conf *agencyConf) readAgencies(filePath string) agencyConf{

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error getting yaml file: #%v", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil{
		log.Fatalf("Unmarshal: %v", err)
	}

	//log.Print(conf)

	return *conf

}
