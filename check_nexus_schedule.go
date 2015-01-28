package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "flag"
import "os"
import "os/user"
import 	"code.google.com/p/go-netrc/netrc"

type SchedulesList struct {
	Data []SchedulesListItem
}

type SchedulesListItem struct {
	ResourceURI    string
	Id             string
	Status         string
	ReadableStatus string
}

func get_content(url string, username string, password string, verbose bool) (items int, err error) {

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")

	if username != "" && password != "" {

		req.SetBasicAuth(username, password)

	}
	
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	// fmt.Printf("get_content:body:%s:\n", body)

	var data SchedulesList
	err = json.Unmarshal(body, &data)
	if err != nil {

		if verbose {

			switch v := err.(type) {
			case *json.SyntaxError:
				fmt.Fprintf(os.Stderr, string(body[v.Offset-40:v.Offset]))
			}

		}
		return
	}

	for i, item := range data.Data {

		if item.ReadableStatus == "BLOCKED_AUTO" {
			items += 1
		}

		if verbose {

			fmt.Fprintf(os.Stderr, "%d: %s %s %s %s\n", i, item.Id, item.ResourceURI, item.Status, item.ReadableStatus)

		}
	}

	return
}

func main() {

	status := "OK"
	rv := 0
	name := "NexusSchedules"

	verbose := flag.Bool("v", false, "Verbose output")
	warn := flag.Int("w", 10, "Warning level for blocked scheduled items")
	crit := flag.Int("c", 20, "Critical level for blocked scheduled items")
	host := flag.String("h", "http://gec-maven-nexus.walmart.com/nexus", "Base url for jenkins api like http://gec-maven-nexus.walmart.com/nexus")

	flag.Parse()

	if len(flag.Args()) > 0 {

		flag.Usage()
		os.Exit(3)

	}

    _url := *host + "/service/local/schedules"
    url := &_url

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(name+" Unknown: ", err)
			os.Exit(3)
		}
	}()

	if *verbose {

		fmt.Printf("checking schedules on:url:%s:warning:%d:critical:%d\n", *url, *warn, *crit)

	}


	self, err := user.Current()
	if err != nil {
		fmt.Printf("Get Current User:err:%t:\n", err)
		os.Exit(3)
	}

	// fmt.Printf("--- ooEnvMerged.Site:\n%s\n\n", *ooEnvMerged.Site)
	creds, err := netrc.FindMachine(self.HomeDir+"/.netrc", *host)
	if err != nil {
		fmt.Printf("NetRC Error:Couldn't find the key from the Login field:netrcfile:%s:machine:%s:\n", self.HomeDir+"/.netrc", *host)
		os.Exit(3)
	}

	jobs, err := get_content(*url, creds.Login, creds.Password, *verbose)
	if err != nil {

		fmt.Printf("%s Unknown: %T %s %#v\n", name, err, err, err)

		os.Exit(3)

	}

	if jobs >= *crit {
		status = "Critical"
		rv = 1
	} else if jobs >= *warn {
		status = "Warning"
		rv = 2
	}

	fmt.Printf("%s %s: %d\n", name, status, jobs)
	os.Exit(rv)

}
