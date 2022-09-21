package main

import (
	"fmt"
	"github/guanhg/syncDB-search/config"
)

func testConfig() {
	cfg := config.LoadJsonFileConfig()

	fmt.Println(cfg)
}
