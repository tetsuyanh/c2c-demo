package repository

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/tetsuyanh/c2c-demo/conf"
)

func TestMain(m *testing.M) {
	cp := conf.Postgres{
		DbUser:     "c2c_test",
		DbPassword: "",
		DbHost:     "localhost",
		DbName:     "c2c_test",
	}
	if eSetup := Setup(cp); eSetup != nil {
		log.Printf("failed to setup: %s", eSetup)
		os.Exit(1)
	}
	defer TearDown()

	clearAll()
	code := m.Run()
	os.Exit(code)
}

func clearAll() {
	// order to delete for constraint
	repo.dbMap.Exec(clear(AuthenticationTableName))
	repo.dbMap.Exec(clear(SessionTableName))
	repo.dbMap.Exec(clear(UserTableName))
}

func clear(m string) string {
	return fmt.Sprintf("delete from %s", m)
}
