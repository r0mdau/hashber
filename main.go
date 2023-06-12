package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/serialx/hashring"
)

var list *memberlist.Memberlist
var ring *hashring.HashRing

const port = ":8090"

func hello(w http.ResponseWriter, req *http.Request) {
	server, _ := ring.GetNode("hello")
	fmt.Println("Selected node", server, "from the hashring")

	if os.Getenv("MY_POD_IP") != server {
		fmt.Println("I am forwarding hello")
		resp, err := http.Get(fmt.Sprintf("http://%s%s/hello", server, port))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("Response status:", resp.Status)

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
			os.Exit(1)
		}
		fmt.Fprint(w, string(respBody))
	} else {
		msg := fmt.Sprintf("Hello from %s\n", server)
		fmt.Fprint(w, msg)
		fmt.Println("Hello from me \\o/")
	}
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func memberBeat() {
	for range time.Tick(time.Second * 3) {
		servers := []string{}
		for _, member := range list.Members() {
			servers = append(servers, member.Addr.String())
		}
		ring = hashring.New(servers)
	}
}

func memberList() {
	var err error
	list, err = memberlist.Create(memberlist.DefaultLocalConfig())
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	// Join an existing cluster by specifying at least one known member.
	_, err = list.Join([]string{"hashber-memberlist.default.svc.cluster.local"})
	if err != nil {
		fmt.Printf("Failed to join cluster: " + err.Error())
	}

	// Ask for members of the cluster
	for _, member := range list.Members() {
		fmt.Printf("Startup member list : %s %s\n", member.Name, member.Addr)
	}
}

func main() {

	// Start the gossip
	go memberList()
	go memberBeat()

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(port, nil)
}
