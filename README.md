Tiny Resource Statistics
=========================
[![Go Report Card](https://goreportcard.com/badge/github.com/alekns/tinyrstats)](https://goreportcard.com/report/github.com/alekns/tinyrstats)

This is tiny demo project (single service solution), to able monitor resources and gathering of some statistics.


Features
=============

- [ ] [HTTP REST API](https://github.com/alekns/tinyrstats/blob/master/http-api.apib)
- [ ] Monitoring of HTTP endpoints
- [ ] Resource statistics (response time, timeouts, errors, etc.)
- [ ] Endpoints statistics
- [ ] Requests tracing


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


Run
========

For host: `$ ./bin/tinyrstats --config-file ./config.common.yml serve`
You may override config by environment variables. Set logging level for example:
For host: `$ TRS_LOGGING_CONSOLE_LEVEL=info ./bin/tinyrstats --config-file ./config.defaults.yml serve --preload-from-file sites.example.txt`
