
## To install on Mac
Install

```brew install postgresql```

Start database

```brew services start postgresql```

Enter commands to create user and database
```
/usr/local/opt/postgres/bin/createuser -s postgres
/usr/local/opt/postgres/bin/createuser lfs -d -P -- enter password when prompted
/usr/local/opt/postgres/bin/createdb lfs -O lfs
```

Run postgress schema script (currently that entails copying and pasting the contents of said script into a runnable postgresql script)

## To install on PC
```
scoop install postgresql
createuser -s postgres
createuser lfs -d -P -- enter password when prompted
createdb lfs -O lfs
```

You will then be able to login as user lfs, password lfs and database lfs.

Change as necessary and remember to change the corresponding items in the configuration file.
