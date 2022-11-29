package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/memberlist"
)

var list *memberlist.Memberlist

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	// Start the gossip
	go memberList()
	go memberBeat()

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8090", nil)
}

func memberList() {
	var err error
	list, err = memberlist.Create(memberlist.DefaultLocalConfig())
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	// Join an existing cluster by specifying at least one known member.
	_, err = list.Join([]string{"hashber.default.svc.cluster.local"})
	if err != nil {
		fmt.Printf("Failed to join cluster: " + err.Error())
	}

	// Ask for members of the cluster
	for _, member := range list.Members() {
		fmt.Printf("Member list youpi : %s %s\n", member.Name, member.Addr)
	}
}

func memberBeat() {
	for range time.Tick(time.Second * 3) {
		fmt.Println("--- Start member beat ---")
		for _, member := range list.Members() {
			fmt.Printf("Member list beat : %s %s\n", member.Name, member.Addr)
		}
	}
}
