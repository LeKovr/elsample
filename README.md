
elsample
========

[![GoCard][1]][2]

[1]: https://goreportcard.com/badge/LeKovr/elsample
[2]: https://goreportcard.com/report/github.com/LeKovr/elsample

[elsample](https://github.com/LeKovr/elsample) - sample [ELSA](https://github.com/LeKovr/elsa) server application.

Features
--------

* app flags inheritance
* sample package with functional options
* Makefile with usefull commands and version attr injection

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
