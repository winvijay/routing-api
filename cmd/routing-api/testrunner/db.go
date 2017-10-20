package testrunner

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"code.cloudfoundry.org/routing-api/config"
	"github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/onsi/ginkgo"
)

type DbAllocator interface {
	Create() (*config.SqlDB, error)
	Reset() error
	Delete() error
	minConfig() *config.SqlDB
}

type mysqlAllocator struct {
	sqlDB      *sql.DB
	schemaName string
}

type postgresAllocator struct {
	sqlDB      *sql.DB
	schemaName string
}

func randSchemaName() string {
	return fmt.Sprintf("test%d%d", time.Now().UnixNano(), GinkgoParallelNode())
}

func NewPostgresAllocator() DbAllocator {
	return &postgresAllocator{schemaName: randSchemaName()}
}

func (a *postgresAllocator) minConfig() *config.SqlDB {
	return &config.SqlDB{
		Username:          "postgres",
		Password:          "",
		Host:              "localhost",
		Port:              5432,
		Type:              "postgres",
		CACert:            os.Getenv("POSTGRES_SERVER_CA_CERT"),
		SkipSSLValidation: os.Getenv("DB_SKIP_SSL_VALIDATION") == "true",
	}
}

func (a *postgresAllocator) connectionString(cfg *config.SqlDB) (string, error) {
	var queryString string
	if cfg.SkipSSLValidation {
		queryString = "?sslmode=require"
	} else if cfg.CACert != "" {
		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			return "", err
		}
		certPath := filepath.Join(tempDir, "postgres_cert.pem")
		err = ioutil.WriteFile(certPath, []byte(cfg.CACert), 0400)
		if err != nil {
			return "", err
		}
		queryString = fmt.Sprintf("?sslmode=verify-full&sslrootcert=%s", certPath)
	} else {
		queryString = "?sslmode=disable"
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Schema,
		queryString,
	), nil
}

func (a *postgresAllocator) Create() (*config.SqlDB, error) {
	var (
		err error
		cfg *config.SqlDB
	)

	cfg = a.minConfig()
	connStr, err := a.connectionString(cfg)
	if err != nil {
		return nil, err
	}
	a.sqlDB, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = a.sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	for i := 0; i < 5; i++ {
		dbExists, err := a.sqlDB.Exec(fmt.Sprintf("SELECT * FROM pg_database WHERE datname='%s'", a.schemaName))
		rowsAffected, err := dbExists.RowsAffected()
		if err != nil {
			return nil, err
		}
		if rowsAffected == 0 {
			_, err = a.sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", a.schemaName))
			if err != nil {
				return nil, err
			}
			cfg.Schema = a.schemaName
			return cfg, nil
		} else {
			a.schemaName = randSchemaName()
		}
	}
	return nil, errors.New("Failed to create unique database ")
}

func (a *postgresAllocator) Reset() error {
	_, err := a.sqlDB.Exec(fmt.Sprintf(`SELECT pg_terminate_backend(pid) FROM pg_stat_activity
	WHERE datname = '%s'`, a.schemaName))
	_, err = a.sqlDB.Exec(fmt.Sprintf("DROP DATABASE %s", a.schemaName))
	if err != nil {
		return err
	}

	_, err = a.sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", a.schemaName))
	return err
}

func (a *postgresAllocator) Delete() error {
	defer func() {
		_ = a.sqlDB.Close()
	}()
	_, err := a.sqlDB.Exec(fmt.Sprintf(`SELECT pg_terminate_backend(pid) FROM pg_stat_activity
	WHERE datname = '%s'`, a.schemaName))
	if err != nil {
		return err
	}
	_, err = a.sqlDB.Exec(fmt.Sprintf("DROP DATABASE %s", a.schemaName))
	return err
}

func NewMySQLAllocator() DbAllocator {
	return &mysqlAllocator{schemaName: randSchemaName()}
}

func (a *mysqlAllocator) minConfig() *config.SqlDB {
	return &config.SqlDB{
		Username:          "root",
		Password:          "password",
		Host:              "localhost",
		Port:              3306,
		Type:              "mysql",
		CACert:            os.Getenv("MYSQL_SERVER_CA_CERT"),
		SkipSSLValidation: os.Getenv("DB_SKIP_SSL_VALIDATION") == "true",
	}
}

func (a *mysqlAllocator) connectionString(cfg *config.SqlDB) string {
	rootCA := x509.NewCertPool()
	queryString := "?parseTime=true"
	configKey := "dbTLSKey"
	if cfg.SkipSSLValidation {
		tlsConfig := tls.Config{}
		tlsConfig.InsecureSkipVerify = cfg.SkipSSLValidation
		mysql.RegisterTLSConfig(configKey, &tlsConfig)
		queryString = fmt.Sprintf("%s&tls=%s", queryString, configKey)
	} else if cfg.CACert != "" {
		tlsConfig := tls.Config{}
		rootCA.AppendCertsFromPEM([]byte(cfg.CACert))
		tlsConfig.ServerName = cfg.Host
		tlsConfig.RootCAs = rootCA
		mysql.RegisterTLSConfig(configKey, &tlsConfig)
		queryString = fmt.Sprintf("%s&tls=%s", queryString, configKey)
	}
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Schema,
		queryString,
	)
}

func (a *mysqlAllocator) Create() (*config.SqlDB, error) {
	var (
		err error
		cfg *config.SqlDB
	)

	cfg = a.minConfig()
	a.sqlDB, err = sql.Open("mysql", a.connectionString(cfg))
	if err != nil {
		return nil, err
	}
	err = a.sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	for i := 0; i < 5; i++ {
		dbExists, err := a.sqlDB.Exec(fmt.Sprintf("SHOW DATABASES LIKE '%s'", a.schemaName))
		rowsAffected, err := dbExists.RowsAffected()
		if err != nil {
			return nil, err
		}
		if rowsAffected == 0 {
			_, err = a.sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", a.schemaName))
			if err != nil {
				return nil, err
			}
			cfg.Schema = a.schemaName
			return cfg, nil
		} else {
			a.schemaName = randSchemaName()
		}
	}
	return nil, errors.New("Failed to create unique database ")
}

func (a *mysqlAllocator) Reset() error {
	_, err := a.sqlDB.Exec(fmt.Sprintf("DROP DATABASE %s", a.schemaName))
	if err != nil {
		return err
	}

	_, err = a.sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", a.schemaName))
	return err
}

func (a *mysqlAllocator) Delete() error {
	defer func() {
		_ = a.sqlDB.Close()
	}()
	_, err := a.sqlDB.Exec(fmt.Sprintf("DROP DATABASE %s", a.schemaName))
	return err
}
