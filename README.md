Tiny Resource Statistics
=========================
[![Go Report Card](https://goreportcard.com/badge/github.com/alekns/tinyrstats)](https://goreportcard.com/report/github.com/alekns/tinyrstats)

This is tiny demo project (single service solution), to able monitor resources and gathering of some statistics.


Features
=============

- [x] [HTTP REST API](https://github.com/alekns/tinyrstats/blob/master/http-api.apib)
- [x] Monitoring of HTTP endpoints
- [x] Resource statistics (response time, timeouts, errors, etc.)
- [x] Endpoints statistics and metrics (**http://localhost:8081/metrics** or **http://localhost:18081/metrics** for docker)
- [x] Requests tracing


Build
======

For host: `make skip_dep=false clean build`


Build and run in Docker
========================

`$ docker-compose up -d` (default port binding :18080 -> :8080, :18081 -> :8081)


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

For host: `$ ./bin/tinyrstats --config-file ./config.defaults.yml serve --preload-from-file sites.example.txt`

You may override config by environment variables. Set logging level for example,
for host: `$ TRS_LOGGING_CONSOLE_LEVEL=debug ./bin/tinyrstats --config-file ./config.defaults.yml serve --preload-from-file sites.example.txt`

Use `--default-protocol` to select between `http` and `https` (override by individual).
You could customize individual protocol in sites.example.txt by prefix `http://` or `https://`.

To enable tracing (in docker-compose enabled by default) for
host (you should have runned jaeger agent): `$ JAEGER_SAMPLER_TYPE=const JAEGER_SAMPLER_PARAM=1 JAEGER_AGENT_HOST=localhost JAEGER_AGENT_PORT=16831 ./bin/tinyrstats monitor --config-file ./config.defaults.yml serve --preload-from-file ./sites.example.txt`

Docker jaeger is available on **http://localhost:16686/**
