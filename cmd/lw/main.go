package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/feduxorg/local-webserver/internal/cli"
	"github.com/gorilla/handlers"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		fmt.Print("Captured CTRL-C. Exit..\n")

		os.Exit(0)
	}()

	conf := cli.Config{}
	conf.ParseArgs()

	if conf.NetworkInterface == "" {
		conf.NetworkInterface = cli.DeterminInterfaceToListenOn()
		fmt.Println("")
	}

	if strings.Count(conf.NetworkInterface, ":") > 1 {
		conf.NetworkInterface = fmt.Sprintf("[%s]", conf.NetworkInterface)
	}

	url := fmt.Sprintf("http://%s:%d/index.html", conf.NetworkInterface, conf.Port)
	listen := fmt.Sprintf("%s:%d", conf.NetworkInterface, conf.Port)

	if conf.OpenBrowser {
		openBrowser(url, conf.Silent)
	}

	startServer(conf.WorkingDirectory, listen, conf.Silent)
}

func openBrowser(url string, silent bool) {
	if silent == false {
		fmt.Printf("Open browser with %s\n", url)
	}

	err := open.Start(url)

	if err != nil {
		log.Fatal(err)
	}
}

func startServer(directory string, bind string, silent bool) {
	if silent == false {
		fmt.Printf("Server listens on %s\n", bind)
		fmt.Printf("\nRequests:\n")
	}

	var http_handler http.Handler

	if silent == true {
		http_handler = http.FileServer(http.Dir(directory))
	} else {
		http_handler = handlers.CombinedLoggingHandler(os.Stdout, http.FileServer(http.Dir(directory)))
	}

	http.Handle("/", http_handler)
	panic(http.ListenAndServe(bind, nil))
}
