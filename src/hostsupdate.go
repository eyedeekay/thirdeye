package main

import (
	"github.com/eyedeekay/gosam"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	//	"time"
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
	Fatal(err, "File I/O errors")
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
	for _, host := range strings.Split(s, "\n") {
		hosts = append(hosts, host)
	}
	return hosts
}

func (updater *hostUpdater) sortHostList() [][]string {
	dat, err := ioutil.ReadFile(updater.hostfile)
	tempHostList := []string{}
	if !Warn(err, "Error reading host file, may take a moment to start up.") {
		Log("Local host file read into slice")
		tempHostList = append(tempHostList, updater.parseNl(string(dat))...)
	}
	sort.Strings(tempHostList)
	var newHostList [][]string
	newHostList = [][]string{[]string{}, []string{}}
	for index, host := range tempHostList {
		if index-1 > 0 {
			if !(host == tempHostList[index-1]) {
				newHostList = append(newHostList, strings.SplitN(host, "=", 2))
			} else {
				Log(host, tempHostList[index-1])
			}
		} else {
			newHostList = append(newHostList, strings.SplitN(host, "=", 2))
		}
	}
	return newHostList
}

func (updater *hostUpdater) pullUpdate(url string) [][]string {
	for _, u := range updater.parseCsv(url) {
		t := updater.retries
		var hostList string
		for t >= 1 {
			Log("Getting updates from new host providers.")
			if done, h := updater.get(u); done {
				hostList += h + "\n"
				//updater.hostList = append(updater.hostList, updater.parseKvp(hostList)...)
				return append(updater.hostList, updater.parseKvp(hostList)...)
			}
		}
		t--
	}
	return updater.hostList
}

func (updater *hostUpdater) hostUpdate() {
	updater.hostList = updater.pullUpdate(updater.tryFirst)
	updater.hostList = updater.pullUpdate(updater.parentList)
	updater.writeHostList()
	updater.hostList = [][]string{nil, nil}
	updater.hostList = updater.sortHostList()
	updater.writeHostList()
	Log("Updates complete.")
}

func (updater *hostUpdater) get(s string) (bool, string) {
	//var samBridgeClient *goSam.Client
	tr := &http.Transport{
		Dial: updater.Dial,
	}

	Log("Fetching updates from: " + s)

	client := &http.Client{Transport: tr}
	resp, err := client.Get(s)

	r := ""
	t := false

	if !Warn(err, "Updater Error: "+s+" ") {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if !Warn(err, "Response error: ") {
			r += string(body)
			t = true
		}
	}

	return t, r
}

func (updater *hostUpdater) getHosts() [][]string {
	return updater.hostList
}

func (updater *hostUpdater) Dial(network, addr string) (net.Conn, error) {
	portIdx := strings.Index(addr, ":")
	if portIdx >= 0 {
		addr = addr[:portIdx]
	}
	addr, err := updater.samBridgeClient.Lookup(addr)
	if err != nil {
		return nil, err
	}

	id, _, err := updater.samBridgeClient.CreateStreamSession("")
	if err != nil {
		return nil, err
	}

	newC, err := goSam.NewClient(updater.samHost + ":" + updater.samPort)
	if err != nil {
		return nil, err
	}

	err = newC.StreamConnect(id, addr)
	if err != nil {
		return nil, err
	}

	return newC.SamConn, nil
}

func (updater *hostUpdater) loadHosts() [][]string {
	dat, err := ioutil.ReadFile(updater.hostfile)
	var hostlist [][]string
	hostlist = [][]string{[]string{}, []string{}}
	if !Warn(err, "Error reading host file, may take a moment to start up.") {
		Log("Local host file read into slice")
		hostlist = append(hostlist, updater.parseKvp(string(dat))...)
	}
	return hostlist
}

func newHostUpdater(samhost string, samport string, retries int, upstream string, parent string, hostfile string, debug bool) *hostUpdater {
	var h hostUpdater
	h.samHost = samhost
	h.samPort = samport
	Log("Looking for SAM bridge on: " + h.samHost)
	Log("At port: " + h.samPort)
	h.samBridgeClient, h.samBridgeErrors = goSam.NewClient(h.samHost + ":" + h.samPort)
	//goSam.ConnDebug = debug
	Log("Connected to the SAM bridge on: " + samhost + ":" + samport)
	h.parentList = upstream
	h.tryFirst = parent
	h.hostfile = hostfile
	Log("Where to store hosts files: " + hostfile)
	h.hostList = h.loadHosts()
	h.retries = retries
	h.writeHostList()
	h.hostList = [][]string{nil, nil}
	h.hostList = h.sortHostList()
	return &h
}

//func newHostUpdaterFromOptions(opts ...func(*hostUpdater) error) (*hostUpdater, error) {
//}
