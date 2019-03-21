Tiny Resource Statistics
=========================
[![Go Report Card](https://goreportcard.com/badge/github.com/alekns/tinyrstats)](https://goreportcard.com/report/github.com/alekns/tinyrstats)

This is tiny demo project (single service solution), to able monitor resources and gathering of some statistics.


Features
=============

- [ ] [HTTP REST API](https://github.com/alekns/tinyrstats/blob/master/http-api.apib)


Build
======

For host: `make skip_dep=false clean build`


Build and run in Docker
========================

`$ docker-compose up -d` (default port binding :18080 -> :8080)


Tests
=======

`$ make test`


Format
========

`$ make format`


Linter
========

`$ make linter`
