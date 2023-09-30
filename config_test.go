package config

import (
	"context"
	"reflect"
	"testing"
)

type NestedStruct struct {
	SubField string `ssm:"/sub_field"`
	Ignored  string
}

type Config struct {
	DatabaseURL string       `ssm:"/database_url"`
	APIKey      string       `ssm:"/api_key,required"`
	Debug       bool         `ssm:"/debug"`
	Nested      NestedStruct `ssm:"/nested"`
}

type ConfigEnv struct {
	DatabaseURL string `env:"DATABASE_URL"`
	APIKey      string `env:"API_KEY,required"`
	Debug       bool   `env:"DEBUG"`
}

func TestLoad(t *testing.T) {
	// Set up a mock SSM client
	client := newMockSSMClient("/myapp")
	// Set up the LoadOpts
	opts := []LoadOptFunc{
		WithPrefix("/myapp"),
		WithSSMClient(client),
	}

	// Load the config
	var cfg Config
	err := Load(context.Background(), &cfg, opts...)
	if err != nil {
		t.Fatalf("Load returned an error: %v", err)
	}

	// Check that the config was loaded correctly
	expected := Config{
		DatabaseURL: "postgres://user:password@localhost/myapp",
		APIKey:      "my-secret-key",
		Debug:       true,
		Nested: NestedStruct{
			SubField: "sub_field_value",
			Ignored:  "",
		},
	}
	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("Config was not loaded correctly. Expected: %#v, got: %#v", expected, cfg)
	}
}

func TestLoadEnv(t *testing.T) {

	newMockEnv()

	// Load the config
	var cfg ConfigEnv
	err := Load(context.Background(), &cfg)
	if err != nil {
		t.Fatalf("Load returned an error: %v", err)
	}

	// Check that the config was loaded correctly
	expected := ConfigEnv{
		DatabaseURL: "postgres://user:password@localhost/myapp",
		APIKey:      "my-secret-key",
		Debug:       true,
	}
	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("Config was not loaded correctly. Expected: %#v, got: %#v", expected, cfg)
	}
}
func TestGetSSMRecusriveTags(t *testing.T) {

	input := Config{}

	inputReflect := reflect.ValueOf(input)

	expected := []*pathConfig{
		{
			name:     "/database_url",
			required: false,
			provider: "ssm",
			value:    inputReflect.Field(0),
		},
		{
			name:     "/api_key",
			required: true,
			provider: "ssm",
			value:    inputReflect.Field(1),
		},
		{
			name:     "/debug",
			required: false,
			provider: "ssm",
			value:    inputReflect.Field(2),
		},
		{
			name:     "/nested/sub_field",
			required: false,
			provider: "ssm",
			value:    inputReflect.Field(3).Field(0),
		},
	}

	actual := getSSMRecursiveTags(inputReflect, "/")

	for i, a := range actual {
		field := expected[i]

		if a.name != field.name {
			t.Errorf("Expected name %s, got %s", field.name, a.name)
		}

		if a.required != field.required {
			t.Errorf("Expected required %t, got %t", field.required, a.required)
		}

		if a.provider != field.provider {
			t.Errorf("Expected provider %s, got %s", field.provider, a.provider)
		}

	}
}

func TestGetEnvRecusriveTags(t *testing.T) {

	input := Config{}

	inputReflect := reflect.ValueOf(input)

	expected := []*pathConfig{
		{
			name:     "/database_url",
			required: false,
			provider: "ssm",
			value:    inputReflect.Field(0),
		},
		{
			name:     "/api_key",
			required: true,
			provider: "ssm",
			value:    inputReflect.Field(1),
		},
		{
			name:     "/debug",
			required: false,
			provider: "ssm",
			value:    inputReflect.Field(2),
		},
		{
			name:     "/nested/sub_field",
			required: false,
			provider: "ssm",
			value:    inputReflect.Field(3).Field(0),
		},
	}

	actual := getSSMRecursiveTags(inputReflect, "/")

	for i, a := range actual {
		field := expected[i]

		if a.name != field.name {
			t.Errorf("Expected name %s, got %s", field.name, a.name)
		}

		if a.required != field.required {
			t.Errorf("Expected required %t, got %t", field.required, a.required)
		}

		if a.provider != field.provider {
			t.Errorf("Expected provider %s, got %s", field.provider, a.provider)
		}

	}
}
