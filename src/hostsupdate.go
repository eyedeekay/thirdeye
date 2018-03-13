package main

import (
	"github.com/eyedeekay/gosam"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type hostUpdater struct {
	samHost         string
	samPort         string
	samBridgeClient *goSam.Client
	samBridgeErrors error

	retries    int
	tryFirst   string
	parentList string

	hostList [][]string
	hostfile string
}

func (updater *hostUpdater) parseCsv(s string) []string {
	hosts := strings.Split(s, ",")
	return hosts
}

func exists(file string) (bool, error) {
	if _, err := os.Stat(file); err == nil {
		return true, err
	} else {
		return false, err
	}
}

func (updater *hostUpdater) parseKvp(s string) [][]string {
	hosts := &[][]string{}
	for _, host := range strings.Split(s, "\n") {
		kv := strings.SplitN(host, "=", 2)
		if len(kv) == 2 {
			*hosts = append(*hosts, kv)
		}
	}
	return *hosts
}

func (updater *hostUpdater) writeHostList() error {
	exist, _ := exists(updater.hostfile)
	if exist {
		os.Remove(updater.hostfile)
	}
	f, err := os.Create(updater.hostfile)
	updater.Fatal(err, "File I/O errors")
	defer f.Close()
	for _, t := range updater.hostList {
		if len(t) == 2 {
			line := t[0] + "=" + t[1] + "\n"
			f.WriteString(line)
		}
	}
	return err
}

func (updater *hostUpdater) parseNl(s string) []string {
	hosts := []string{}
	for index, host := range strings.Split(s, "\n") {
		updater.Log(host)
		if index-1 > 0 {
			if host != hosts[index-1] {
				hosts = append(hosts, host)
			}
		} else {
			hosts = append(hosts, host)
		}
	}
	return hosts
}

func (updater *hostUpdater) sortHostList() [][]string {
	dat, err := ioutil.ReadFile(updater.hostfile)
	tempHostList := []string{}
	if !updater.Warn(err, "Error reading host file, may take a moment to start up.") {
		updater.Log("Local host file read into slice")
		tempHostList = append(tempHostList, updater.parseNl(string(dat))...)
	}
	sort.Strings(tempHostList)
	updater.hostList = nil
	updater.hostList = [][]string{[]string{}, []string{}}
	for index, host := range tempHostList {
		if index-1 > 0 {
			if !(host == tempHostList[index-1]) {
				updater.hostList = append(updater.hostList, strings.SplitN(host, "=", 2))
			}
		} else {
			updater.hostList = append(updater.hostList, strings.SplitN(host, "=", 2))
		}
	}
	return updater.hostList
}

func (updater *hostUpdater) hostUpdate() {
	t := updater.retries
	var hostList string
	for t >= 1 {
		updater.Log("Getting updates from new host providers.")
		for _, u := range updater.parseCsv(updater.tryFirst) {
			if done, h := updater.get(u); done {
				hostList += h + "\n"
				updater.hostList = append(updater.hostList, updater.parseKvp(hostList)...)
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
				updater.hostList = append(updater.hostList, updater.parseKvp(hostList)...)
				break
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
		t--
	}
	updater.writeHostList()
	updater.sortHostList()
	updater.writeHostList()
	updater.Log("Updates complete.")
}

func (updater *hostUpdater) get(s string) (bool, string) {
	tr := &http.Transport{
		Dial: updater.samBridgeClient.Dial,
	}

	updater.Log("Fetching updates from: " + s)

	client := &http.Client{Transport: tr}

	resp, err := client.Get(s)

	r := ""

	t := false

	if !updater.Warn(err, "Updater Error: "+s+" ") {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if !updater.Warn(err, "Response error: ") {
			r += string(body)
			t = true
		}
	}

	return t, r
}

func (updater *hostUpdater) getHosts() [][]string {
	return updater.hostList
}

func (updater *hostUpdater) Log(s string) {
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

func (updater *hostUpdater) Fatal(err error, s string) bool {
	if loglevel > 0 {
		if err != nil {
			log.Println(s, err)
			return true
		}
		return false
	}
	return false
}

func (updater *hostUpdater) loadHosts() [][]string {
	dat, err := ioutil.ReadFile(updater.hostfile)
	updater.hostList = [][]string{[]string{}, []string{}}
	if !updater.Warn(err, "Error reading host file, may take a moment to start up.") {
		updater.Log("Local host file read into slice")
		updater.hostList = append(updater.hostList, updater.parseKvp(string(dat))...)
	}
	updater.sortHostList()
	updater.writeHostList()
	return updater.hostList
}

func newHostUpdater(samhost string, samport string, retries int, upstream string, parent string, hostfile string, debug bool) *hostUpdater {
	var h hostUpdater
	h.samHost = samhost
	h.samPort = samport
	h.Log("Looking for SAM bridge on: " + samhost)
	h.Log("At port: " + samport)
	h.samBridgeClient, h.samBridgeErrors = goSam.NewClient(samhost + ":" + samport)
	goSam.ConnDebug = debug
	h.Log("Connected to the SAM bridge on: " + samhost + ":" + samport)
	h.parentList = upstream
	h.tryFirst = parent
	h.hostfile = hostfile
	h.Log("Where to store hosts files: " + hostfile)
	h.loadHosts()
	h.retries = retries
	return &h
}
