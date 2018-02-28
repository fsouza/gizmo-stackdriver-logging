// Copyright 2018 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logging

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func TestLoadConfigFromEnvironment(t *testing.T) {
	var tests = []struct {
		name           string
		envs           map[string]string
		expectedConfig Config
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setEnvs(test.envs)
			var cfg Config
			err := envconfig.Process("", &cfg)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(cfg, test.expectedConfig); diff != "" {
				t.Errorf("invalid config returned\nwant %#v\ngot  %#v\n\ndiff: %v", test.expectedConfig, cfg, diff)
			}
		})
	}
}

func TestLogger(t *testing.T) {
	cfg := Config{
		Level:                       "debug",
		EnvironmentVariables:        []string{"HOME"},
		SendToStackDriver:           true,
		StackDriverErrorServiceName: "logging-test",
	}
	logger, err := cfg.Logger()
	if err != nil {
		t.Fatal(err)
	}
	if logger.Level != logrus.DebugLevel {
		t.Errorf("wrong level\nwant %v\ngot  %v", logrus.DebugLevel, logger.Level)
	}
	if logger.Out != os.Stdout {
		t.Errorf("wrong log output\nshould send logs to stdout, it's sending to: %#v", logger.Out)
	}
}

func TestLoggerInvalidLevel(t *testing.T) {
	cfg := Config{Level: "bananas"}
	logger, err := cfg.Logger()
	if err == nil {
		t.Error("unexpected <nil> error")
	}
	if logger != nil {
		t.Errorf("unexpected non-nil logger: %#v", logger)
	}
}

func setEnvs(envs map[string]string) {
	os.Clearenv()
	for name, value := range envs {
		os.Setenv(name, value)
	}
}
