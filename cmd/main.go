package main

import (
	"encoding/gob"
	"fmt"
	"gdn/config"
	"gdn/internal/guardian"
	"os"

	"github.com/patrickmn/go-cache"
)

var gdn *guardian.Guardian

func init() {
	gob.Register(map[string]string{})
	cache := cache.New(cache.NoExpiration, cache.NoExpiration)
	if _, err := os.Stat(config.CachePath); err != nil {
		if os.IsExist(err) {
			cache.LoadFile(config.CachePath)
		}
	} else {
		cache.LoadFile(config.CachePath)
	}
	gdn = guardian.NewGuardian(cache, config.LogPath)
}

func main() {
	if len(os.Args) == 1 {
		echoHelp()
		return
	}
	defer gdn.SaveList("/etc/gdn/cache")
	switch os.Args[1] {
	case "list":
		gdn.ShowList()
	case "stop":
		gdn.StopProc(os.Args[2])
	case "log":
		gdn.ShowLog(os.Args[2])
	case "help":
		echoHelp()
	case "version":
	default:
		gdn.StartProc(os.Args[1], os.Args[2:])
	}
}

func echoHelp() {
	fmt.Printf(
		"\n gdn: run command as daemon" +
			"\n\n     <command>       run your command" +
			"\n     list            show running commands" +
			"\n     stop <proc>     stop a command by SIG" +
			"\n     log  <proc>     view log of command" +
			"\n\n     help            show help" +
			"\n     version         show version\n",
	)
}
