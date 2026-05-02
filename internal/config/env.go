package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var envPrefix = "FREECODE"

func applyEnvVars(v *viper.Viper) error {
	envVars := []string{
		"FREECODE_SHELL",
		"FREECODE_LOG_LEVEL",
		"FREECODE_YOLO",
		"FREECODE_QUIET",
		"FREECODE_VERBOSE",
		"FREECODE_DEBUG",
		"FREECODE_COLOR",
		"FREECODE_THEME",
		"FREECODE_EDITOR",
		"FREECODE_PAGER",
		"FREECODE_TIMEOUT",
		"FREECODE_RETRIES",
		"FREECODE_SERVER_HOST",
		"FREECODE_SERVER_PORT",
		"FREECODE_MCP_HOST",
		"FREECODE_MCP_PORT",
		"FREECODE_WEBUI_HOST",
		"FREECODE_WEBUI_PORT",
		"FREECODE_AGENT_DEFAULT",
		"FREECODE_AGENT_THINKING",
		"FREECODE_AGENT_TIMEOUT",
		"FREECODE_AGENT_MAX_TURNS",
		"FREECODE_SESSION_DIR",
		"FREECODE_SESSION_HISTORY_SIZE",
		"FREECODE_HTTP_PROXY",
		"FREECODE_HTTPS_PROXY",
		"FREECODE_NO_PROXY",
	}

	for _, envVar := range envVars {
		if val := os.Getenv(envVar); val != "" {
			key := strings.ToLower(strings.TrimPrefix(envVar, envPrefix+"_"))
			key = strings.ReplaceAll(key, "_", ".")
			v.Set(key, val)
		}
	}

	return nil
}

func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(envPrefix + "_" + key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(envPrefix + "_" + key); val != "" {
		return strings.ToLower(val) == "true" || val == "1"
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(envPrefix + "_" + key); val != "" {
		var result int
		if _, err := fmt.Sscanf(val, "%d", &result); err == nil {
			return result
		}
	}
	return defaultVal
}

func (c *Config) ApplyEnvOverrides() {
	c.Shell = getEnv("SHELL", c.Shell)
	c.LogLevel = getEnv("LOG_LEVEL", c.LogLevel)
	c.Yolo = getEnvBool("YOLO", c.Yolo)
	c.Quiet = getEnvBool("QUIET", c.Quiet)
	c.Verbose = getEnvBool("VERBOSE", c.Verbose)
	c.Debug = getEnvBool("DEBUG", c.Debug)
	c.Color = getEnvBool("COLOR", c.Color)
	c.Theme = getEnv("THEME", c.Theme)
	c.Editor = getEnv("EDITOR", c.Editor)
	c.Pager = getEnv("PAGER", c.Pager)
	c.Timeout = getEnvInt("TIMEOUT", c.Timeout)
	c.Retries = getEnvInt("RETRIES", c.Retries)

	if host := os.Getenv("FREECODE_SERVER_HOST"); host != "" {
		c.Server.Host = host
	}
	if port := os.Getenv("FREECODE_SERVER_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &c.Server.Port)
	}
	if agent := os.Getenv("FREECODE_AGENT_DEFAULT"); agent != "" {
		c.Agent.Default = agent
	}
}
