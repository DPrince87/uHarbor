package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type AssetData struct {
	ID          uint `gorm:"primaryKey"`
	AssetID     uint
	AttributeID uint
	Value       float64
	Date        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Valid       int
	Asset       Assets                    `gorm:"foreignKey:AssetID"`
	Attribute   ServiceProviderAttributes `gorm:"foreignKey:AttributeID"`
}

type AssetDataResources struct {
	ID          uint `gorm:"primaryKey"`
	AssetDataID uint
	Name        string
	Value       float64
	Extra       string
	AssetData   AssetData `gorm:"foreignKey:AssetDataID"`
}

type Assets struct {
	ID       uint `gorm:"primaryKey"`
	SystemID uint
	CMDBID   uint
	Type     string
	Name     string
	Date     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Path     string
	Excludes string
	Active   int
	System   Systems `gorm:"foreignKey:SystemID"`
}

type Customers struct {
	ID          uint `gorm:"primaryKey"`
	BillingID   uint
	BillingCode string
	Name        string
	ManagedBy   string
	Date        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Status      string
}

type Locations struct {
	ID      string `gorm:"primaryKey"`
	Country string
	Zip     string
	State   string
	City    string
	Street  string
}

type ReportTemplates struct {
	Module  string `gorm:"primaryKey"`
	Enabled int
}

type Reports struct {
	ID           uint      `gorm:"primaryKey"`
	Date         time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ContentType  string
	ReportConfig string
	Report       []byte
	Name         string
	Status       string `gorm:"default:COMPLETE"`
}

type ServiceAsset struct {
	ServiceID    uint
	AssetID      uint
	SnapshotDate time.Time `gorm:"default:'0000-00-00 00:00:00'"`
	Frozen       int       `gorm:"default:0"`
}

type ServiceBillingPolicy struct {
	ServiceBillingPolicyID int `gorm:"primaryKey"`
	ServiceID              uint
	SofNumber              string
	BillingCommit          float64
	BillingBurst           float64
	BillingAttribute       string
	BillingMethod          string
	OverageIncrement       float64
	OveragePriceUSD        float64
	UnitOfMeasure          string
	ActivatedOn            time.Time `gorm:"default:'0000-00-00 00:00:00'"`
	Active                 int
}

type ServiceProviderAttributes struct {
	ID                uint `gorm:"primaryKey"`
	ServiceProviderID uint
	Attribute         string
	Unit              string
	AggregationPeriod string `gorm:"default:'Daily'"`
	Billable          int
	SystemOnly        int
	CapacityPlanning  int `gorm:"default:0"`
	Active            int
	ServiceProvider   ServiceProvider `gorm:"foreignKey:ServiceProviderID"`
}

type ServiceProvider struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	ServiceType string
	Plugin      string
	Date        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Active      int
}

type ServiceTypes struct {
	Type        string `gorm:"primaryKey"`
	Term        string
	Enabled     int
	Description string
}

type Services struct {
	ID               uint `gorm:"primaryKey"`
	BillingServiceId string
	CustomerID       uint
	Name             string
	SKU              string
	ServiceType      string
	Location         string
	StartBilling     time.Time
	EndBilling       time.Time
	Active           int
	Date             time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ServiceClass     string    `gorm:"default:'burstable'"`
	ParentID         uint
}

type Sessions struct {
	ID          string `gorm:"primaryKey"`
	SessionData string
	Expires     int
}

type SystemData struct {
	ID          uint `gorm:"primaryKey"`
	SystemID    uint
	AttributeID uint
	Value       float64
	Date        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Valid       int
	Attribute   ServiceProviderAttributes `gorm:"foreignKey:AttributeID"`
	System      Systems                   `gorm:"foreignKey:SystemID"`
}

type SystemDataResources struct {
	ID           uint `gorm:"primaryKey"`
	SystemDataID uint
	Name         string
	Value        float64
	Extra        string
	SystemData   SystemData `gorm:"foreignKey:SystemDataID"`
}

type Systems struct {
	ID                uint `gorm:"primaryKey"`
	ServiceProviderID uint
	CMDBID            string
	Name              string
	URI               string
	Parameters        string
	Username          string
	Password          string
	Date              time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CapacityPlanning  int       `gorm:"default:1"`
	Active            int
	ServiceProvider   ServiceProvider `gorm:"foreignKey:ServiceProviderID"`
}

type Task struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Type      string
	Email     string
	UserID    uint
	Date      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Object    string
	Enabled   int `gorm:"default:1"`
	Frequency string
	User      Users `gorm:"foreignKey:UserID"`
}

type Users struct {
	ID       uint `gorm:"primaryKey"`
	Username string
	Password string
	Date     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Role     string    `gorm:"default:'USER'"`
}

func (app *application) db() *gorm.DB {
	// if DB is defined, then return it
	if app.DB != nil {
		return app.DB
	}
	// Use SQLite as the database driver
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.Exec("PRAGMA foreign_keys = ON")

	if err != nil {
		panic("failed to connect to database")
	}

	// Auto Migrate the schema
	err = db.AutoMigrate(
		&AssetData{},
		&AssetDataResources{},
		&Assets{},
		&Customers{},
		&Locations{},
		&ReportTemplates{},
		&Reports{},
		&ServiceAsset{},
		&ServiceBillingPolicy{},
		&ServiceProviderAttributes{},
		&ServiceProvider{},
		&ServiceTypes{},
		&Services{},
		&Sessions{},
		&SystemData{},
		&SystemDataResources{},
		&Systems{},
		&Task{},
		&Users{},
	)
	if err != nil {
		fmt.Println(err)
	}

	app.DB = db

	return db
}
