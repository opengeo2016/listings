package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/opengeo2016/listings/db"

	"gopkg.in/mgo.v2"
)

var (
	dbAddr      = flag.String("db_addr", "localhost", "The address to connect to mongo db.")
	servicePort = flag.Int("service_port", 8001, "The port number for listing service.")
)

func main() {
	flag.Parse()

	// Creates mongo session.
	session, err := mgo.Dial(*dbAddr)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Sets up http handler.
	http.HandleFunc("/listings", MakeListHandler(db.NewListingDetailsAccessor(session)))

	log.Println("Start serving on: ", *servicePort)
	http.ListenAndServe(":"+strconv.Itoa(*servicePort), nil)
}
