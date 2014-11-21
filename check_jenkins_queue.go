package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"
import  "flag"
import "os"

type Queue struct {

    Items []Job

}

type Job struct {
    Stuck bool
    Buildable bool
    InQueueSince int64
    Url string
    Blocked bool
    Why string
    Id int32
}

func get_content(url string, verbose bool) (items int, err error) {

    defer func() {
        if err := recover(); err != nil {
            return
        }
    }()

    res, err := http.Get(url)
    if err != nil {
        return
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return
    }

    var data Queue
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

    for i, job := range data.Items {
        items += 1
        if verbose {

            fmt.Fprintf(os.Stderr, "%d: %s %s\n", i, job.Url, job.Why)
            
        }
    }

    return
}

func main() {

    status := "OK"
    rv := 0

    verbose := flag.Bool("v", false, "verbose output")
    warn := flag.Int("w", 10, "warning level for job queue depth")
    crit := flag.Int("c", 20, "critical level for job queue depth")
    host := flag.String("h", "http://ci.walmartlabs.com/jenkins", "base url for jenkins  like http://ci.walmartlabs.com/jenkins")

    url := *host + "/queue/api/json"

    flag.Parse()

    if len(flag.Args()) > 0 {

        flag.Usage()
        os.Exit(3)

    }

    defer func() {
        if err := recover(); err != nil {
            fmt.Println("Unknown: ", err)
            os.Exit(3)
        }
    }()

    if *verbose {

        fmt.Printf("checking queues on:%s:warning:%d:critical:%d\n", url, *warn, *crit)
        
    }
    jobs, err := get_content(url, *verbose)
    if err != nil {

        fmt.Printf("Unknown: %T %s %#v\n",err, err, err)

        os.Exit(3)

    }

    if jobs >= *crit {
        status = "Critical"
        rv = 1
    } else if jobs >= *warn {
        status = "Warning"
        rv = 2
    } 

    fmt.Printf("%s: %d\n", status, jobs)
    os.Exit(rv)

}
