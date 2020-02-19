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
	Id                  string
	Name                string
	ServiceProviderCode string
	HostId              string
	LogLevel            int
	Enabled             bool
	DisplayName         string
	IsForDemo           bool
	State               string
	Country             string
	IsDefaultAppActive  bool
	DatabaseType        int
}

type agencyConf struct {
	Agencies []map[string]string `yaml:agencies`
}

var filePath *string
var spc *string
var target *string
var enable *bool
var disable *bool
var set_enabled int  // default is 0

func main() {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get working directory #%v", err)
	}

	filePath := flag.String("path", wd+"\\agencies.yml", "Path to the agencies.yml file")
	spc := flag.String("agency", "", "Agency to move from current host to target")
	target := flag.String("target", "", "Target hostid to move the agency to")
	enable := flag.Bool("enable", false, "flag to enable")
	disable := flag.Bool("disable", false, "flag to disable")
	
	flag.Parse()

	var aConf agencyConf

	var wg sync.WaitGroup

	if *enable && *disable {
		log.Println("can't ask to both enable and disable")
	} else {
		if *enable {
			set_enabled = 2
		}
		if *disable {
			set_enabled = 1
		}

		if *target == "" && *spc == "" {
			aConf.readAgencies(*filePath)
			// Get the values from the agencies collection
			for _, v0 := range aConf.Agencies {
				// Get the elements in each agency entry
				for k1, v1 := range v0 {
					// Check if the element is the agency_name and process if so, ignore the rest
					//TODO Unmarshal the yml so agency_name can be checked directly
					if k1 == "agency_name" {
						log.Printf("Prcessing agency %v", v1)
						agencyData := agencyResponseData{}
						wg.Add(1)
						go updateAgency(v1, set_enabled, &agencyData, &wg)
					}
				}
			}
		} else {
			log.Println("target and agency flags aren't implemented yet, sorry")
		}
	}

	//TODO Add logic for target and agency flags
	wg.Wait()
	log.Print("Done processing all agencies")
}

func updateAgency(a string, enabledflag int, targetRes *agencyResponseData, wg *sync.WaitGroup) {

	defer wg.Done()
	//TODO Make this global
	subsystemAccessKey := os.Getenv("CONSTRUCT_ACCESS_KEY")
	url := "https://admin.accela.com/apis/v4/agencies/" + a
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-accela-subsystem-accesskey", subsystemAccessKey)

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error getting agency data: $%v", err)
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(targetRes)
	agencyData := targetRes.Result

	//TODO more logic to handle different use-cases - here, either enable/disable or re-home
	if enabledflag > 0 {
		//agencyData.Enabled = enabledflag - 1
		agencyData.Enabled = !((enabledflag- 1) == 0)
		log.Printf("setting enabled flag to: %v", agencyData.Enabled)
	} else {
		agencyData.DatabaseType = 2
		//TODO Accept as an arg instead of hardcoded
		agencyData.HostId = "52f80df3-26db-462a-8c4c-125fdb29ce97"
		log.Printf("changing to SQLServer and Azure-MT-Host")
	}

	agencyData.updateForAzureMT()
}

// Updates the called on agencyData and makes an api call to submit the changes to Construct
func (agencyData *agency) updateForAzureMT() {
	//TODO Make this global
	subsystemAccessKey := os.Getenv("CONSTRUCT_ACCESS_KEY")

	url := "https://admin.accela.com/apis/v4/agencies/" + agencyData.Id
	client := &http.Client{}
	body, jsonErr := json.Marshal(agencyData)
	if jsonErr != nil {
		log.Fatalf("Error marshaling json: #%v", jsonErr)
	}
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	req.Header.Set("x-accela-subsystem-accesskey", subsystemAccessKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error updating agency data: $%v", err)
	}
	defer res.Body.Close()

	resbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(resbody))
}

func (conf *agencyConf) readAgencies(filePath string) agencyConf {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error getting yaml file: #%v", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return *conf
}
