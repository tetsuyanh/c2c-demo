# c2c-demo

c2c demo application by golang

## setup
only macOS, and you need to install homebrew and golang

postgresql v9.6
```
$ brew update
$ brew install postgresql@9.6
$ echo 'export PATH="/usr/local/Cellar/postgresql@9.6/9.6.6/bin:$PATH"' >> ~/.bash_profile
$ source ~/.bash_profile
$ postgres --version
postgres (PostgreSQL) 9.6.6

$ psql postgres
=# create role c2c_demo with login superuser;
=# create database c2c_demo owner c2c_demo;
```

migration
```
$ curl -L https://github.com/mattes/migrate/releases/download/v3.0.1/migrate.darwin-amd64.tar.gz | tar xvz
$ mv migrate.darwin-amd64 /usr/local/bin/migrate

$ migrate --version
3.0.1

$ migrate -database postgres://localhost:5432/c2c_demo?sslmode=disable -path ./migration up
```

install golang packages
```
$ make setup
$ make install
```

# development

```
$ export C2C_DEMO_CONF_PATH=$(pwd)/conf
$ export C2C_DEMO_CONF_NAME=conf
$ make run/api
```
