# File Hasher

File hasher is a simple tool that walk through the files contained in a given folder and in all its subfolders and calculate the sha1 hash of any file found along the way. 
The collected information for every file (`file path, file name, file size, hash`) are then appended to a csv file stored locally. 

This tool was thought to help me find duplicated mp3 in my Music collection. 

## How to build
Simply build the package for your platform.

## How to use the tool
The usage is pretty simple:
`./file-hasher <config file location>`

## Configuration
The only parameter taken by the tool is the location of the configuration file. The configuration file is a json file that contains the following:
- "csvseparator": the field separator to be used in the csv
- "csvdelimiter": the character to be used as a line delimiter in the csv
- "targetdir": the directory where the scan should start
- "outfile": the csv output file

