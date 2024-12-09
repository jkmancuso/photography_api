package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const secretName = "salt"

func GetSalt(cfg aws.Config) (string, error) {
	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Println("Retrieved Salt from secretsmgr")

	// Decrypts secret using the associated KMS key.
	return *result.SecretString, nil

}
