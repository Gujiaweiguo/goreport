package datasource

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ConnectionBuilder struct{}

func NewConnectionBuilder() *ConnectionBuilder {
	return &ConnectionBuilder{}
}

func (b *ConnectionBuilder) BuildDSN(ctx context.Context, ds *models.DataSource) (string, *SSHTunnel, error) {
	var tunnel *SSHTunnel

	if ds.SSHHost != "" && ds.SSHPort > 0 {
		tunnelConfig := &SSHTunnelConfig{
			Host:     ds.SSHHost,
			Port:     ds.SSHPort,
			Username: ds.SSHUsername,
		}

		if ds.SSHPassword != "" {
			tunnelConfig.Password = ds.SSHPassword
		}

		if ds.SSHKey != "" {
			tunnelConfig.Key = []byte(ds.SSHKey)
			tunnelConfig.Phrase = ds.SSHKeyPhrase
		}

		tunnel = NewSSHTunnel(tunnelConfig)

		localAddr, err := tunnel.Connect(ctx, tunnelConfig, ds.Host, ds.Port)
		if err != nil {
			return "", nil, fmt.Errorf("failed to establish SSH tunnel: %w", err)
		}

		host, portStr, err := parseHostPort(localAddr)
		if err != nil {
			tunnel.Close()
			return "", nil, fmt.Errorf("failed to parse SSH tunnel local address: %w", err)
		}

		ds.Host = host
		ds.Port = portStr
	}

	port := strconv.Itoa(ds.Port)
	effectiveHost := ds.Host
	if tunnel == nil {
		effectiveHost = ResolveHost(ds.Host)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ds.Username,
		ds.Password,
		effectiveHost,
		port,
		ds.Database,
	)

	return dsn, tunnel, nil
}

func (b *ConnectionBuilder) Connect(ctx context.Context, ds *models.DataSource) (*gorm.DB, *SSHTunnel, error) {
	dsn, tunnel, err := b.BuildDSN(ctx, ds)
	if err != nil {
		return nil, nil, err
	}

	var db *gorm.DB
	var dialErr error

	dialCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	dialChan := make(chan *gorm.DB, 1)
	errChan := make(chan error, 1)

	go func() {
		openedDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
			NowFunc: func() time.Time {
				return time.Now().Local()
			},
		})
		if err != nil {
			errChan <- err
			return
		}
		dialChan <- openedDB
	}()

	select {
	case <-dialCtx.Done():
		if tunnel != nil {
			tunnel.Close()
		}
		return nil, nil, fmt.Errorf("database connection timeout: %w", dialCtx.Err())
	case db = <-dialChan:
		if ds.QueryTimeoutSeconds > 0 {
			sqlDB, err := db.DB()
			if err != nil {
				if tunnel != nil {
					tunnel.Close()
				}
				return nil, nil, fmt.Errorf("failed to get underlying database connection: %w", err)
			}
			sqlDB.SetConnMaxLifetime(time.Duration(ds.QueryTimeoutSeconds) * time.Second)
		}

		if ds.MaxConnections > 0 {
			sqlDB, err := db.DB()
			if err != nil {
				if tunnel != nil {
					tunnel.Close()
				}
				return nil, nil, fmt.Errorf("failed to get underlying database connection: %w", err)
			}
			sqlDB.SetMaxOpenConns(ds.MaxConnections)
			sqlDB.SetMaxIdleConns(min(ds.MaxConnections/2, 5))
		}

		return db, tunnel, nil
	case dialErr = <-errChan:
		if tunnel != nil {
			tunnel.Close()
		}
		return nil, nil, fmt.Errorf("failed to connect to database: %w", dialErr)
	}
}

func (b *ConnectionBuilder) TestConnection(ctx context.Context, ds *models.DataSource) error {
	db, tunnel, err := b.Connect(ctx, ds)
	if err != nil {
		return err
	}
	defer func() {
		if tunnel != nil {
			tunnel.Close()
		}
	}()

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

func parseHostPort(addr string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port: %w", err)
	}

	return host, port, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
