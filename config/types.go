package config

import "fmt"

type Server struct {
	Id         int    `json:"id"`
	Code       string `json:"code"`
	Type       string `json:"type"`
	Ip         string `json:"main_ip"`
	MainUser   string `json:"server_main_user"`
	NginxName  string `json:"nginx_name"`
	MySQL      *Db
	PostgreSQL *Db
	Installs   []string
	Postfix    *Postfix
}

func (this Server) String() string {
	return fmt.Sprintf("%s : %s", this.Code, this.Ip)
}

func (this *Server) GetPostfixConfig() *Postfix {
	var p Postfix
	return &p
}

type ServerConf struct {
	FileName string `json:"filename"`
	FilePath string `json:"filepath"`
	Data     string `json:"data"`
}

func (this ServerConf) String() string {
	return fmt.Sprintf("%s length:%d %s", this.FileName, len(this.Data), this.FilePath)
}

type Installs []*Install

type Install struct {
	Item     string `json:"item"`
	ServerId int    `json:"server"`
}

type Config struct {
	LocalLinuxUser string
	Localhost      *Server
}

type Databases []*Db

type Db struct {
	Id         int    `json:"id"`
	ServerId   int    `json:"server"`
	Type       string `json:"type"`
	TypeDb     string `json:"type_db"`
	TypeDbName string `json:"type_db_name"`
	Version    string `json:"version"`
	User       *UserItem
}

func (this Db) String() string {
	return fmt.Sprintf("Id:%d ServerId:%d | %s %s", this.Id, this.ServerId, this.TypeDbName, this.Version)
}

type Users []*User

// no password
type User struct {
	Id       int    `json:"id"`
	ServerId int    `json:"server"`
	DbId     int    `json:"db"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func (this User) String() string {
	return fmt.Sprintf("Id:%d ServerId:%d DbId:%d Name:%s", this.Id, this.ServerId, this.DbId, this.Name)
}

// with password
type UserItem struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Postfix struct {
	Hostname      string `json:"hostname"`
	MySQLDb       string `json:"mysql_db"`
	MySQLUser     string `json:"mysql_user"`
	MySQLPassword string `json:"mysql_password"`
	MySQLSalt     string `json:"mysql_salt"`
	Emails        []*PostfixEmail
}

func (this Postfix) String() string {
	return fmt.Sprintf("%s %s %s", this.Hostname, this.MySQLDb, this.MySQLUser)
}

type PostfixEmail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Alias    string `json:"alias"`
}

type Project struct {
	Id                  int      `json:"id"`
	Executables         []string `json:"executables"`
	Exclude             []string `json:"exclude"`
	Type                string   `json:"type"`
	Name                string   `json:"name"`
	Domain              string   `json:"domain"`
	ProjectDirServer    string   `json:"project_dir_server"`
	ProjectDirLocal     string   `json:"project_dir_local"`
	IsGit               bool     `json:"is_git"`
	PythonVersion       string   `json:"python_version"`
	DjangoVersion       string   `json:"django_version"`
	VenvDirServer       string   `json:"venv_dir_server"`
	VenvDirLocal        string   `json:"venv_dir_local"`
	StaticDirServer     string   `json:"static_dir_server"`
	StaticDirLocal      string   `json:"static_dir_local"`
	MediaDirServer      string   `json:"media_dir_server"`
	MediaDirLocal       string   `json:"media_dir_local"`
	IsStaticDirSeparate bool     `json:"is_static_dir_separate"`
	RequirementsDir     string   `json:"requirements_dir"`
	UwsgiPort           int      `json:"uwsgi_port"`
	PythonPathServer    string   `json:"python_path_server"`
	PythonPathLocal     string   `json:"python_path_local"`
	ReloadIniPath       string   `json:"reload_ini_path"`
	Server              *Server  `json:"server"`
	NginxConfs          ProjectConfs
	SupervisorConfs     ProjectConfs
}

func (this Project) String() string {
	return fmt.Sprintf("%s (%s) %s", this.Name, this.Type, this.Server.Code)
}

type ProjectConfs []*ProjectConf

type ProjectConf struct {
	FileName string `json:"filename"`
	FilePath string `json:"filepath"`
	Data     string `json:"data"`
}

func (this ProjectConf) String() string {
	return fmt.Sprintf("%s length:%d %s", this.FileName, len(this.Data), this.FilePath)
}
