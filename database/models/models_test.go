package models

import (
	"github.com/LightningFootball/backend/base/exit"
	"github.com/LightningFootball/backend/database"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	defer database.SetupDatabaseForTest()()
	defer exit.SetupExitForTest()()
	os.Exit(m.Run())
}
