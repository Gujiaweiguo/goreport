package datasource

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProfileValidator(t *testing.T) {
	validator := NewProfileValidator()
	assert.NotNil(t, validator)
}

func TestProfileValidator_Validate(t *testing.T) {
	validator := NewProfileValidator()

	tests := []struct {
		name    string
		dsType  string
		config  map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name:   "valid mysql config",
			dsType: "mysql",
			config: map[string]interface{}{
				"host":     "localhost",
				"port":     3306,
				"database": "test",
			},
			wantErr: false,
		},
		{
			name:   "missing required field",
			dsType: "mysql",
			config: map[string]interface{}{
				"host": "localhost",
			},
			wantErr: true,
			errMsg:  "missing required field",
		},
		{
			name:    "unsupported datasource type",
			dsType:  "unsupported",
			config:  map[string]interface{}{},
			wantErr: true,
			errMsg:  "unsupported datasource type",
		},
		{
			name:   "valid mysql with SSH",
			dsType: "mysql",
			config: map[string]interface{}{
				"host":         "localhost",
				"port":         3306,
				"database":     "test",
				"ssh_host":     "ssh.example.com",
				"ssh_port":     22,
				"ssh_username": "user",
				"ssh_password": "pass",
			},
			wantErr: false,
		},
		{
			name:   "SSH without port",
			dsType: "mysql",
			config: map[string]interface{}{
				"host":         "localhost",
				"port":         3306,
				"database":     "test",
				"ssh_host":     "ssh.example.com",
				"ssh_username": "user",
			},
			wantErr: true,
			errMsg:  "ssh_port is required",
		},
		{
			name:   "SSH without username",
			dsType: "mysql",
			config: map[string]interface{}{
				"host":     "localhost",
				"port":     3306,
				"database": "test",
				"ssh_host": "ssh.example.com",
				"ssh_port": 22,
			},
			wantErr: true,
			errMsg:  "ssh_username is required",
		},
		{
			name:   "SSH without password or key",
			dsType: "mysql",
			config: map[string]interface{}{
				"host":         "localhost",
				"port":         3306,
				"database":     "test",
				"ssh_host":     "ssh.example.com",
				"ssh_port":     22,
				"ssh_username": "user",
			},
			wantErr: true,
			errMsg:  "either ssh_password or ssh_key is required",
		},
		{
			name:   "valid postgres config",
			dsType: "postgres",
			config: map[string]interface{}{
				"host":     "localhost",
				"port":     5432,
				"database": "test",
			},
			wantErr: false,
		},
		{
			name:   "valid api config",
			dsType: "api",
			config: map[string]interface{}{
				"url": "https://api.example.com",
			},
			wantErr: false,
		},
		{
			name:   "valid csv config",
			dsType: "csv",
			config: map[string]interface{}{
				"file_path": "/data/test.csv",
			},
			wantErr: false,
		},
		{
			name:   "valid excel config",
			dsType: "excel",
			config: map[string]interface{}{
				"file_path": "/data/test.xlsx",
			},
			wantErr: false,
		},
		{
			name:   "max_connections too low",
			dsType: "mysql",
			config: map[string]interface{}{
				"host":            "localhost",
				"port":            3306,
				"database":        "test",
				"max_connections": 0,
			},
			wantErr: true,
			errMsg:  "max_connections must be between 1 and 100",
		},
		{
			name:   "max_connections too high",
			dsType: "mysql",
			config: map[string]interface{}{
				"host":            "localhost",
				"port":            3306,
				"database":        "test",
				"max_connections": 101,
			},
			wantErr: true,
			errMsg:  "max_connections must be between 1 and 100",
		},
		{
			name:   "query_timeout too low",
			dsType: "mysql",
			config: map[string]interface{}{
				"host":                  "localhost",
				"port":                  3306,
				"database":              "test",
				"query_timeout_seconds": 4,
			},
			wantErr: true,
			errMsg:  "query_timeout_seconds must be between 5 and 300",
		},
		{
			name:   "query_timeout too high",
			dsType: "mysql",
			config: map[string]interface{}{
				"host":                  "localhost",
				"port":                  3306,
				"database":              "test",
				"query_timeout_seconds": 301,
			},
			wantErr: true,
			errMsg:  "query_timeout_seconds must be between 5 and 300",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.dsType, tt.config)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProfileValidator_GetProfile(t *testing.T) {
	validator := NewProfileValidator()

	t.Run("existing profile", func(t *testing.T) {
		profile, err := validator.GetProfile("mysql")
		require.NoError(t, err)
		assert.Equal(t, "mysql", profile.Name)
		assert.Equal(t, "MySQL", profile.DisplayName)
		assert.True(t, profile.SupportsSSH)
	})

	t.Run("non-existing profile", func(t *testing.T) {
		_, err := validator.GetProfile("nonexistent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported datasource type")
	})
}

func TestProfileValidator_ListProfiles(t *testing.T) {
	validator := NewProfileValidator()

	profiles := validator.ListProfiles()
	assert.NotEmpty(t, profiles)

	profileNames := make(map[string]bool)
	for _, p := range profiles {
		profileNames[p.Name] = true
	}

	assert.True(t, profileNames["mysql"])
	assert.True(t, profileNames["postgres"])
	assert.True(t, profileNames["api"])
	assert.True(t, profileNames["csv"])
	assert.True(t, profileNames["excel"])
}

func TestConnectorProfiles(t *testing.T) {
	t.Run("mysql profile", func(t *testing.T) {
		profile := ConnectorProfiles["mysql"]
		assert.Equal(t, "MySQL", profile.DisplayName)
		assert.Contains(t, profile.Required, "host")
		assert.Contains(t, profile.Required, "port")
		assert.Contains(t, profile.Required, "database")
		assert.True(t, profile.SupportsSSH)
	})

	t.Run("api profile", func(t *testing.T) {
		profile := ConnectorProfiles["api"]
		assert.Equal(t, "API", profile.DisplayName)
		assert.Contains(t, profile.Required, "url")
		assert.False(t, profile.SupportsSSH)
	})

	t.Run("csv profile", func(t *testing.T) {
		profile := ConnectorProfiles["csv"]
		assert.Equal(t, "CSV", profile.DisplayName)
		assert.Contains(t, profile.Required, "file_path")
		assert.Contains(t, profile.Optional, "delimiter")
		assert.False(t, profile.SupportsSSH)
	})
}
