
## To install on Mac

brew install postgresql

* /usr/local/opt/postgres/bin/createuser -s postgres 

* /usr/local/opt/postgres/bin/createuser lfs -d -P -- enter password when prompted

* /usr/local/opt/postgres/bin/createdb lfs -O lfs

You will then be able to login as user lfs, password lfs and database lfs.

Change as necessary and remember to change the corresponsing items in the configuration file.

## To install on PC

scoop install postgresql
