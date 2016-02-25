package config

import (
	"djcontrol/funcs"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	API_HOST   = "http://localhost:8001"
	API_URL_01 = "/api/server/"
	API_URL_02 = "/api/project/"
	API_URL_03 = "/api/db/"
	API_URL_04 = "/api/install/"
	API_URL_05 = "/api/conf/"
	API_URL_06 = "/api/user/"

	API_URL_21 = "/api/server-conf/"
	API_URL_22 = "/api/postfix/"
	API_URL_23 = "/api/local-linux-username/"
	API_URL_24 = "/api/project-by-name/"
	API_URL_25 = "/api/project-conf/"
)

func GetConfig() *Config {
	c := new(Config)
	c.LocalLinuxUser = getLocalLinuxUsername()
	c.Localhost = getLocalhost()
	return c
}

func GetServer(code string) (*Server, error) {
	var s Server
	var url = API_HOST + API_URL_01 + code + "/"
	b := urlRequest(url)

	var err error
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, errors.New("Server not found.")
	}
	return &s, nil
}

func PrepareServerToInstall(s *Server) {
	s.Installs = getServerInstalls(s.Id)

	if funcs.IsSliceContainsString(s.Installs, "mysql") {
		s.MySQL = getMySQLDb(s.Id)
	}

	if funcs.IsSliceContainsString(s.Installs, "postgresql") {
		s.PostgreSQL = getPostgreSQLDb(s.Id)
	}

	if funcs.IsSliceContainsString(s.Installs, "postgresql") {
		s.Postfix = getPostfix(s.Id)
	}
}

func GetServerConf(serverId int, item string, filename string) *ServerConf {
	url := fmt.Sprintf("%s%s%d/%s/%s", API_HOST, API_URL_21, serverId, item, filename)
	b := urlRequest(url)

	var conf ServerConf

	var err error
	err = json.Unmarshal(b, &conf)
	checkErr(err, "GetServerConf")
	return &conf
}

func getSettingsValue(url string) string {
	b := urlRequest(url)

	res := map[string]string{}
	var err error
	err = json.Unmarshal(b, &res)
	checkErr(err, "getSettingsValue")
	return res["value"]
}

func getLocalLinuxUsername() string {
	var url = API_HOST + API_URL_23
	return getSettingsValue(url)
}

func getLocalhost() *Server {
	s, err := GetServer("localhost")
	checkErr(err, "getLocalhost")
	s.MySQL = getMySQLDb(s.Id)
	s.PostgreSQL = getPostgreSQLDb(s.Id)
	return s
}

func getInstalls() []*Install {
	var list Installs
	var url = API_HOST + API_URL_04
	b := urlRequest(url)

	var err error
	err = json.Unmarshal(b, &list)
	checkErr(err, "getInstalls")

	return list
}

func getServerInstalls(serverId int) []string {
	var list []string
	for _, x := range getInstalls() {
		if x.ServerId == serverId {
			list = append(list, x.Item)
		}
	}
	return list
}

func getDatabases() []*Db {
	var list Databases
	var url = API_HOST + API_URL_03
	b := urlRequest(url)

	var err error
	err = json.Unmarshal(b, &list)
	checkErr(err, "getDatabases")
	return list
}

func getUsers() []*User {
	// without passwords
	var list Users
	var url = API_HOST + API_URL_06
	b := urlRequest(url)

	var err error
	err = json.Unmarshal(b, &list)
	checkErr(err, "getUsers")
	return list
}

func getUserItem(id int) *UserItem {
	var user UserItem
	var url = API_HOST + API_URL_06 + strconv.Itoa(id) + "/"
	b := urlRequest(url)

	var err error
	err = json.Unmarshal(b, &user)
	checkErr(err, "getUserItem")
	return &user
}

func getMySQLDb(serverId int) *Db {
	for _, x := range getDatabases() {
		if x.ServerId == serverId && x.Type == "S" && x.TypeDb == "M" {
			x.User = getMySQLRootUser(x.Id)
			return x
		}
	}
	panic("Server doesn't have MySQLDb.")
}

func getPostgreSQLDb(serverId int) *Db {
	for _, x := range getDatabases() {
		if x.ServerId == serverId && x.Type == "S" && x.TypeDb == "P" {
			x.User = getPostgreSQLPostgresUser(x.Id)
			return x
		}
	}
	panic("Server doesn't have PostgreSQLDb.")
}

func getMySQLRootUser(dbId int) *UserItem {
	for _, x := range getUsers() {
		if x.Type == "db" && x.DbId == dbId && x.Name == "root" {
			return getUserItem(x.Id)
		}
	}
	panic("MySQL root user not found for server.")
}

func getPostgreSQLPostgresUser(dbId int) *UserItem {
	for _, x := range getUsers() {
		if x.Type == "db" && x.DbId == dbId && x.Name == "postgres" {
			return getUserItem(x.Id)
		}
	}
	panic("PostgreSQL postgres user not found for server.")
}

func getPostfix(serverId int) *Postfix {
	var p Postfix
	var url = fmt.Sprintf("%s%s%d/", API_HOST, API_URL_22, serverId)
	b := urlRequest(url)

	var err error
	err = json.Unmarshal(b, &p)
	checkErr(err, "getPostfix")

	return &p
}

func InputServer() (*Server, error) {
	fmt.Print("Input server name> ")

	var input string
	fmt.Scanln(&input)

	server, err := GetServer(input)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func GetProject(name string) (*Project, error) {
	var p Project
	var url = API_HOST + API_URL_24 + name + "/"
	b := urlRequest(url)

	var err error
	err = json.Unmarshal(b, &p)
	if err != nil {
		return nil, errors.New("Project not found.")
	}
	return &p, nil
}

func getProjectConfs(projectId int, item string) []*ProjectConf {
	var list ProjectConfs
	var url = fmt.Sprintf("%s%s%d/%s/", API_HOST, API_URL_25, projectId, item)
	b := urlRequest(url)

	var err error
	err = json.Unmarshal(b, &list)
	checkErr(err, "getProjectConfs")
	return list
}

func InputDjangoProject() (*Project, error) {
	fmt.Print("Input django project> ")

	var input string
	fmt.Scanln(&input)

	project, err := GetProject(input)
	if err != nil {
		return nil, err
	}

	if project.Type != "DJ" {
		return nil, errors.New("Incorrect django project.")
	}

	// initially project has server with small info
	// get server with full info
	server, err := GetServer(project.Server.Code)
	if err != nil {
		return nil, err
	}
	project.Server = server

	project.NginxConfs = getProjectConfs(project.Id, "nginx")
	project.SupervisorConfs = getProjectConfs(project.Id, "supervisor")

	return project, nil
}

func urlRequest(url string) []byte {
	fn := "urlRequest"
	req, err := http.NewRequest("GET", url, nil)
	checkErr(err, fn)

	client := &http.Client{}
	resp, err := client.Do(req)
	checkErr(err, fn)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err, fn)
	return body
}

func checkErr(err error, funcName string) {
	if err != nil {
		log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
		message := err.Error() + " " + funcName
		log.Fatal(errors.New(message))
	}
}
