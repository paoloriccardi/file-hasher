# RELEASE NOTE

### Build
Succesfully tested on windows, in order to build the tool correctly I had to: 

Install [Msys64][mysis] and gcc on mingw64 with pacman:
```
pacman -S mingw-w64-x86_64-gcc
```

Add *C:\msys64\mingw64\bin* to *$PATH$* in windows environment variables.

Enable CGO before downloading the needed package to use sqlite (*github.com/mattn/go-sqlite3*), otherwise the following error is produced on a Windows 11 machine if CGO is not enabled.
```
Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
```

After that, the sqlite driver can be installed via:
```
go get github.com/mattn/go-sqlite3
go install github.com/mattn/go-sqlite3
```

For more details, see:
- [Install without gcc (not possible)][mattnforum]
- [Getting errors when running go project that uses github.com/mattn/go-sqlite3 library][def]

### Select datasink output type
Used strategy pattern for adding new datasink output type, actually supported both:
- csv single file
- sqlite database

### Find duplicates
The duplicates are found with a sql query and they are inserted in a table called duplicates.

## TODO

- Export duplicates found in the file collection in other format than sqlite
  - Find a more elegant way than sha1DuplicatesInSqliteToTable (use of interface in order to decouple the way to export the duplicates csv|table, plus pass source and sink as parameter)


[mattnforum]: https://github.com/mattn/go-sqlite3/issues/212
[def]: https://forum.golangbridge.org/t/getting-errors-when-running-go-project-that-uses-github-com-mattn-go-sqlite3-library/31800
[mysis]: https://www.msys2.org/
