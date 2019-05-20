package mysql

import (
	"database/sql"
	"fmt"
	"gconst"
	"lobbyserver/lobby"
	"strconv"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql" //不能去掉，不然连接数据库的时候提示找不到mssql
	log "github.com/sirupsen/logrus"
)

// DBConfig 数据库配置
type DBConfig struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	UserName string `json:"userName"`
	Password string `json:"password"`
	DBName   string `json:"dbName"`
}

func newDbConnect(ip string, port int, user string, password string, database string) (*sql.DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?", user, password, ip, port, database)

	log.Printf("mysql connString:%s\n", connString)
	dbCon, err := sql.Open("mysql", connString)
	if err != nil {
		log.Println("Open mssql connection failed:", err.Error())
		return nil, err
	}

	err = dbCon.Ping()
	if err != nil {
		log.Println(database, "Cannot ping: ", err.Error())
		return nil, err
	}

	return dbCon, nil
}

func loadDBConfigFromRedis() *DBConfig {
	conn := lobby.Pool().Get()
	defer conn.Close()

	fields, err := redis.Strings(conn.Do("HMGET", gconst.LobbyDatabaseConfig, "ip", "port", "userName", "password", "dbName"))
	if err != nil {
		log.Println("loadUserInfoFromRedis, error", err)
		return nil
	}

	dbConfig := &DBConfig{}
	dbConfig.IP = fields[0]
	dbConfig.Port, _ = strconv.Atoi(fields[1])
	dbConfig.UserName = fields[2]
	dbConfig.Password = fields[3]
	dbConfig.DBName = fields[4]

	return dbConfig
}

// StartMssql 启动mssql,只能调用一次
func startMySQL() (*sql.DB, error) {
	dbConfig := loadDBConfigFromRedis()

	if dbConfig.IP == "" || dbConfig.Port == 0 || dbConfig.UserName == "" || dbConfig.Password == "" || dbConfig.DBName == "" {
		log.Error("DBConfig not setting, Plase setting database ip, port, userName, password, dbName on redis")

		return nil, fmt.Errorf("DBConfig not setting, Plase setting database ip, port, userName, password, dbName on redis")
	}

	gameDBCon, err := newDbConnect(dbConfig.IP, dbConfig.Port, dbConfig.UserName, dbConfig.Password, dbConfig.DBName)

	return gameDBCon, err
}
