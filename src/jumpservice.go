package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type jumpService struct {
	title    string
	desc     string
	host     string
	port     string
	logwl    string
	hostfile string

	serverErr error
	mux       *http.ServeMux
	hostList  [][]string
}

func (jumpsite *jumpService) initMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, you hit bar!")
	})

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, you hit favicon.ico!")
	})

	mux.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, you hit style.css!")
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u := strings.Split(strings.Replace(r.URL.Path, "/", "", 1), "/")
		l := len(u)
		if test, handle := jumpsite.handle404(l, u, w, r); test {
			if handle == 1 {
				jumpsite.provideHosts(u[0], w, r)
			} else if handle == 2 {
				jumpsite.handleJump(u[1], w, r)
			}
		}
	})
	return mux
}

func (jumpsite *jumpService) emitHeader(w http.ResponseWriter, r *http.Request) bool {
	fmt.Fprintln(w, "<!DOCTYPE html>")
	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "  <head>")
	fmt.Fprintln(w, "    <title>"+jumpsite.title+"</title>")
	fmt.Fprintln(w, "    <meta charset=\"utf-8\">")
	fmt.Fprintln(w, "    <meta name=\"description\" content=\""+jumpsite.desc+"\">")
	fmt.Fprintln(w, "    <link rel=\"stylesheet\" href=\"/style.css\">")
	fmt.Fprintln(w, "  </head>")
	fmt.Fprintln(w, "  <body>")
	return true
}

func (jumpsite *jumpService) getHosts() [][]string {
	return jumpsite.hostList
}

func (jumpsite *jumpService) emitFooter(w http.ResponseWriter, r *http.Request) bool {
	fmt.Fprintln(w, "  </body>")
	fmt.Fprintln(w, "</html>")
	return true
}

func (jumpsite *jumpService) checkWhiteList(s string) bool {
	for _, lwl := range jumpsite.parseCsv(jumpsite.logwl) {
		if s == lwl {
			return true
		}
	}
	return false
}

func (jumpsite *jumpService) hosts(w http.ResponseWriter, r *http.Request) bool {
	for _, t := range jumpsite.getHosts() {
		log.Println(len(t))
		if len(t) == 2 {
			line := t[0] + "=" + t[1]
			fmt.Fprintln(w, line)
		}
	}
	return true
}

func (jumpsite *jumpService) doJump(test string, w http.ResponseWriter, r *http.Request) bool {
	b := false
	for _, t := range jumpsite.getHosts() {
		if len(t) == 2 {
			if t[0] == test {
				val := strings.SplitN(t[1], "#", 2)
				line := "http://" + t[0] + "/?i2paddresshelper=" + val[0]
				w.Header().Set("Location", line)
				w.WriteHeader(301)
				jumpsite.emitHeader(w, r)
				fmt.Fprintln(w, "<h1>", "Looking up:", test, "... checking", jumpsite.length(), "hosts", "</h1>")
				fmt.Fprintln(w, "<pre><code>")
				fmt.Fprintln(w, "    ", line)
				fmt.Fprintln(w, "</pre></code>")
				jumpsite.emitFooter(w, r)
				b = true
				return b
			}
		}
	}
	if !b {
		w.WriteHeader(200)
		jumpsite.emitHeader(w, r)
		fmt.Fprintln(w, "<h1>", "Looking up:", test, "... checking", jumpsite.length(), "hosts", "</h1>")
		jumpsite.emitFooter(w, r)
	}
	return b
}

func (jumpsite *jumpService) provideHosts(s string, w http.ResponseWriter, r *http.Request) bool {
	if jumpsite.checkWhiteList(s) {
		jumpsite.hosts(w, r)
	} else {
		jumpsite.hosts(w, r)
		jumpsite.Log("Hosts were requested across the following URL " + s)
	}
	return true
}

func (jumpsite *jumpService) handleJump(s string, w http.ResponseWriter, r *http.Request) bool {
	jumpsite.doJump(s, w, r)
	return true
}

func (jumpsite *jumpService) handle404(l int, s []string, w http.ResponseWriter, r *http.Request) (bool, int) {
	if l == 1 {
		return true, l
	} else if l == 2 {
		if s[0] == "jump" {
			if s[1] != "" {
				return true, l
			} else {
				return true, 1
			}
		} else {
			w.WriteHeader(404)
			fmt.Fprintln(w, r.URL.Path)
			fmt.Fprintln(w, "You're lost, go home")
		}
	} else {
		w.WriteHeader(404)
		fmt.Fprintln(w, r.URL.Path)
		fmt.Fprintln(w, "You're lost, go home")
	}
	return false, 0
}

func (jumpsite *jumpService) length() int {
	log.Println("Data length", len(jumpsite.getHosts()))
	return len(jumpsite.getHosts())
}

func (jumpsite *jumpService) parseCsv(s string) []string {
	urls := strings.Split(s, ",")
	return urls
}

func (jumpsite *jumpService) parseKvp(s string) [][]string {
	hosts := &[][]string{}
	for _, host := range strings.Split(s, "\n") {
		kv := strings.SplitN(host, "=", 2)
		if len(kv) == 2 {
			jumpsite.Log(kv[0])
			*hosts = append(*hosts, kv)
		}
	}
	return *hosts
}

func (jumpsite *jumpService) fullAddress() string {
	return jumpsite.host + ":" + jumpsite.port
}

func (jumpsite *jumpService) Log(s string) {
	if loglevel > 2 {
		log.Println("LOG: ", s)
	}
}

func (jumpsite *jumpService) Warn(err error, s string) bool {
	if loglevel > 1 {
		if err != nil {
			log.Println("WARN: ", s)
			return true
		}
		return false
	}
	return false
}

func (jumpsite *jumpService) Fatal(err error, s string) bool {
	if loglevel > 0 {
		if err != nil {
			log.Println("FATAL: ", s)
			return true
		}
		return false
	}
	return false
}

func (jumpsite *jumpService) Serve() {
	jumpsite.Log("Initializing handler handle")
	jumpsite.length()
	if err := http.ListenAndServe(jumpsite.fullAddress(), jumpsite.mux); err != nil {
		jumpsite.Warn(err, "Fatal Error: server not started")
		jumpsite.serverErr = err
	}
}

func (jumpsite *jumpService) loadHosts() [][]string {
	dat, err := ioutil.ReadFile(jumpsite.hostfile)
    jumpsite.hostList = [][]string{nil, nil}
    var hostlist [][]string
	hostlist = [][]string{[]string{}, []string{}}
	if !jumpsite.Warn(err, "Error reading host file, may take a moment to start up.") {
		jumpsite.Log("Local host file read into slice")
		hostlist = append(hostlist, jumpsite.parseKvp(string(dat))...)
	}
	return hostlist
}

func newJumpService(host string, port string, title string, desc string, hostfile string, logwl string) *jumpService {
	var j jumpService
	j.logwl = logwl
	j.Log("setting log whitelist: ")
	j.title = title
	j.Log("setting title: " + j.title)
	j.desc = desc
	j.Log("setting description: " + j.desc)
	j.host = host
	j.Log("setting host: " + host)
	j.port = port
	j.Log("setting port: " + port)
	j.mux = j.initMux()
	j.Log("Listening at: " + j.host + " At port: " + j.port)
	log.Println("loading local jump service data:")
	j.hostfile = hostfile
	j.hostList = j.loadHosts()
	log.Println("Starting jump service web site: ", j.fullAddress())
	return &j
}
