# Description

This is a very simple "store" for Android applications.  It is useful for serving customers who do not have access to traditional Android stores such as Google Play or Amazon Appstore, or for developers who do not wish to use those stores to distribute their products.

The previous version was a Django app, and the heritage shows!

# Installation

This package can be installed using go get:

````
go get github.com/mathuin/androidapps
````

Change to the appropriate directory and build the app:

````
go build
````

# Documentation

## Settings

### Environment variables

The best/easiest way to configure this app is with environment variables.

````
export ANDROIDAPPS_DBFILE="androidapps.db"
export ANDROIDAPPS_HOST="0.0.0.0"
export ANDROIDAPPS_PORT="4000"
export ANDROIDAPPS_NAME="Jane Doe"
export ANDROIDAPPS_EMAIL="jane@example.net"
````

### Flags

Sometimes environment variables just won't do.  In those cases, use flags.

| Flag | Meaning |
| ~~~~ | ~~~~~~~ |
| -dbfile | Database file |
| -host | Host |
| -port | Port |
| -name | Developer name |
| -email | Developer email |

## Subcommands

| Subcommand | Purpose | Arguments |
| ~~~~~~~~~~ | ~~~~~~~ | ~~~~~~~~~ |
| runserver | Run the server | |
| export | Export the database to standard output | |
| import | Import database in export format | |
| check | Check for files corresponding to products | verbose?  string match? |
| rebuild | for each product, re-extract title, version, icon. | all? string match? |
| list | list products in database | string match? enabled? |
| enable | enable product (will need flag added to database) | all? one? |
| disable | disable product (will need flag added to database) | all? one? |
| add | add product to database | APK file |
| remove | remove product from database | force? |
| upgrade | upgrade product in database (upload new APK) | APK file |

At this time, very few commands are implemented.  Please be patient. :-)

# TODO

* Install the web stuff (static, media, templates) somewhere
* Build tools to manipulate the database
* Add QR code support

# License

This software is released under the [MIT license](http://opensource.org/licenses/mit-license.php).
