package main

import (
	"fmt"

	"github.com/yakabuff/discord-dl/archiver"
	"github.com/yakabuff/discord-dl/db"
)

func main() {
	var archiver = archiver.Archiver{}
	args := archiver.InitCli()
	archiver.Args = args
	db, err := db.Init_db(archiver.Args.Output)
	if err != nil {
		panic(err.Error())
	}

	errDg, dg := archiver.CreateConnection()
	if errDg != nil {
		fmt.Println(errDg)
	}
	archiver.Dg = dg

	archiver.Db = *db
	err = archiver.ParseCliArgs()
	if err != nil {
		fmt.Println(err)
	}

}
