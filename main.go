package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/nilesh0729/OrdinaryBank/api"
	Anuskh "github.com/nilesh0729/OrdinaryBank/db/Result"
	"github.com/nilesh0729/OrdinaryBank/util"
)


func main(){
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot Load Config", err)
	}
	
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot Connect To db: ", err)
	}


	store := Anuskh.NewTxConn(conn)
	server, err := api.NewServer(store, config)
	if err !=nil{
		fmt.Printf("cannot create Server: %v", err)
	}


	err = server.Start(config.ServerAddress)
	if err != nil{
		log.Fatal("Cannot Start Server : ", err)
	}
}