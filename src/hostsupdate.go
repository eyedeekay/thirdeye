package main

import (
    "log"
    "github.com/eyedeekay/gosam"
)

type hostUpdater struct{
    retries int
    parentList []string
    samBridgeClient *goSam.Client

}

func (updater *hostUpdater) parseCsv(s string){

}

func (updater *hostUpdater) hostUpdate(){
    t := updater.retries
    for t < 1 {
        updater.Log("Getting updates")
    }
}

func (updater *hostUpdater) Log(s string){
    if loglevel > 3 {
        log.Println(s)
    }
}

func (updater *hostUpdater) Warn(s string){
    if loglevel > 2 {
        log.Println(s)
    }
}

func (updater *hostUpdater) Fatal(s string){
    if loglevel > 1 {
        log.Println(s)
    }
}

func newHostUpdater(samhost string, samport string, retries int, upstream string, parent string, hostfile string) *hostUpdater{
    var h hostUpdater
    log.Println("Looking for SAM bridge on: ", samhost)
    log.Println("At port: ", samport)

    log.Println("Upstream host providers: ", upstream)
    log.Println("Authoritative host provider: ", parent)
    log.Println("Where to store hosts files: ", hostfile)

    log.Println("Retry limit for requesting new hosts: ", retries)
    h.retries = retries
    return &h
}
