# File Hasher

File hasher is a simple tool that walk through the files contained in a given folder and in all its subfolders and calculate the sha1 hash of any file found along the way.

At this point there are two possibilities, the collected information for every file (*directory, filename, size, checksum, date of retrieval*):
- are appended to a csv file stored locally
- are inserted in a local sqlite database 

This tool was thought to help me find duplicated mp3 in my Music collection. 

## How to build

On windows systems gcc is required in order to succesfully install and use the sqlite go driver *github.com/mattn/go-sqlite3*.

### Install  mingw64 with pacman

Install [Msys64][mysis] and then add *gcc* with *pacman* using the command: 
```
pacman -S mingw-w64-x86_64-gcc
```

### Add gcc bin folder to Path

To allow application to call directly gcc from any place of the filesystem, add *C:\msys64\mingw64\bin* to *PATH* in windows environment variables.

### Enable CGO

Before downloading the needed package to use sqlite (*github.com/mattn/go-sqlite3*), otherwise the following error is produced on a Windows 11 machine if CGO is not enabled.
```
Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
```

It is required by the sqlite *mattn/sqlite3* library to have CGO enabled during the build via: `go env -w CGO_ENABLED=1`  otherwise the sqlite library won't work.

### Install the driver
```
go get github.com/mattn/go-sqlite3
go install github.com/mattn/go-sqlite3
```

On Linux, apart from enabling CGO the same way, there should be no problem provided that you have gcc installed which should be done via `apt install build-essentials`

## How to use the tool
The usage is pretty simple:

`./file-hasher <config file location>`

## Configuration
The only parameter taken by the tool is the location of the configuration file. The configuration file is a json file that contains the following:

- `csv_separator` the field separator to be used in the csv
- `csv_delimiter` the character to be used as a line delimiter in the csv
- `target_dir` the directory where the scan should start
- `out_file` the csv output file, used if the attribute `export_to` is equal to *csv*  
- `db_path` the sqlite3 database path (usually the `data` folder)
- `db_file` the sqlite3 database file
- `export_to` can be either *sqlite|csv*


[mysis]: https://www.msys2.org/