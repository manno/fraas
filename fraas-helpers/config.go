package fraas_helpers

import (
	"errors"
	"log"

	"github.com/gobuffalo/envy"
	"gopkg.in/yaml.v2"
)

type SiteConfig struct {
	Database Database `yaml:"database"`
	Google   Google   `yaml:"google"`
	Mail     Mail     `yaml:"mail"`

	DockerImage string `yaml:"docker_image"`
	Domain      string `yaml:"domain"`
}

type Database struct {
	Instance      string `yaml:"instance"`
	AdminUser     string `yaml:"admin_user"`
	AdminPassword string `yaml:"admin_password"`
}

type Google struct {
	ProjectID string `yaml:"project_id"`
	Zone      string `yaml:"zone"`
	Region    string `yaml:"region"`
	ClusterID string `yaml:"cluster_id"`
	DNSZone   string `yaml:"dnszone"`
}

type Mail struct {
	SMTPServer     string `yaml:"smtp_server"`
	SMTPServerPort string `yaml:"smtp_server_port"`
	SMTPNOTLS      string `yaml:"smtp_notls"`
	SMTPUsername   string `yaml:"smtp_user_name"`
	SMTPPassword   string `yaml:"smtp_password"`
	ExceptionEMail string `yaml:"exception_email"`
}

var ENV = envy.Get("FRAAS_CONFIG", "{}")
var config *SiteConfig

func Config() *SiteConfig {
	if config == nil {
		if err := ConfigFromEnv(); err != nil {
			log.Fatal(err)
		}
	}
	return config
}

func ConfigFromEnv() error {
	if err := yaml.Unmarshal([]byte(ENV), &config); err != nil {
		return err
	}
	if config.Domain == "" {
		return errors.New("FRAAS site config is missing values")
	}
	return nil
}
