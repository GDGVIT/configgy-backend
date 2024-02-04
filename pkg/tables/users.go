package tables

import (
	"time"

	"gorm.io/gorm"
)

// A struct on which the methods are defined
type DB struct {
	gormDB *gorm.DB
}

func NewDB(gormDB *gorm.DB) *DB {
	return &DB{gormDB: gormDB}
}

type Users struct {
	ID       int    `gorm:"column:user_id;primaryKey;autoIncrement"`
	PID      string `gorm:"column:user_pid;unique;type:varchar(100)"`
	Name     string `gorm:"column:user_name;not null;type:varchar(100)"`
	Email    string `gorm:"column:user_email;unique;type:varchar(100)"`
	Password []byte `gorm:"column:user_password;not null"`
	TOTPKey  []byte `gorm:"column:user_totp_key;"`
	// public key used to encrypt in shares
	PublicKey  []byte `gorm:"column:public_key;not null"`
	IsAdmin    bool   `gorm:"column:is_admin;not null;default:false"`
	IsEnabled  bool   `gorm:"column:is_enabled;not null;default:true"`
	IsVerified bool   `gorm:"column:is_verified;not null;default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UserVerification struct {
	ID        int    `gorm:"column:user_verification_id;primaryKey;autoIncrement"`
	UserID    int    `gorm:"column:user_id;not null"`
	Token     string `gorm:"column:token;not null"`
	Completed bool   `gorm:"column:user_verification_completed;not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Users) TableName() string {
	return "users"
}

// Create a new user
func (db *DB) CreateUser(user *Users) error {
	user.PID = UUIDWithPrefix("usr")
	return db.gormDB.Create(user).Error
}

// Get a user by id
func (db *DB) GetUserByID(id int) (*Users, error) {
	var user Users
	err := db.gormDB.Where("user_id = ?", id).First(&user).Error
	return &user, err
}

// Get a user by pid
func (db *DB) GetUserByPID(pid string) (*Users, error) {
	var user Users
	err := db.gormDB.Where("user_pid = ?", pid).First(&user).Error
	return &user, err
}

// Get user by email
func (db *DB) GetUserByEmail(email string) (*Users, error) {
	var user Users
	err := db.gormDB.Where("user_email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

// Create a new user verification
func (db *DB) CreateUserVerification(userVerification *UserVerification) error {
	return db.gormDB.Create(userVerification).Error
}

// Get a user verification by userID
func (db *DB) GetUserVerificationByUserID(userID int) (*UserVerification, error) {
	var userVerification UserVerification
	err := db.gormDB.Where("user_id = ?", userID).Where("user_verification_completed = ?", false).First(&userVerification).Error
	return &userVerification, err
}

func (db *DB) UpdateUserVerification(userID int, isVerified bool) error {
	return db.gormDB.Model(&UserVerification{}).Where("user_id = ?", userID).Update("user_verification_completed", isVerified).Error
}
