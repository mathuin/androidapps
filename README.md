# Description

This is a very simple "store" for Android applications.  It is useful for serving customers who do not have access to traditional Android stores such as Google Play or Amazon Appstore, or for developers who do not wish to use those stores to distribute their products.

The previous version was a Django app, and the heritage shows!

# Installation

This package can be installed using go get:

````
go get github.com/mathuin/androidapps
````

# Documentation

## Populating the database

I'm currently using a database made with other tools.  I still need to write some new tools to manipulate the database as well as add or update existing applications in the store.

## Running the server

First, assign the appropriate environment variables:

````
export ANDROIDAPPS_DB="androidapps.db"
export ANDROIDAPPS_HOST="0.0.0.0"
export ANDROIDAPPS_PORT="4000"
export ANDROIDAPPS_NAME="Jane Doe"
export ANDROIDAPPS_EMAIL="jane@example.net"
````

Next, run the binary from the source directory:

````
./androidapps runserver
````

# TODO

* Install the web stuff (static, media, templates) somewhere
* Build tools to manipulate the database
* Add QR code support

# License

This software is released under the [MIT license](http://opensource.org/licenses/mit-license.php).
