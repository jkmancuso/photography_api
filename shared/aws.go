package shared

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	log "github.com/sirupsen/logrus"
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

	log.Debugf("Getting secret %s from secretsmgr", secretName)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)

	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Debugf("Retrieved secret %s from secretsmgr", secretName)

	return *result.SecretString, nil

}
