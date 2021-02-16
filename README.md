# golang-api

A simple JSON API that returns data from an external data source (NYC Citibikes)

## Pre-requisites
* Download Go 1.11 or later: https://golang.org/
    * Needed for Go modules: https://github.com/golang/go/wiki/Modules#how-to-use-modules 
* Ensure sure you can run `make`
    * Linux: https://www.gnu.org/software/make/
    * MacOS: https://formulae.brew.sh/formula/make
    * Windows: http://gnuwin32.sourceforge.net/packages/make.htm
* Make sure `localhost:4000` is not being used
    * Also check your firewall to make sure you can run locally and expose port 4000
* Clone this repository and `cd` into the root directory

## Building, testing, and running the app
Runs clean, fmt, build, and test  
`$ make all`

Runs program at port 4000  
`$ make run`
