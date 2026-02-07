package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gujiaweiguo/goreport/internal/cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CachedMetadataService struct {
	cache *cache.Cache
}

func NewCachedMetadataService(cache *cache.Cache) *CachedMetadataService {
	return &CachedMetadataService{
		cache: cache,
	}
}

func (s *CachedMetadataService) GetTables(ctx context.Context, tenantID, datasourceID, database string) ([]string, error) {
	domain := "datasource:tables"
	identity := datasourceID

	if cached, hit, err := s.cache.Get(ctx, tenantID, domain, identity, nil); err == nil && hit {
		var result []string
		if err := json.Unmarshal(cached, &result); err == nil {
			return result, nil
		}
	}

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s?charset=utf8mb4&parseTime=True&loc=Local", database)), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	tables, err := GetTables(ctx, db, database)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(tables)
	_ = s.cache.Set(ctx, tenantID, domain, identity, nil, data, 30*time.Minute)

	return tables, nil
}

func (s *CachedMetadataService) GetFields(ctx context.Context, tenantID, datasourceID, database, tableName string) ([]FieldInfo, error) {
	domain := "datasource:fields"
	identity := datasourceID + ":" + tableName

	if cached, hit, err := s.cache.Get(ctx, tenantID, domain, identity, nil); err == nil && hit {
		var result []FieldInfo
		if err := json.Unmarshal(cached, &result); err == nil {
			return result, nil
		}
	}

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s?charset=utf8mb4&parseTime=True&loc=Local", database)), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	fields, err := GetFields(ctx, db, database, tableName)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(fields)
	_ = s.cache.Set(ctx, tenantID, domain, identity, nil, data, 30*time.Minute)

	return fields, nil
}
