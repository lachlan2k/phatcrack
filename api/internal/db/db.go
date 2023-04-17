package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

// Same as gorm default, except uses uuid instead of uint
type UUIDBaseModel struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SimpleBaseModel struct {
	ID uint `gorm:"primarykey"`
}

var dbInstance *gorm.DB = nil

var ErrNotFound = gorm.ErrRecordNotFound

func Connect(dsn string) error {
	var err error
	dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to db: %v", err)
	}

	runMigrations()
	return upsertConfig()
}

func GetInstance() *gorm.DB {
	return dbInstance
}

func runMigrations() {
	instance := GetInstance()

	instance.AutoMigrate(&Agent{})

	instance.AutoMigrate(&Job{})
	instance.AutoMigrate(&JobRuntimeData{})
	instance.AutoMigrate(&JobCrackedHash{})

	instance.AutoMigrate(&Wordlist{})
	instance.AutoMigrate(&RuleFile{})

	instance.AutoMigrate(&PotfileEntry{})

	instance.AutoMigrate(&Project{})
	instance.AutoMigrate(&ProjectShare{})
	instance.AutoMigrate(&Hashlist{})
	instance.AutoMigrate(&HashlistHash{})
	instance.AutoMigrate(&Attack{})

	instance.AutoMigrate(&User{})
}

type pgJSONBArray[T interface{}] struct {
	arr  pq.GenericArray
	Data []datatypes.JSONType[T]
}

func (a *pgJSONBArray[T]) Scan(src interface{}) error {
	a.Data = make([]datatypes.JSONType[T], 0)
	a.arr.A = &a.Data
	return a.arr.Scan(src)
}

func (a *pgJSONBArray[T]) GormDataType() string {
	return "jsonb[]"
}
