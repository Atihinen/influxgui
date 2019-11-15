# InfluxGUI

[![Build Status](https://travis-ci.com/Atihinen/influxgui.svg?branch=master)](https://travis-ci.com/Atihinen/influxgui)

Standalone GUI for influxDB for debugging and verify Influx Query


# Installation

Download correct binary for your OS from releases.

## macOS

You need to give permissions to the binary with `chmod u+x` before running the application.

# Usage

After download, douple click the binary to open the GUI.

Connect to database by inputting the hostname, port and credentials (if there are none set, leave empty). Then press "Connect".

Then run queries against correct database.


## Example queries

`SHOW MEASUREMENTS`

`SELECT * FROM <measurement> LIMIT 1;`

# Development

## Setup dev env

1. [Install GO](https://golang.org/doc/install)
2. Clone this repository
3. Run: `go get'

## Automated test cases

Run `make test`

## Build

* To run application: `make run`
* To build dev build (linux only): `make build-dev`
* To build to all platforms (linux deb, macOS, windows exe): `make build`
  * requires docker-ce





