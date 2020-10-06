# Gopherprint

Simple menubar app to display Octprint status and progress

## Installing

Add your host and API key in a file `Makefile.vars`:

```make
SERVER_URL=http://...

API_KEY=...
```

Then build and install the app with `make install`, that will create Gopherprint.app

```shell
$ make build
```

## Developing

Follow the steps for installation, execute make build instead to not copy to Applications folder

## Libraries

- [caseymrm/menuet](https://github.com/caseymrm/menuet)
- [mcuadros/go-octoprint](https://github.com/mcuadros/go-octoprint)
- [haklop/gnotifier](https://github.com/haklop/gnotifier)
