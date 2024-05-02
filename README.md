# File Hasher

File hasher is a simple tool that walk through the files contained in a given folder and in all its subfolders and calculate the sha1 hash of any file found along the way.

At this point there are two possibilities, the collected information for every file (*directory, filename, size, checksum, date of retrieval*):
- are appended to a csv file stored locally
- are inserted in a local sqlite database 

This tool was thought to help me find duplicated mp3 in my Music collection. 

## How to build

It is required by the sqlite *mattn/sqlite3* library to have CGO enabled during the build via: `go env -w CGO_ENABLED=1`  otherwise the sqlite library won't work. 

After that, the sqlite driver can be installed via:
```
go get github.com/mattn/go-sqlite3
go install github.com/mattn/go-sqlite3
```

On windows there are issues with using the sqlite *mattn/sqlite3* library without having gcc installed, so I wasn't able to test it on my windows machine. 

On Linux there should be no problem provided that you have gcc installed.

## How to use the tool
The usage is pretty simple:

`./file-hasher <config file location>`

## Configuration
The only parameter taken by the tool is the location of the configuration file. The configuration file is a json file that contains the following:

- `csvseparator` the field separator to be used in the csv
- `csvdelimiter` the character to be used as a line delimiter in the csv
- `targetdir` the directory where the scan should start
- `outfile` the csv output file, if this is set to something different than `""` the output will be produced in csv, otherwise the next two parameters will be used to insert data in a sqlite database 
- `dbpath` the sqlite3 database path (usually the `data` folder)
- `dbfile` the sqlite3 database file

