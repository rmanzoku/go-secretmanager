package sm

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

func List(ctx context.Context, svc *secretsmanager.Client) ([]string, error) {
	params := &secretsmanager.ListSecretsInput{MaxResults: 100}

	ret := []string{}
	for {
		out, err := svc.ListSecrets(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, s := range out.SecretList {
			ret = append(ret, *s.Name)
		}

		if out.NextToken == nil {
			break
		}

		params.NextToken = out.NextToken
	}

	return ret, nil
}

func Get(ctx context.Context, svc *secretsmanager.Client, key string) (string, error) {
	params := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}
	out, err := svc.GetSecretValue(ctx, params)
	if err != nil {
		return "", err
	}
	return *out.SecretString, nil
}

func Set(ctx context.Context, svc *secretsmanager.Client, key, value string) error {
	_, err := Get(ctx, svc, key)
	if err == nil {
		return update(ctx, svc, key, value)
	}

	var e *types.ResourceNotFoundException
	if errors.As(err, &e) {
		return create(ctx, svc, key, value)
	}

	return err
}

func Del(ctx context.Context, svc *secretsmanager.Client, key string) error {
	return delete(ctx, svc, key)
}

func delete(ctx context.Context, svc *secretsmanager.Client, key string) error {
	params := &secretsmanager.DeleteSecretInput{
		SecretId: aws.String(key),
	}
	_, err := svc.DeleteSecret(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func create(ctx context.Context, svc *secretsmanager.Client, key, value string) error {
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

func update(ctx context.Context, svc *secretsmanager.Client, key, value string) error {
	params := &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(key),
		SecretString: aws.String(value),
	}
	_, err := svc.UpdateSecret(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
