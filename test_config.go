package main

import (
	"github.com/LightningFootball/backend/base/log"
)

func testConfig() {
	// TODO: test config using config.Get
	readConfig()
	initGorm(false)
	initLog()
	initRedis()
	initStorage()
	initWebAuthn()
	log.Fatalf("should success.")
}
