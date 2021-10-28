package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	sm "github.com/rmanzoku/go-secretmanager"
	"github.com/rmanzoku/go-secretmanager/smutil"
)

func usage() {
	fmt.Println("Usage of awssecrets:")
	fmt.Println("   get [key]          Get secret value")
	fmt.Println("   set [key] [value]  Create/Update key")
	fmt.Println("   del [key]          Delete secret")
}

func main() {
	var err error
	_ = flag.NewFlagSet("get", flag.ExitOnError)
	_ = flag.NewFlagSet("set", flag.ExitOnError)
	_ = flag.NewFlagSet("del", flag.ExitOnError)

	if len(os.Args) < 2 {
		usage()
		return
	}

	svc, err := smutil.NewSMClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	switch os.Args[1] {
	case "get":
		if len(os.Args) != 3 {
			usage()
			return
		}
		key := os.Args[2]
		ret, err := sm.Get(ctx, svc, key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ret)

	case "set":
		if len(os.Args) != 4 {
			usage()
			return
		}

		key := os.Args[2]
		value := os.Args[3]
		err = sm.Set(ctx, svc, key, value)
		if err != nil {
			log.Fatal(err)
		}

	case "del":
		if len(os.Args) != 3 {
			usage()
			return
		}
		key := os.Args[2]
		err = sm.Del(ctx, svc, key)
		if err != nil {
			log.Fatal(err)
		}

	default:
		usage()
	}

}
