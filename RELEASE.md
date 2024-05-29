# RELEASE NOTE

### Backup output to AWS S3
The database or the csv produced by the file scan can now be saved to an S3 bucket. [go-s3]

### Config file
Added a config_sample.json file and added config.json to .gitignore.

## TODO

- Export duplicates found in the file collection in other format than sqlite
  - Find a more elegant way than sha1DuplicatesInSqliteToTable (use of interface in order to decouple the way to export the duplicates csv|table, plus pass source and sink as parameter)


[mattnforum]: https://github.com/mattn/go-sqlite3/issues/212
[def]: https://forum.golangbridge.org/t/getting-errors-when-running-go-project-that-uses-github-com-mattn-go-sqlite3-library/31800
[mysis]: https://www.msys2.org/
[go-s3] : https://docs.aws.amazon.com/sdk-for-go/api/ 
