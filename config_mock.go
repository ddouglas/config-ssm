package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func newMockEnv() {
	for k, v := range sharedParams {
		os.Setenv(strings.ToUpper(k), v)
	}
}

var sharedParams = map[string]string{
	"database_url":     "postgres://user:password@localhost/myapp",
	"api_key":          "my-secret-key",
	"debug":            "true",
	"nested/sub_field": "sub_field_value",
}

func newMockSSMClient(prefix string) *mockSSMClient {

	params := make(map[string]string, len(sharedParams))
	for k, v := range sharedParams {
		params[fmt.Sprintf("%s/%s", prefix, k)] = v
	}

	return &mockSSMClient{
		params: params,
	}

}

type mockSSMClient struct {
	params map[string]string
}

func (c *mockSSMClient) GetParameters(ctx context.Context, params *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error) {

	out := new(ssm.GetParametersOutput)

	out.InvalidParameters = make([]string, 0, len(params.Names))

	for _, name := range params.Names {
		_, ok := c.params[name]
		if !ok {
			out.InvalidParameters = append(out.InvalidParameters, name)
			continue
		}

		out.Parameters = append(out.Parameters, ssmTypes.Parameter{
			Name:  aws.String(name),
			Value: aws.String(c.params[name]),
		})

	}

	return out, nil

}
