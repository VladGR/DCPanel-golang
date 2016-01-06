## About

GoLang program for [DCPanel](https://github.com/vladgr/DCPanel).

1. It installs following software on clean Debian 8 server:

* IPTables
* MySQL
* Nginx
* PHP (not much tested - just works)
* Postfix (not much tested)
* PostgreSQL
* Python2
* Python3
* Redis
* Supervisor
* Squid

2. It deploys django projects and updates nginx and supervisor configuration.

## Dependencies

```
go get golang.org/x/crypto/ssh
```
