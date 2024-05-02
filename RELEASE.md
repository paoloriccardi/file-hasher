## RELEASE NOTE

*02/05/2024*

Implemented sqlite sink for file data, still need to test it though because of a problem with sqlite library on windows (need gcc to compile the packagemattn/go-sqlite3). For more details, see:
- [Install without gcc (not possible)][def1]
- [Getting errors when running go project that uses github.com/mattn/go-sqlite3 library][def2]

## TODO

- Have a cleaner way than an `if` statement to tell the tool to either send the output to a csv or to a sqlite database
- Export duplicates found in the file collection


[def1]: https://github.com/mattn/go-sqlite3/issues/212
[def2]: https://forum.golangbridge.org/t/getting-errors-when-running-go-project-that-uses-github-com-mattn-go-sqlite3-library/31800