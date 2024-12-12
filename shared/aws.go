package shared

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func NewAWSCfg() (aws.Config, error) {

	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		log.Println(err)
		return aws.Config{}, err
	}

	return cfg, nil
}

func GetSecret(cfg aws.Config, secretName string) (string, error) {
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

	log.Printf("Retrieved %s from secretsmgr", secretName)

	// Decrypts secret using the associated KMS key.
	return *result.SecretString, nil

}
