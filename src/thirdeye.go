package main

import (
    "flag"
    "log"
    "time"
)

var loglevel int
var wait time.Duration

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
    newhosts := flag.String("newhosts", "http://stats.i2p/cgi-bin/newhosts.txt", "Fetch new hosts from here")
    upstream := flag.String("upstream", "http://i2p2.i2p/hosts.txt,http://i2host.i2p/cgi-bin/i2hostetag", "Fetch more hosts from here")
    hostfile := flag.String("hostfile", "etc/thirdeye/localhosts.txt", "Local hosts file")
    debug := flag.Bool("debug", false, "Print connection debug info" )
    verbosity := flag.Int("verbosity", 4, "Verbosity level: 0=Quiet 1=Fatal 2=Warning 3=Debug")

    flag.Parse()

    Title := *title
    Description := *desc
    LogWhiteList := *logwl
    Host := *host
    Port := *port

    Retries := *retries
    Interval := *interval
    NewHosts := *newhosts
    Upstream := *upstream
    HostFile := *hostfile

    SamConnHost := *samhost
    SamConnPort := *samport
    Debug := *debug
    Verbosity := *verbosity


    loglevel = Verbosity

    log.Println("Log level: ", *verbosity)
    wait = time.Duration(Interval) * time.Hour
    log.Println("Updater Interval: ", Interval, wait )
    log.Println("Updating hosts...")
    hostUpdater := newHostUpdater(SamConnHost,
        SamConnPort,
        Retries,
        Upstream,
        NewHosts,
        HostFile,
        Debug)

    hostUpdater.Log("Hostupdater created.")

    jumpService := newJumpService(Host,
        Port,
        Title,
        Description,
        HostFile,
        LogWhiteList)
    go jumpService.Serve()
    for true {
        hostUpdater.hostUpdate()
        jumpService.loadHosts()
        time.Sleep(wait)
    }

}
