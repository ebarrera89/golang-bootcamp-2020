package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type Config struct {
	DSN      string
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	DB       *sql.DB
	Addr     *string
}

// Init the application's configuration
func Init(infoLog, errorLog *log.Logger) (*Config, error) {

	//Read config
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			panic(fmt.Errorf("fatal error config file not found: %s", err))

		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
	}

	// Get the server address/port
	addr := viper.GetString("addr")

	// Create the Data Source Name
	// *We need to use the parseTime=true parameter in our
	// DSN to force it to convert TIME and DATE fields to time.Time. Otherwise it returns these as
	// []byte objects.
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	dbHost := viper.GetString("db.host")
	dbPort := viper.GetString("db.port")
	dbName := viper.GetString("db.name")
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", username, password, dbHost, dbPort, dbName)
	fmt.Println(dsn)
	//Open db connection
	DB, err := openDB(dsn)
	if err != nil {
		errorLog.Fatalf("message: unable to open db connection, type: database, err: %v", err)
		return nil, err
	}

	return &Config{dsn, infoLog, errorLog, DB, &addr}, nil
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	//Check if the DB is responding
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}