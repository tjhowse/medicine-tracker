package main

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// This stores the medicine-logger database.
type MedicineLoggerDB struct {
	db *gorm.DB
}

type MedicineTypeDB struct {
	gorm.Model
	MedicineType
	User userGUID
}

// This stores a record of all medicine loags.
type MedicineLogEntryDB struct {
	gorm.Model
	MedicineLogEntry
	User userGUID
}

type userGUID string

type UsersDB struct {
	gorm.Model
	User         userGUID
	Email        string
	Password     string
	Settings     string
	OneTimeToken string
}

// Open the DB and migrate if required.
func (f *MedicineLoggerDB) Init(filename string) error {
	var err error
	f.db, err = gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		return err
	}

	// Migrate the schema
	f.db.AutoMigrate(&MedicineTypeDB{})
	f.db.AutoMigrate(&MedicineLogEntryDB{})
	f.db.AutoMigrate(&UsersDB{})

	return nil
}

func NewMedicineLoggerDB(filename string) (*MedicineLoggerDB, error) {
	fdb := &MedicineLoggerDB{}
	if err := fdb.Init(filename); err != nil {
		return fdb, err
	}

	return fdb, nil
}

// Close the DB.
func (f *MedicineLoggerDB) Close() error {
	sqlDB, err := f.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Add/update a weight
func (f *MedicineLoggerDB) AddMedicine(m MedicineTypeDB) error {

	// Check if the type already exists
	var medicine MedicineTypeDB
	log.Println("Checking for medicine: ", m.MedicineId)
	// This weird Limit(1).Find(... business is because .First(), which ostensibly does the same thing,
	// also outputs a "record not found" error to console, which is annoying.
	if err := f.db.Limit(1).Find(&medicine, "user = ? AND medicine_id = ?", m.User, m.MedicineId).Error; err != nil {
		// Doesn't exist, add it
		log.Println("Adding new medicine: ", m)
		if err := f.db.Create(&m).Error; err != nil {
			return err
		}
	} else {
		// Exists, update it
		log.Println("Updating medicine: ", m)
		medicine.MedicineType = m.MedicineType
		medicine.User = m.User
		if err := f.db.Save(&medicine).Error; err != nil {
			return err
		}
	}
	return nil
}

// Get all the available medicines
func (f *MedicineLoggerDB) GetMedicines(u userGUID) ([]MedicineTypeDB, error) {

	var medicines []MedicineTypeDB
	if err := f.db.Find(&medicines, "user = ?", u).Error; err != nil {
		return medicines, err
	}
	return medicines, nil
}

// Add a medicine log entry
func (f *MedicineLoggerDB) AddMedicineLog(u userGUID, m MedicineLogEntry) error {
	log.Println("Adding medicine log entry: ", m)
	return f.db.Create(&MedicineLogEntryDB{
		User:             u,
		MedicineLogEntry: m,
	}).Error
}

// Get all the medicine log entries
func (f *MedicineLoggerDB) GetMedicineLog(u userGUID, start, end time.Time) ([]MedicineLogEntry, error) {
	var logDBs []MedicineLogEntryDB
	var logs []MedicineLogEntry
	log.Println("Getting log entries for ", u, " between ", start, " and ", end)
	if err := f.db.Find(&logDBs, "user = ? AND time >= ? AND time <= ?", u, start, end).Error; err != nil {
		return logs, err
	}
	// Copy the internal DB ID into the external Log id
	log.Println("Got ", len(logDBs), " log entries")
	for i, l := range logDBs {
		logDBs[i].LogId = LogID(l.ID)
		logs = append(logs, logDBs[i].MedicineLogEntry)
	}
	log.Println("Returning ", len(logs), " log entries")
	return logs, nil
}
func (f *MedicineLoggerDB) DeleteMedicineLog(u userGUID, id uint) error {
	return f.db.Delete(&MedicineLogEntryDB{}, id).Error
}

func (f *MedicineLoggerDB) GetSettings(u userGUID) (UserSettings, error) {
	var user UsersDB
	if err := f.db.Find(&user, "user = ?", u).Error; err != nil {
		return UserSettings{}, err
	}
	var settings UserSettings
	if err := json.Unmarshal([]byte(user.Settings), &settings); err != nil {
		return UserSettings{}, err
	}
	return settings, nil
}

// Warning: This overwrites all settings. It should apply some intelligence to what
// settings get overwritten.
func (f *MedicineLoggerDB) UpdateSettings(u userGUID, settings UserSettings) error {
	var user UsersDB
	if err := f.db.Find(&user, "user = ?", u).Error; err != nil {
		return err
	}
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	user.Settings = string(settingsJSON)
	return f.db.Save(&user).Error
}

func GetDefaultUserSettings(name string) UserSettings {
	return UserSettings{
		Name: name,
	}
}

func GetDefaultUserSettingsJSON(name string) string {
	s := GetDefaultUserSettings(name)
	settingsJSON, _ := json.Marshal(s)
	return string(settingsJSON)
}

func GetDefaultMedicineTypes() []MedicineType {
	return []MedicineType{
		{
			Name:       "Paracetamol",
			Dose:       500,
			MedicineId: 0,
		},
		{
			Name:       "Oxycodone",
			Dose:       5,
			MedicineId: 1,
		},
		{
			Name:       "Tramadol",
			Dose:       50,
			MedicineId: 2,
		},
		{
			Name:       "Asprin",
			Dose:       300,
			MedicineId: 3,
		},
		{
			Name:       "Pantoprazole",
			Dose:       300,
			MedicineId: 4,
		},
	}
}

// This adds a user to the database. It accepts their username and password,
// and stores them securely in the database.
func (f *MedicineLoggerDB) AddUser(username, email, password string) error {
	// Check to see if the user already exists
	if extant, _, _ := f.ValidateUser(username, password); extant {
		return errors.New("user already exists")
	}
	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Create a new user
	user := UsersDB{
		User:     userGUID(username),
		Email:    email,
		Password: string(hash),
		Settings: GetDefaultUserSettingsJSON(username),
	}
	// Add default available weights
	for _, medicine := range GetDefaultMedicineTypes() {
		var db MedicineTypeDB
		db.User = user.User
		db.MedicineType = medicine
		if err := f.db.Save(&db).Error; err != nil {
			return err
		}
	}
	// Add the user to the database
	return f.db.Save(&user).Error
}

// TODO this should probably return a user object
func (f *MedicineLoggerDB) ValidateUser(username string, password string) (userExists, pwValid bool, err error) {
	var user UsersDB
	if err := f.db.Find(&user, "user = ?", userGUID(username)).Error; err != nil {
		return false, false, err
	}
	if user.User == "" {
		return false, false, nil
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
		return true, true, nil
	}
	return true, false, nil
}

// This deletes a user from all tables
func (f *MedicineLoggerDB) DeleteUser(username userGUID) error {
	var user UsersDB
	if err := f.db.Find(&user, "user = ?", username).Error; err != nil {
		return err
	}
	if err := f.db.Unscoped().Delete(&user).Error; err != nil {
		return err
	}
	var logEntries []MedicineLogEntryDB
	if err := f.db.Find(&logEntries, "user = ?", username).Error; err != nil {
		return err
	}
	if err := f.db.Unscoped().Delete(&logEntries).Error; err != nil {
		return err
	}
	var medicines []MedicineTypeDB
	if err := f.db.Find(&medicines, "user = ?", username).Error; err != nil {
		return err
	}
	if err := f.db.Unscoped().Delete(&medicines).Error; err != nil {
		return err
	}
	return nil
}

func (f *MedicineLoggerDB) ValidateOneTimeToken(token string) (userGUID, error) {
	var user UsersDB
	if err := f.db.Find(&user, "onetimetoken = ?", token).Error; err != nil {
		return "", err
	}
	return user.User, nil
}
