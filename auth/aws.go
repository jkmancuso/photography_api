package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func NewAWSCfg() (aws.Config, error) {

	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		log.Println(err)
		return aws.Config{}, err
	}

	return cfg, nil
}
