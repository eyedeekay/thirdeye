package main

import (
	"flag"
	"log"
    "os"
	"time"
)

var loglevel int
var wait time.Duration

func main() {
	title := flag.String("title", "Thirdeye Based Jump Service", "Title of the service.")
	desc := flag.String("desc", "Thirdeye based jump service", "Brief description of the service.")
	logwl := flag.String("logwl", "", "Whitelist of urls to never log")
	samhost := flag.String("samhost", "127.0.0.1", "Host address to attach to SAM bridge.")
	samport := flag.String("samport", "7656", "SAM port.")
	host := flag.String("host", "0.0.0.0", "Host address to listen on.")
	port := flag.String("port", "8053", "Port to listen on.")
	retries := flag.Int("retries", 0, "Number of attempts to fetch new hosts")
	interval := flag.Int("interval", 6, "Hours between updatess")
	newhosts := flag.String("newhosts", "http://inr.i2p/export/alive-hosts.txt", "Fetch new hosts from here")
	upstream := flag.String("upstream", "http://inr.i2p/export/alive-hosts.txt", "Fetch more hosts from here")

	hostfile := flag.String("hostfile", "etc/thirdeye/localhosts.txt", "Local hosts file")

	cssfile := flag.String("cssfile", "etc/thirdeye/style.css", "Local css file")

	icofile := flag.String("icofile", "etc/thirdeye/favicon.ico", "Local favicon file")

	debug := flag.Bool("debug", false, "Print connection debug info")
	verbosity := flag.Int("verbosity", 2, "Verbosity level: 0=Quiet 1=Fatal 2=Warning 3=Debug")

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

	CssFile := *cssfile
	IconFile := *icofile

	SamConnHost := *samhost
	SamConnPort := *samport
	Debug := *debug
	Verbosity := *verbosity

	loglevel = Verbosity

	log.Println("Log level: ", *verbosity)
	wait = time.Duration(Interval) * time.Hour
	log.Println("Updater Interval: ", Interval, wait)
	log.Println("Updating hosts...")
	hostUpdater := newHostUpdater(SamConnHost,
		SamConnPort,
		Retries,
		Upstream,
		NewHosts,
		HostFile,
		Debug)

	Log("Hostupdater created.")

	jumpService := newJumpService(Host,
		Port,
		Title,
		Description,
		HostFile,
		LogWhiteList,
		CssFile,
		IconFile)
	go jumpService.Serve()
	for true {
		hostUpdater.hostUpdate()
		jumpService.hostList = jumpService.loadHosts()
		time.Sleep(wait)
	}

}

func Log(s ...string) {
	if loglevel > 2 {
		log.Println("LOG: ", s)
	}
}

func Warn(err error, s string) bool {
	if loglevel > 1 {
		if err != nil {
			log.Println("WARN:", s, err)
			return true
		}
		return false
	}
	return false
}

func Fatal(err error, s string) bool {
	if loglevel > 0 {
		if err != nil {
			log.Fatal("FATAL: ", s, err)
			return true
		}
		return false
	}
	return false
}

func exists(file string) (bool, error) {
	if _, err := os.Stat(file); err == nil {
		return true, err
	} else {
		return false, err
	}
}
