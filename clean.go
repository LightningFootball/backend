package main

import (
	"github.com/LightningFootball/backend/base/exit"
	"github.com/LightningFootball/backend/base/log"
	"github.com/LightningFootball/backend/base/utils"
	"os"
)

func clean() {
	readConfig()
	initGorm()
	initLog()
	err := utils.CleanUpExpiredTokens()
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	exit.Close()
	exit.QuitWG.Wait()
	log.Fatal("Clean succeed!")
}
