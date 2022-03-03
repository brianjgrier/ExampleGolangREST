package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	myDB      *sql.DB
	dbConnErr error
	myIP      string
)

func main() {
	myIP, _ = externalIP()

	fmt.Println("Rest to SQL example @", myIP)

//	myDB, dbConnErr = sql.Open("mysql", "golang-test:test-golang@tcp(sds-server1:3306)/golang_rest_sql_test")
	myDB, dbConnErr = sql.Open("mysql", "golang-test:test-golang@tcp(mysqldb:3306)/golang_rest_sql_test")
	if dbConnErr != nil {
		fmt.Println("Cannot connect to DB using only log statements")
		//panic(err.Error())
	} else {
		fmt.Println("Connected to DB")
		defer myDB.Close()
	}

	now := time.Now()
	fmt.Println("Current Timestamp: ", now)

	//r := mux.NewRouter()
	r := http.NewServeMux()
	r.HandleFunc("/", handler)
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		fmt.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("Could not start http server: Fatal(err)")
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)

}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	//os.Exit(0)
}

//
// Get the ip addrsess of the first interface that you come accross
//
func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
			//fmt.Println("Found Ip addr: ", ip.String())
		}
	}
	return "", errors.New("are you connected to the network?")
}

func recordRequest(payload string, theTime time.Time) {
	fmt.Println("Writing to SQL db")
	stmt, err := myDB.Prepare("INSERT INTO eventlog (now, payload, ipaddr) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		return
	}
	res, execErr := stmt.Exec(theTime, payload, myIP)

	if execErr != nil {
		fmt.Println("Found: ", err)
	}
	fmt.Println("Inserted ", res)
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}

	fmt.Printf("Received request for %s\n", name)

	if myDB != nil {
		recordRequest(name, time.Now())
	}

	w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
}
