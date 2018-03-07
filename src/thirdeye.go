package main

import (
    "flag"
    "log"
    "time"
)

var loglevel int

func main(){
    samhost := flag.String("samhost", "127.0.0.1", "Host address to listen on.")
    samport := flag.String("samport", "7656", "Port to listen on.")
    host := flag.String("host", "127.0.0.1", "Host address to listen on.")
    port := flag.String("port", "8053", "Port to listen on.")
    retries := flag.Int("retries", 5, "Number of attempts to fetch new hosts")
    interval := flag.Int("interval", 6, "Hours between updatess")
    upstream := flag.String("upstream", "http://i2p2.i2p/hosts.txt,http://i2host.i2p/cgi-bin/i2hostetag", "")
    parent := flag.String("parent", "http://stats.i2p/cgi-bin/newhosts.txt", "Fetch new hosts from here")
    hostfile := flag.String("hostfile", "localhosts.txt", "Local hosts file")
    verbosity := flag.Int("verbosity", 0, "Verbosity level: 0=Quiet 1=Fatal 2=Warning 3=Debug")

    flag.Parse()


    loglevel = *verbosity
    log.Println("Log level: ", *verbosity)
    wait := time.Duration(*interval) * time.Second
    log.Println("Updater Interval: ", *interval, wait )

    hostsData := newHostUpdater(*samhost, *samport, *retries, *upstream, *parent, *hostfile)

    jumpService := newJumpService(*host, *port)
    log.Println("Starting jump service web site: ", jumpService.fullAddress())

    for true {
        log.Println("Getting updates")
        hostsData.hostUpdate()
        time.Sleep(wait)
    }

}
