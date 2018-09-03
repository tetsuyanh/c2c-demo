# c2c-demo

[![Build Status](https://travis-ci.org/tetsuyanh/c2c-demo.svg?branch=master)](https://travis-ci.org/tetsuyanh/c2c-demo)

c2c demo application by golang

- user scope (anonymous, authencated by email)
- virtual point
- item management
- deal with transaction

## setup
for macOS, and you need to install homebrew and golang

postgresql v9.6
```
$ brew update
$ brew install postgresql@9.6
$ echo 'export PATH="/usr/local/Cellar/postgresql@9.6/9.6.6/bin:$PATH"' >> ~/.bash_profile
$ source ~/.bash_profile
$ postgres --version
postgres (PostgreSQL) 9.6.6
```

migration
```
$ curl -L https://github.com/mattes/migrate/releases/download/v3.0.1/migrate.darwin-amd64.tar.gz | tar xvz
$ mv migrate.darwin-amd64 /usr/local/bin/migrate

$ migrate --version
3.0.1
```

### development

run api
```
$ psql postgres
=# create role c2c_demo with login superuser;
=# create database c2c_demo owner c2c_demo;

$ migrate -database postgres://localhost:5432/c2c_demo?sslmode=disable -path ./migration up
$ source env.sh

$ make setup
$ make install
$ make run
```

test
```
$ psql postgres
=# create role c2c_test with login superuser;
=# create database c2c_test owner c2c_test;

$ migrate -database postgres://localhost:5432/c2c_test?sslmode=disable -path ./migration up

$ make test
```
