package main

import (
	"fmt"

	nebula "github.com/vesoft-inc/nebula-go/v2"
)

const (
	address = "10.134.*.**"
	// 3699 is only for testing.
	port     = 9669
	username = "root"
	password = "******"
)

// Initialize logger
var log = nebula.DefaultLogger{}

func main() {
	hostAddress := nebula.HostAddress{Host: address, Port: port}
	hostList := []nebula.HostAddress{hostAddress}
	testPoolConfig := nebula.GetDefaultConf()

	pool, err := nebula.NewConnectionPool(hostList, testPoolConfig, log)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to initialize the connection pool, host: %s, port: %d, %s", address, port, err.Error()))
	}
	defer pool.Close()
	session, err := pool.GetSession(username, password)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to create a new session from connection pool, username: %s, password: %s, %s",
			username, password, err.Error()))
	}
	defer session.Release()

	checkResultSet := func(prefix string, res *nebula.ResultSet) {
		if !res.IsSucceed() {
			log.Fatal(fmt.Sprintf("%s, ErrorCode: %v, ErrorMsg: %s", prefix, res.GetErrorCode(), res.GetErrorMsg()))
		}
	}

	{
		// 用户信息
		createSchema := "CREATE SPACE IF NOT EXISTS test_user(vid_type=FIXED_STRING(20)); " +
			"USE test_user;" +
			"CREATE TAG IF NOT EXISTS user(name string, age int, sex int, phone string, address string, ip string, os string);" +
			"CREATE EDGE IF NOT EXISTS relation(rel double)"

		// Excute a query
		resultSet, err := session.Execute(createSchema)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		checkResultSet(createSchema, resultSet)
	}

	{
		// 用户写入mock数据, sex: 0男， 1 女
		insertVertext := "INSERT VERTEX user(name, age, sex,  phone, address,  ip, os) VALUES '1':('xiaohua', 20, 1, '16693311132', '北京市', '10.134.15.120', 'ios');\n" +
			"INSERT VERTEX user(name, age, sex,  phone, address,  ip, os) VALUES '2':('xiaoming', 30, 0, '17693311132', '北京市', '10.134.15.122', 'android');\n" +
			"INSERT VERTEX user(name, age, sex,  phone, address,  ip, os) VALUES '3':('xiaolan', 40, 1, '12693311132', '山东', '10.124.15.120', 'ios');\n" +
			"INSERT VERTEX user(name, age, sex,  phone, address,  ip, os) VALUES '4':('xiaolv', 35, 0, '17693311832', '浙江', '10.124.15.122', 'android'); \n" +
			"INSERT VERTEX user(name, age, sex,  phone, address,  ip, os) VALUES '5':('xiaohong', 45, 1, '12312445551', '天津', '10.114.15.120', 'ios'); \n" +
			"INSERT VERTEX user(name, age, sex,  phone, address,  ip, os) VALUES '6':('xiaohuang', 22, 0, '17673311132', '北京市', '10.133.15.122', 'android');"
		resultSet, err := session.Execute(insertVertext)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		checkResultSet(insertVertext, resultSet)
	}

	{
		// 用户边建立
		insertEdgit := "INSERT EDGE relation(rel) VALUES '1' -> '2':(100); \n" +
			"INSERT EDGE relation(rel) VALUES '2' -> '3':(70);\n" +
			"INSERT EDGE relation(rel) VALUES '1' -> '3':(50); \n" +
			"INSERT EDGE relation(rel) VALUES '2' -> '4':(100); \n" +
			"INSERT EDGE relation(rel) VALUES '3' -> '4':(80); \n" +
			"INSERT EDGE relation(rel) VALUES '5' -> '6':(100); \n" +
			"INSERT EDGE relation(rel) VALUES '6' -> '1':(20); \n"
		resultSet, err := session.Execute(insertEdgit)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		checkResultSet(insertEdgit, resultSet)

	}

	{
		// query 查询用户之间的关系
		querySimple := " GO FROM '1', '2', '3', '4', '5', '6' OVER relation yield $$.user.name AS name, relation._dst AS rel, relation.rel as re;"
		resultSet, err := session.Execute(querySimple)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		checkResultSet(querySimple, resultSet)
	}

	{
		// query 查询2跳数据总和
		querySql := "GO 2 STEPS FROM '1' OVER relation | YIELD COUNT(*);"
		resultSet, err := session.Execute(querySql)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		checkResultSet(querySql, resultSet)
	}

	// Drop space
	{
		query := "DROP SPACE IF EXISTS test_user"
		// Send query
		resultSet, err := session.Execute(query)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		checkResultSet(query, resultSet)
	}

	log.Info("Nebula Go Client user info")
}
