package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v4"
)

func main() {
	log.Printf("Parsing the Command line Arguments ... ")

	port := flag.Int("port", -1, "Port to run the HTTP server on")
	dataDir := flag.String("dataDir", "/tmp/data", "Directory to store the BadgerDB data")
	verbose := flag.Bool("v", false, "Enable verbose logging")
	flag.Parse()

	if *port == -1 {
		log.Fatalf("Port must be specified using -port flag")
		panic("Port not specified")
	}

	

	log.Printf("Starting BadgerDB at %s ... ", *dataDir)
	opts := badger.DefaultOptions(*dataDir).WithLogger(nil)
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalf("Failed to open BadgerDB: %v", err)
		panic(err)
	}
	defer db.Close()

	
	
	ctx , cancel := context.WithCancel(context.Background())

	go runGarbageCollector(ctx , db)

	log.Printf("Starting HTTP server on port %d ... ", *port)

	app := newApp(db)

	if !*verbose {
		log.SetOutput(io.Discard)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), app)

	
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
		if err == http.ErrServerClosed {
			log.Printf("Server closed")
		}
		panic(err)
	}


	quit := make(chan os.Signal , 1)

	signal.Notify(quit , syscall.SIGINT , syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")
	cancel()	

    log.Println("Server exited properly")

}


func runGarbageCollector(ctx context.Context, db *badger.DB) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            
            log.Println("GC loop stopping...")
            return
        case <-ticker.C:
            // 2. Or if the timer ticked, run GC
        again:
            err := db.RunValueLogGC(0.5)
            if err == nil {
                // GC cleaned something, try again immediately
                goto again
            }
        }
    }
}