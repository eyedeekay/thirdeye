package main

import (
    "flag"
    "log"
    "time"
)

var loglevel int
var debug bool

func main(){
    title := flag.String("title", "Thirdeye Based Jump Service", "Title of the service.")
    desc := flag.String("desc", "Thirdeye based jump service", "Brief description of the service.")
    logwl := flag.String("logwl", "", "Whitelist of urls to never log")
    samhost := flag.String("samhost", "127.0.0.1", "Host address to listen on.")
    samport := flag.String("samport", "7656", "Port to listen on.")
    host := flag.String("host", "127.0.0.1", "Host address to listen on.")
    port := flag.String("port", "8053", "Port to listen on.")
    retries := flag.Int("retries", 2, "Number of attempts to fetch new hosts")
    interval := flag.Int("interval", 6, "Hours between updatess")
    upstream := flag.String("upstream", "http://i2p2.i2p/hosts.txt,http://i2host.i2p/cgi-bin/i2hostetag", "Fetch more hosts from here")
    parent := flag.String("newhosts", "http://stats.i2p/cgi-bin/newhosts.txt", "Fetch new hosts from here")
    hostfile := flag.String("hostfile", "etc/thirdeye/localhosts.txt", "Local hosts file")
    debug = *flag.Bool("debug", false, "Print connection debug info" )
    verbosity := flag.Int("verbosity", 4, "Verbosity level: 0=Quiet 1=Fatal 2=Warning 3=Debug")

    flag.Parse()

    loglevel = *verbosity

    log.Println("Log level: ", *verbosity)
    wait := time.Duration(*interval) * time.Hour
    log.Println("Updater Interval: ", *interval, wait )

    hostsData := newHostUpdater(*samhost, *samport, *retries, *upstream, *parent, *hostfile)

    jumpService := newJumpService(*host, *port, *title, *desc, *logwl)

    log.Println("loading local jump service data:")
    jumpService.newData(hostsData.getHosts())

    log.Println("Starting jump service web site: ", jumpService.fullAddress())

    for true {
        log.Println("Starting loop")
        hostsData.hostUpdate()
        jumpService.newData(hostsData.getHosts())
        time.Sleep(wait)
    }

}
