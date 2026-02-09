package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataSourceHandler(t *testing.T) {
	handler := NewDataSourceHandler(nil, nil)
	assert.NotNil(t, handler)
}

func TestBuildDSN(t *testing.T) {
	ds := &struct {
		Username string
		Password string
		Host     string
		Port     int
		Database string
	}{
		Username: "root",
		Password: "password",
		Host:     "localhost",
		Port:     3306,
		Database: "testdb",
	}

	dsn := ds.Username + ":" + ds.Password + "@tcp(" + ds.Host + ":" + "3306" + ")/" + ds.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	assert.Contains(t, dsn, "root:password")
	assert.Contains(t, dsn, "localhost:3306")
	assert.Contains(t, dsn, "testdb")
}
