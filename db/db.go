package db

import (
	"context"
	"fmt"
	"log"

	"em/config"

	"github.com/jackc/pgx/v5"
)


var Conn *pgx.Conn

func InitDB(cfg *config.Config){
	var err error

	Conn,err = pgx.Connect(context.Background(),cfg.DataName())
	if err != nil{
		log.Fatalf("Failed to load DB: %v", err)
	}

	fmt.Println("Connected",Conn)
}