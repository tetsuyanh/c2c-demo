#!/bin/bash
set -ex
psql -c 'create user c2c_test superuser;' -U postgres
psql -c 'create database c2c_test owner c2c_test;' -U postgres
wget https://github.com/mattes/migrate/releases/download/v3.0.1/migrate.linux-amd64.tar.gz -P /tmp
tar -xzvf /tmp/migrate.linux-amd64.tar.gz
./migrate.linux-amd64 -database postgres://localhost:5432/c2c_test?sslmode=disable -path ./migration up
