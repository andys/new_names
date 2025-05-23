package config

import (
	"os"
	"testing"

	"github.com/frankban/quicktest"
)

func TestLoadConfig_ParsesFieldsCorrectly(t *testing.T) {
	c := quicktest.New(t)
	// Create a temporary config file
	content := `
anonymize:
  users: email, name, phone
  orders: address
skip:
  - logs
  - audit
sample:
  users: 0.1
  orders: 0.5
`
	tmpfile, err := os.CreateTemp("", "testconfig*.conf")
	c.Assert(err, quicktest.IsNil)
	defer os.Remove(tmpfile.Name())
	_, err = tmpfile.WriteString(content)
	c.Assert(err, quicktest.IsNil)
	tmpfile.Close()

	cfg := &Config{}
	err = LoadConfig(cfg, tmpfile.Name())
	c.Assert(err, quicktest.IsNil)
	c.Assert(cfg.AnonymizeFields, quicktest.DeepEquals, map[string][]string{
		"users":  {"email", "name", "phone"},
		"orders": {"address"},
	})
	c.Assert(cfg.SkipTables, quicktest.DeepEquals, []string{"logs", "audit"})
	c.Assert(cfg.SampleTables, quicktest.DeepEquals, map[string]float64{
		"users":  0.1,
		"orders": 0.5,
	})
}

func TestLoadConfig_HandlesEmptyFile(t *testing.T) {
	c := quicktest.New(t)
	tmpfile, err := os.CreateTemp("", "testconfig*.conf")
	c.Assert(err, quicktest.IsNil)
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	cfg := &Config{}
	err = LoadConfig(cfg, tmpfile.Name())
	c.Assert(err, quicktest.IsNil)
	c.Assert(cfg.AnonymizeFields, quicktest.DeepEquals, map[string][]string{})
}
