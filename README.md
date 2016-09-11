
elsample
========

[![GoCard][1]][2]

[1]: https://goreportcard.com/badge/LeKovr/elsample
[2]: https://goreportcard.com/report/github.com/LeKovr/elsample

[elsample](https://github.com/LeKovr/elsample) - sample [ELSA](https://github.com/LeKovr/elsa) server application.

This web application skeleton aimed to be a base for other web applications.

Do not reinvent wheels, just fork this code and add your custom application logic.

Features
--------

### Basic

* [x] [negroni](https://github.com/urfave/negroni)
* [x] [recovery](https://github.com/urfave/negroni#recovery)
* [x] [static files](https://github.com/urfave/negroni#static)
* [ ] tests

### Extended

* [x] [colog](https://github.com/comail/colog) logger
* [x] config via [go-flags](https://github.com/jessevdk/go-flags)
* [x] [graceful](https://gopkg.in/tylerb/graceful.v1) shutdown

### ELSA Middlewares

* [x] [render](https://gopkg.in/unrolled/render.v1) templates
* [x] [ace](https://github.com/yosssi/ace) templates
* [x] [sample](sample/) - sample with context (go 1.7) usage demo
* [x] [cors](github.com/rs/cors)
* [x] [stats](github.com/thoas/stats)

### Addons

* [x] Makefile with usefull commands and version attr injection

Install
-------

You need to install [consup](https://github.com/LeKovr/consup) for running test site.
After installing, symlink it to $GOPATH/src/github.com/LeKovr/consup

```
sudo echo "127.0.0.1 app.dev.lan" >> /etc/hosts
go get github.com/LeKovr/elsample
cd $GOPATH/src/github.com/LeKovr/elsample
make start
grep "random pass" log/app.dev.lan/sample-stderr.log
```

Open in browser http://app.dev.lan/my with user/password from grep results

License
-------

The MIT License (MIT), see [LICENSE](LICENSE).

Copyright (c) 2016 Alexey Kovrizhkin ak@elfire.ru
