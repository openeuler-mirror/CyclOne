package middleware

import (
	"fmt"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/config/jsonconf"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model/mysqlrepo"
	"idcos.io/cloudboot/utils"
)

func TestReflectRepo(t *testing.T) {
	conf := GetConf()
	repo, err := GetRepo(conf, GetLog(conf))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if repo == nil {
		t.Errorf("fail load repo")
	}

	respMap := GetRespVal(repo, "GetImageTemplateByID", []reflect.Value{reflect.ValueOf(1)})
	fmt.Println(utils.ToJsonString(respMap))
}

func GetConf() *config.Config {
	conf, err := jsonconf.New("../../../../doc/conf/cloudboot-server.conf").Load()
	if err != nil {
		fmt.Printf("load config file error,%s", err.Error())
		return nil
	}

	return conf
}

func GetLog(conf *config.Config) logger.Logger {
	return logger.NewBeeLogger(&conf.Logger)

}

func GetRepo(conf *config.Config, log logger.Logger) (repo *mysqlrepo.MySQLRepo, err error) {
	repo, err = mysqlrepo.NewRepo(conf, log)
	if err != nil {
		return nil, err
	}
	return
}
