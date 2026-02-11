package datasource

import (
	"errors"
	"fmt"
)

type ConnectorProfile struct {
	Name        string
	DisplayName string
	Required    []string
	Optional    []string
	SupportsSSH bool
}

var ConnectorProfiles = map[string]ConnectorProfile{
	"mysql": {
		Name:        "mysql",
		DisplayName: "MySQL",
		Required:    []string{"host", "port", "database"},
		Optional:    []string{"username", "password"},
		SupportsSSH: true,
	},
	"postgres": {
		Name:        "postgres",
		DisplayName: "PostgreSQL",
		Required:    []string{"host", "port", "database"},
		Optional:    []string{"username", "password"},
		SupportsSSH: true,
	},
	"mongodb": {
		Name:        "mongodb",
		DisplayName: "MongoDB",
		Required:    []string{"host", "port", "database"},
		Optional:    []string{"username", "password"},
		SupportsSSH: true,
	},
	"excel": {
		Name:        "excel",
		DisplayName: "Excel",
		Required:    []string{"file_path"},
		Optional:    []string{},
		SupportsSSH: false,
	},
	"csv": {
		Name:        "csv",
		DisplayName: "CSV",
		Required:    []string{"file_path"},
		Optional:    []string{"delimiter"},
		SupportsSSH: false,
	},
	"api": {
		Name:        "api",
		DisplayName: "API",
		Required:    []string{"url"},
		Optional:    []string{"username", "password", "token"},
		SupportsSSH: false,
	},
}

type ProfileValidator struct {
	profiles map[string]ConnectorProfile
}

func NewProfileValidator() *ProfileValidator {
	return &ProfileValidator{
		profiles: ConnectorProfiles,
	}
}

func (v *ProfileValidator) Validate(dsType string, config map[string]interface{}) error {
	profile, exists := v.profiles[dsType]
	if !exists {
		return fmt.Errorf("unsupported datasource type: %s", dsType)
	}

	for _, field := range profile.Required {
		if _, ok := config[field]; !ok {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	if profile.SupportsSSH {
		if sshHost, ok := config["ssh_host"]; ok && sshHost != "" {
			if sshPort, ok := config["ssh_port"]; !ok || sshPort == nil || sshPort == 0 {
				return errors.New("ssh_port is required when ssh_host is provided")
			}
			if sshUsername, ok := config["ssh_username"]; !ok || sshUsername == nil || sshUsername == "" {
				return errors.New("ssh_username is required when ssh_host is provided")
			}
			if _, hasPassword := config["ssh_password"]; !hasPassword {
				if _, hasKey := config["ssh_key"]; !hasKey {
					return errors.New("either ssh_password or ssh_key is required for SSH tunnel")
				}
			}
		}
	}

	if maxConnections, ok := config["max_connections"]; ok {
		if maxConn, ok := maxConnections.(int); ok {
			if maxConn < 1 || maxConn > 100 {
				return fmt.Errorf("max_connections must be between 1 and 100, got: %d", maxConn)
			}
		}
	}

	if queryTimeout, ok := config["query_timeout_seconds"]; ok {
		if timeout, ok := queryTimeout.(int); ok {
			if timeout < 5 || timeout > 300 {
				return fmt.Errorf("query_timeout_seconds must be between 5 and 300, got: %d", timeout)
			}
		}
	}

	return nil
}

func (v *ProfileValidator) GetProfile(dsType string) (ConnectorProfile, error) {
	profile, exists := v.profiles[dsType]
	if !exists {
		return ConnectorProfile{}, fmt.Errorf("unsupported datasource type: %s", dsType)
	}
	return profile, nil
}

func (v *ProfileValidator) ListProfiles() []ConnectorProfile {
	profiles := make([]ConnectorProfile, 0, len(v.profiles))
	for _, profile := range v.profiles {
		profiles = append(profiles, profile)
	}
	return profiles
}
