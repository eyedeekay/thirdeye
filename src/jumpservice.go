package main

import (
    "fmt"
    "net/http"
    "log"
)

type jumpService struct{
    host string
    port string
    serverErr error
}

func (jumpsite jumpService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
}

func (jumpsite jumpService) parseCsv(s string){

}

func (jumpsite jumpService) fullAddress() string{
    return jumpsite.host + ":" + jumpsite.port
}

func (jumpsite jumpService) Log(s string){
    if loglevel > 2 {
        log.Println(s)
    }
}

func (jumpsite jumpService) Warn(err error, s string){
    if loglevel > 1 {
        if err != nil {
            log.Println(s)
        }
    }
}

func (jumpsite jumpService) Fatal(err error,s string){
    if loglevel > 0 {
        if err != nil {
            log.Println(s)
        }
    }
}

func (jumpsite jumpService) Serve() error{
    jumpsite.Log("Initializing handler handle")
    if err := http.ListenAndServe(jumpsite.fullAddress(), jumpService{}); err != nil {
        jumpsite.Warn(err, "Fatal Error: server not started")
        jumpsite.serverErr = err
        return err
    }
    return nil
}

func newJumpService(host string, port string) *jumpService{
    var j jumpService
    j.host = host
    j.port = port
    j.Log("Listening at: " + host + "At port: " + port)
    go j.Serve()
	j.Fatal(j.serverErr, "Error starting server.")
    return &j
}
