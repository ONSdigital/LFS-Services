
## LFS Imports REST API

Provides the APIs for the LFS services.

### Dependencies
The API used the readstat library to load and save SPSS sav files. Install it from:

https://github.com/WizardMac/ReadStat

Note that there is currently an error in the windows install instructions.

> pacman -S autoconf automake libtool mingw-w64-x86_64-toolchain ingw-w64-x86_64-cmake mingw-w64-x86_64-libiconv

should be

> pacman -S autoconf automake libtool mingw-w64-x86_64-toolchain mingw-w64-x86_64-cmake mingw-w64-x86_64-libiconv make

For linux the LD_LIBRARY_PATH environment variable must include 
_/usr/local/lib_ otherwise the app will be unable to find the Readstat C libraries.

### Configuration

A configuration framework is included to set various configuration properties. See the _config/config.development.toml_ file
for current configuration options.

### Running

You will need a suitable PostgreSql installation. The schema is under the _scripts_ directory and the configuration is set in 
the configuration file. Note that the database configuration can, and really must, be 
overidden by the following environment variables:

	DB_SERVER
	DB_USER
	DB_PASSWORD
	DB_DATABASE

### Dockerfile

Two dockerfiles are provided. The first `dockerfile.debug` is for running a delve server in docker and the second, 
`Dockerfile.dev` is for running the service under docker in development mode. 
