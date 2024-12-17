package mysqlrepo

import (
	"database/sql"


	"github.com/jinzhu/gorm"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
)

// MySQLRepo mysql数据库实现
type MySQLRepo struct {
	log logger.Logger
	db  *gorm.DB
}

// NewRepo 创建mysql数据实现实例
func NewRepo(conf *config.Config, log logger.Logger) (repo *MySQLRepo, err error) {
	connection := conf.Repo.Connection

	db, err := gorm.Open("mysql", connection)
	if err != nil {
		log.Errorf("database connection failed:%s", err.Error())
		return nil, err
	}
	if conf.Repo.Debug {
		db.LogMode(true)
		if conf.Repo.LogDestination != "" && conf.Repo.LogDestination != config.ConsoleLog {
			db.SetLogger(logger.NewBeeLogger(&config.Logger{
				Level:    "debug",
				LogFile:  conf.Repo.LogDestination,
				FilePerm: conf.Logger.FilePerm,
			}))
		}
	}
	return &MySQLRepo{
		log: log,
		db:  db,
	}, nil
}

// NewRepoWithDB 创建新实例
func NewRepoWithDB(log logger.Logger, rawDB *sql.DB, debug bool) (*MySQLRepo, error) {
	db, err := gorm.Open("mysql", rawDB)
	if err != nil {
		log.Errorf("database connection failed:%s", err.Error())
		return nil, err
	}
	if debug {
		db.LogMode(true)
	}
	return &MySQLRepo{
		// conf: conf,
		log: log,
		db:  db,
	}, nil
}

// Close 关闭mysql连接
func (repo *MySQLRepo) Close() error {
	return repo.db.Close()
}

// DropDB 删除表
func (repo *MySQLRepo) DropDB() error {
	return nil
}
