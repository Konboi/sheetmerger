package main

import (
	"flag"
	"log"
	"time"

	"github.com/Konboi/sheetmerger"
	"github.com/soh335/sliceflag"
)

func main() {
	var mergeSheetNames []string
	var baseSheetKeyID, diffSheetKeyID string
	var config string

	flag.StringVar(&baseSheetKeyID, "base", "", "base sheet key")
	flag.StringVar(&diffSheetKeyID, "diff", "", "diff sheet key")
	flag.StringVar(&config, "c", "", "config file path")
	sliceflag.StringVar(flag.CommandLine, &mergeSheetNames, "name", []string{}, "set merge sheet names")
	flag.Parse()

	conf, err := sheetmerger.NewConfig(config)
	if err != nil {
		log.Fatalln("error read config file", err.Error())
	}
	sm, err := sheetmerger.New(conf)
	if err != nil {
		log.Fatalln("error new sheet merger", err.Error())
	}

	if err := sm.Backup(
		baseSheetKeyID,
		time.Now().Format("2006/01/02 15:04:05"),
	); err != nil {
		log.Fatalln("error backup base file", err.Error())
	}

	if err := sm.Merge(baseSheetKeyID, diffSheetKeyID, mergeSheetNames...); err != nil {
		log.Fatalln("error merge file", err.Error())
	}

	flag.Parse()
}
