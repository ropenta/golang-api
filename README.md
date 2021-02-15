# golang-api

A simple JSON API that returns data from an external data source (NYC Citibikes)

## Pre-requisites
* Download Go 1.11 or later: https://golang.org/
    * Needed for Go modules: https://github.com/golang/go/wiki/Modules#how-to-use-modules 
* Ensure sure you can run `make`
    * Linux: https://www.gnu.org/software/make/
    * MacOS: https://formulae.brew.sh/formula/make
    * Windows: http://gnuwin32.sourceforge.net/packages/make.htm
* Clone this repository and `cd` into the root directory
* Check your local firewall to make sure you can run locally and expose port 4000

## Building, unit testing, and running the app
`$ make build`
`$ make test`
`$ make run`
