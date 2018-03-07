package main

import (
    "log"
)

type jumpService struct{
    host string
    port string
}

func (jumpsite *jumpService) parseCsv(s string){

}

func (jumpsite *jumpService) fullAddress() string{
    return jumpsite.host + ":" + jumpsite.port
}

func (jumpsite *jumpService) Log(s string){
    if loglevel > 3 {
        log.Println(s)
    }
}

func (jumpsite *jumpService) Warn(s string){
    if loglevel > 2 {
        log.Println(s)
    }
}

func (jumpsite *jumpService) Fatal(s string){
    if loglevel > 1 {
        log.Println(s)
    }
}

func newJumpService(host string, port string) *jumpService{
    var j jumpService
    log.Println("Listening at: ", host)
    j.host = host
    log.Println("At port: ", port)
    j.port = port
    return &j
}
