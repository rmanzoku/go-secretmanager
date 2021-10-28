package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/rmanzoku/go-secretmanager/smutil"
)

func get(ctx context.Context, svc *secretsmanager.Client, key string) (string, error) {
	params := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}
	out, err := svc.GetSecretValue(ctx, params)
	if err != nil {
		return "", err
	}
	return *out.SecretString, nil
}

func set(ctx context.Context, svc *secretsmanager.Client, key, value string) error {
	params := &secretsmanager.CreateSecretInput{
		Name:         aws.String(key),
		SecretString: aws.String(value),
	}
	_, err := svc.CreateSecret(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func usage() {
}

func main() {
	var err error
	_ = flag.NewFlagSet("get", flag.ExitOnError)
	_ = flag.NewFlagSet("set", flag.ExitOnError)
	_ = flag.NewFlagSet("del", flag.ExitOnError)

	svc, err := smutil.NewSMClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	switch os.Args[1] {
	case "get":
		key := os.Args[2]
		ret, err := get(ctx, svc, key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ret)
	case "set":
		key := os.Args[2]
		value := os.Args[3]
		err = set(ctx, svc, key, value)
		if err != nil {
			log.Fatal(err)
		}
	// case "del":
	// 	err = del(svc)

	default:
		usage()
	}

	if err != nil {
		log.Fatal(err)
	}
}
