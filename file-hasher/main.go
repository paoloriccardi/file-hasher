package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// ConfigFile is the struct that contains the configuration parameters specified in the config.json file
type ConfigFile struct {
	CsvSeparator string `json:"csvseparator"`
	CsvDelimiter string `json:"csvdelimiter"`
	TargetDir    string `json:"targetdir"`
	OutFile      string `json:"outfile"`
	DBPath       string `json:"dbpath"`
	DBFile       string `json:"dbfile"`
}

var Cfg ConfigFile

// initialize the ConfigFile, the path to the config file is specified as the first command line argument
func init() {
	confFile := os.Args[1]
	confJson, err := os.Open(confFile)
	if err != nil {
		log.Fatal(err)
	}
	defer confJson.Close()

	byteValue, _ := io.ReadAll(confJson)

	// we unmarshal our byteArray which contains our
	// confFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &Cfg)

	// If the outfile field is an empty string in the config file it means we want to use sqlite as the output sink for our data.
	if Cfg.OutFile == "" {
		dbfile, err := os.OpenFile(filepath.Join(Cfg.DBPath, Cfg.DBFile), os.O_RDWR|os.O_CREATE, 0755)
		if errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
		defer dbfile.Close()

		db, err := sql.Open("sqlite3", filepath.Join(Cfg.DBPath, Cfg.DBFile))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		defer db.Close()

		createFileTableIfnotExist := "CREATE TABLE IF NOT EXISTS `files` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `directory` TEXT NULL, `filename` TEXT NULL, `size` text NULL, `sha1` TEXT NULL, `date` DATETIME NULL)"

		_, err = db.Exec(createFileTableIfnotExist)
		if err != nil {
			log.Fatal(err)
		}
	}

}

// FileData contains the relevant file information
type FileData struct {
	FilePath     string
	FileName     string
	FileSize     string
	FileChecksum string
	ScanTime     time.Time
}

func (r FileData) appendToCsvFile() error {
	f, err := os.OpenFile(Cfg.OutFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(r.toCsvRow(Cfg.CsvSeparator, Cfg.CsvDelimiter, true))
	return nil
}

func (r FileData) toSqliteDB(db *sql.DB) error {
	insertFileData := `INSERT INTO files (directory,filename,size,sha1,date) VALUES (?,?,?,?,?)`
	_, err := db.Exec(insertFileData, r.FilePath, r.FileName, r.FileSize, r.FileChecksum, r.ScanTime)
	if err != nil {
		return err
	}
	return nil
}

// toCsvRow returns a formatted string which represent FileData content on a single csv row
func (r FileData) toCsvRow(separator string, delimiter string, encapsulated bool) string {
	if encapsulated {
		return strconv.Quote(r.FilePath) + separator + strconv.Quote(r.FileName) + separator + strconv.Quote(r.FileSize) + separator + strconv.Quote(r.FileChecksum) + separator + r.ScanTime.String() + delimiter
	} else {
		return r.FilePath + separator + r.FileName + separator + r.FileSize + separator + r.FileChecksum + separator + r.ScanTime.String() + delimiter
	}
}

// GenerateFileChecksum calculate the sha1 checksum if the resource specified by filepath is a file, it returns a FileData, plus a boolean to signal to the caller that filepath was a Directory, error if it fails
func GenerateFileChecksum(filepath string, d os.DirEntry) (FileData, bool, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return FileData{"", "", "", "", time.Time{}}, false, err
	}
	defer f.Close()

	fInfo, err := f.Stat()
	if err != nil {
		return FileData{"", "", "", "", time.Time{}}, false, err
	}

	if !fInfo.IsDir() {
		h := sha1.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Println(err)
			return FileData{"", "", "", "", time.Time{}}, true, nil
		}
		fdata := FileData{strings.ReplaceAll(filepath, fInfo.Name(), ""), fInfo.Name(), strconv.FormatInt(fInfo.Size(), 10), hex.EncodeToString(h.Sum(nil)), time.Now()}
		return fdata, false, nil
	} else {
		return FileData{"", "", "", "", time.Time{}}, true, nil
	}
}

type DataSink struct {
	dsOutputType DSOutputType
}

func (ds *DataSink) setDsOutputType(ot DSOutputType) {
	ds.dsOutputType = ot
}

type DSOutputType interface {
	scanFile(path string, d os.DirEntry, err error) error
}

type SingleCsvFileOutType struct {
}

func (c *SingleCsvFileOutType) scanFile(path string, d os.DirEntry, err error) error {
	if err != nil {
		return err
	}

	fdata, isdir, err := GenerateFileChecksum(path, d)
	if err != nil {
		return err
	}

	if !isdir {
		if Cfg.OutFile != "" {
			err := fdata.appendToCsvFile()
			if err != nil {
				return err
			}
		} else {
			log.Println("SCnfig file could be malformed, outfile null during csv export")
		}

	}
	return nil
}

type SqliteDatabaseOutType struct {
}

func (s *SqliteDatabaseOutType) scanFile(path string, d os.DirEntry, err error) error {
	if err != nil {
		return err
	}

	fdata, isdir, err := GenerateFileChecksum(path, d)
	if err != nil {
		return err
	}

	if !isdir {
		if Cfg.OutFile != "" {
			log.Println("Config file could be malformed, outfile not null during sqlite export")
		} else {
			db, err := sql.Open("sqlite3", filepath.Join(Cfg.DBPath, Cfg.DBFile))
			if err != nil {
				log.Println(err)
				return err
			}
			defer db.Close()
			fdata.toSqliteDB(db)
		}

	}
	return nil
}

func main() {

	ds := DataSink{}
	ot := SqliteDatabaseOutType{}
	ds.setDsOutputType(&ot)

	err := filepath.WalkDir(Cfg.TargetDir, ds.dsOutputType.scanFile)
	if err != nil {
		log.Println(err)
	}

}
