package main

import (
    "log"
    "strings"
    "time"
    "github.com/eyedeekay/gosam"
)

type hostUpdater struct{
    retries int
    tryFirst []string
    parentList []string
    samBridgeClient *goSam.Client
    samBridgeErrors error

}

func (updater *hostUpdater) parseCsv(s string) []string{
    hosts := strings.Split(s, ",")
    return hosts
}

func (updater *hostUpdater) hostUpdate(){
    t := updater.retries
    for t > 1 {
        updater.Log("Getting updates")
        time.Sleep(time.Duration(100) * time.Second)
        t--
    }
}

func (updater *hostUpdater) Log(s string){
    if loglevel > 2 {
        log.Println(s)
    }
}

func (updater *hostUpdater) Warn(err error, s string){
    if loglevel > 1 {
        if err != nil {
            log.Println(s)
        }
    }
}

func (updater *hostUpdater) Fatal(err error,s string){
    if loglevel > 0 {
        if err != nil {
            log.Println(s)
        }
    }
}

func newHostUpdater(samhost string, samport string, retries int, upstream string, parent string, hostfile string) *hostUpdater{
    var h hostUpdater
    h.Log("Looking for SAM bridge on: " + samhost)
    h.Log("At port: " + samport)
    h.samBridgeClient, h.samBridgeErrors = goSam.NewClient(samhost + ":" + samport)
    h.Log("Connected to the SAM bridge on" + samhost + ":" + samport)

    h.parentList = h.parseCsv(upstream)
    for _, url := range h.parentList {
        h.Log("Upstream host providers: " + url)
    }
    h.tryFirst = h.parseCsv(parent)
    for _, url := range h.tryFirst {
        h.Log("Authoritative host providers: " + url)
    }
    h.Log("Where to store hosts files: " + hostfile)
    log.Println("Retry limit for requesting new hosts: ", retries)
    h.retries = retries
    return &h
}
