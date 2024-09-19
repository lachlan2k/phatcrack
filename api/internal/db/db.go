package db

import (
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lachlan2k/phatcrack/api/internal/roles"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
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

func HardDelete[T any](obj *T) error {
	return GetInstance().Unscoped().Delete(obj).Error
}

func SoftDelete[T any](obj *T) error {
	return GetInstance().Delete(obj).Error
}

func Save[T any](obj *T) error {
	return GetInstance().Save(obj).Error
}

func GetByID[T any](id string) (*T, error) {
	var res T
	err := GetInstance().First(&res, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func GetAll[T any]() ([]T, error) {
	res := []T{}
	err := GetInstance().Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

var dbInstance *gorm.DB = nil

var ErrNotFound = gorm.ErrRecordNotFound

func Connect(dsn string) error {
	l := gormlogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gormlogger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  gormlogger.Warn,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	})

	var err error
	dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: l,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	runMigrations()
	return seed()
}

func seed() error {
	var userCount int64
	err := GetInstance().Model(&User{}).Count(&userCount).Error
	if err != nil {
		return err
	}

	if userCount == 0 {
		_, err := RegisterUserWithCredentials("admin", "changeme", []string{roles.UserRoleAdmin, roles.UserRoleRequiresPasswordChange})
		if err != nil {
			return err
		}
	}

	return nil
}

func GetInstance() *gorm.DB {
	return dbInstance
}

func runMigrations() {
	instance := GetInstance()

	instance.AutoMigrate(&Agent{})
	instance.AutoMigrate(&AgentRegistrationKey{})

	instance.AutoMigrate(&Job{})
	instance.AutoMigrate(&JobRuntimeData{})

	instance.AutoMigrate(&Listfile{})

	instance.AutoMigrate(&PotfileEntry{})

	instance.AutoMigrate(&Project{})
	instance.AutoMigrate(&ProjectShare{})
	instance.AutoMigrate(&Hashlist{})
	instance.AutoMigrate(&HashlistHash{})
	instance.AutoMigrate(&Attack{})
	instance.AutoMigrate(&AttackTemplate{})
	instance.AutoMigrate(&AttackTemplateSet{})

	instance.AutoMigrate(&User{})

	instance.AutoMigrate(&Config{})
}

func WipeEverything() error {
	instance := GetInstance()

	toDelete := []interface{}{&Agent{}, &Job{}, &JobRuntimeData{}, &Listfile{}, &PotfileEntry{}, &Project{}, &ProjectShare{}, &Hashlist{}, &HashlistHash{}, &Attack{}, &User{}, &Config{}}

	return instance.Transaction(func(tx *gorm.DB) error {
		for _, d := range toDelete {
			err := tx.Unscoped().Where("1 = 1").Delete(d).Error
			if err != nil {
				return err
			}
		}

		return tx.Create(&User{
			Username:     "admin",
			Roles:        []string{roles.UserRoleAdmin, roles.UserRoleRequiresPasswordChange},
			PasswordHash: "$2a$10$6bms9eKxGvegFOd7XTA.XORrgn/ulqvWAVTcUGnXjVmvrR2O9/ViK",
		}).Error
	})
}

type pgJSONBArray[T interface{}] struct {
	arr  pq.GenericArray
	Data []datatypes.JSONType[T]
}

func (a *pgJSONBArray[T]) Init() {
	a.Data = make([]datatypes.JSONType[T], 0)
}

// Value implements the driver.Valuer interface.
func (a pgJSONBArray[T]) Value() (driver.Value, error) {
	a.arr.A = a.Data
	return a.arr.Value()
}

func (a *pgJSONBArray[T]) Unwrap() []T {
	arr := make([]T, len(a.Data))
	for i, el := range a.Data {
		arr[i] = el.Data()
	}

	return arr
}

func (a *pgJSONBArray[T]) Scan(src interface{}) error {
	a.Data = make([]datatypes.JSONType[T], 0)
	a.arr.A = &a.Data
	return a.arr.Scan(src)
}

func (a *pgJSONBArray[T]) GormDataType() string {
	return "jsonb[]"
}
