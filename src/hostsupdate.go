package main

import (
    "io/ioutil"
    "log"
    "net/http"
    "strings"
    "os"
    "time"
    "github.com/eyedeekay/gosam"
)

type hostUpdater struct{
    retries int
    tryFirst string
    parentList string
    samBridgeClient *goSam.Client
    samBridgeErrors error

    hostList [][][]string
    hostfile string
}

func (updater *hostUpdater) parseCsv(s string) []string{
    hosts := strings.Split(s, ",")
    return hosts
}

func (updater *hostUpdater) parseKvp(s string) [][]string{
    hosts := &[][]string{}
    for _, host := range strings.Split(s, "\n"){
        kv := strings.SplitN(host, "=", 2)
        if len(kv) == 2 {
            *hosts = append(*hosts, kv)
        }
    }
    return *hosts
}

func (updater *hostUpdater) writeHostList() error {
    f, err := os.Create(updater.hostfile)
    updater.Fatal(err, "File I/O errors")
    defer f.Close()
    for _, s := range updater.hostList {
        for _, t := range s {
            if len(t) == 2 {
                line := t[0] +"="+ t[1] +"\n"
                f.WriteString(line)
            }
        }
        f.WriteString("\n")
    }
    return err
}

func (updater *hostUpdater) hostUpdate(){
    t := updater.retries
    var hostList string
    for t >= 1 {
        updater.Log("Getting updates from new host providers.")
        for _, u := range updater.parseCsv(updater.tryFirst) {
            if done, h := updater.get(u); done {
                hostList += h + "\n"
                updater.hostList = append(updater.hostList, updater.parseKvp(hostList))
                updater.writeHostList()
                t = 0
                break
            }
            time.Sleep(time.Duration(1) * time.Second)
        }
        t--
    }
    t = updater.retries
    for t >= 1 {
        updater.Log("Getting updates from upstream host providers.")
        for _, u := range updater.parseCsv(updater.parentList) {
            if done, h := updater.get(u); done {
                hostList += h + "\n"
                updater.hostList = append(updater.hostList, updater.parseKvp(hostList))
                updater.writeHostList()
                t = 0
                break
            }
            time.Sleep(time.Duration(1) * time.Second)
        }
        t--
    }
    updater.Log("Updates complete.")
}

func (updater *hostUpdater) get(s string) (bool, string){
    tr := &http.Transport{
		Dial: updater.samBridgeClient.Dial,
    }

    updater.Log("Fetching updates from: " + s)

   	client := &http.Client{Transport: tr}

    resp, err := client.Get(s)

    r := ""

    t := false

    if ! updater.Warn(err, "Updater Error: " + s + " ") {
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if ! updater.Warn(err, "Response error: ") {
            r += string(body)
            t = true
        }
    }

    return t, r
}

func (updater *hostUpdater) getHosts() [][][]string{
    return updater.hostList
}

func (updater *hostUpdater) Log(s string){
    if loglevel > 2 {
        log.Println(s)
    }
}

func (updater *hostUpdater) Warn(err error, s string) bool {
    if loglevel > 1 {
        if err != nil {
            log.Println(s, err)
            return true
        }
        return false
    }
    return false
}

func (updater *hostUpdater) Fatal(err error,s string) bool{
    if loglevel > 0 {
        if err != nil {
            log.Println(s, err)
            return true
        }
        return false
    }
    return false
}

func newHostUpdater(samhost string, samport string, retries int, upstream string, parent string, hostfile string) *hostUpdater{
    var h hostUpdater
    h.Log("Looking for SAM bridge on: " + samhost)
    h.Log("At port: " + samport)
    h.samBridgeClient, h.samBridgeErrors = goSam.NewClient(samhost + ":" + samport)
    h.Log("Connected to the SAM bridge on: " + samhost + ":" + samport)
    goSam.ConnDebug = debug
    h.parentList = upstream
    for _, url := range h.parseCsv(h.parentList) {
        h.Log("Upstream host providers: " + url)
    }
    h.tryFirst = parent
    for _, url := range h.parseCsv(h.tryFirst) {
        h.Log("New host providers: " + url)
    }
    h.hostfile = hostfile
    h.Log("Where to store hosts files: " + hostfile)
    dat, err := ioutil.ReadFile(h.hostfile)
    if ! h.Warn(err, "Error reading host file, may take a moment to start up.") {
        h.hostList = append(h.hostList, h.parseKvp(string(dat)))
    }
    h.retries = retries
    log.Println("Retry limit for requesting new hosts: ", retries)
    return &h
}
