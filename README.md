negroni-gomol
=============

[![GoDoc](https://godoc.org/github.com/aphistic/negroni-gomol?status.svg)](https://godoc.org/github.com/aphistic/negroni-gomol)
[![Build Status](https://img.shields.io/travis/aphistic/negroni-gomol.svg)](https://travis-ci.org/aphistic/negroni-gomol)
[![Code Coverage](https://img.shields.io/codecov/c/github/aphistic/negroni-gomol.svg)](http://codecov.io/github/aphistic/negroni-gomol?branch=master)

Negroni middleware to log with the Gomol log library.

Usage
-----
Once gomol is initialized as you'd normally do (see the
[gomol](https://www.github.com/aphistic/gomol) documentation) you can
create the logging middleware for Negroni in one of two ways.

The first way is to just use the current default gomol logger, like so:
```
import ng "github.com/aphistic/negroni-gomol"

...

n := negroni.New()
n.Use(ng.NewLogger())
```

If you have a different gomol Base logger you'd like to use instead, you can
also use it in this way:
```
import ng "github.com/aphistic/negroni-gomol"

...

base := gomol.NewBase()

...

n := negroni.New()
n.Use(ng.NewLoggerForBase(base))
```