package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ConfigFile struct {
	CsvSeparator string `json:"csvseparator"`
	CsvDelimiter string `json:"csvdelimiter"`
	TargetDir    string `json:"targetdir"`
	OutFile      string `json:"outfile"`
}

var Cfg ConfigFile

func init() {
	confFile := os.Args[1]
	confJson, err := os.Open(confFile)
	if err != nil {
		log.Fatal(err)
	}
	defer confJson.Close()

	byteValue, _ := io.ReadAll(confJson)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &Cfg)
}

// FileData contains the relevant file information
type FileData struct {
	FilePath     string
	FileName     string
	FileSize     string
	FileChecksum string
}

func (r FileData) toCsvRow(separator string, delimiter string, encapsulated bool) string {
	if encapsulated {
		return strconv.Quote(r.FilePath) + separator + strconv.Quote(r.FileName) + separator + strconv.Quote(r.FileSize) + separator + strconv.Quote(r.FileChecksum) + delimiter
	} else {
		return r.FilePath + separator + r.FileName + separator + r.FileSize + separator + r.FileChecksum + delimiter
	}
}

// GenerateFileChecksum calculate the sha1 checksum if the resource specified by filepath is a file, it returns a FileData, plus a boolean to signal to the caller that filepath was a Directory, error if it fails
func GenerateFileChecksum(filepath string, d os.DirEntry) (FileData, bool, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return FileData{"", "", "", ""}, false, err
	}
	defer f.Close()

	fInfo, err := f.Stat()
	if err != nil {
		return FileData{"", "", "", ""}, false, err
	}

	if !fInfo.IsDir() {
		h := sha1.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal()
		}
		fdata := FileData{strings.ReplaceAll(filepath, fInfo.Name(), ""), fInfo.Name(), strconv.FormatInt(fInfo.Size(), 10), hex.EncodeToString(h.Sum(nil))}
		return fdata, false, nil
	} else {
		return FileData{"", "", "", ""}, true, nil
	}
} // we initialize our Users array

// ScanFileCsvOut is the callback function used by WalkDir, it calls GenerateFileChecksums to retrieve the file info as a FileData struct and then it append them to outfile.csv, in case the specified path is a File, it skips if it's a Directory
func ScanFileCsvOut(path string, d os.DirEntry, err error) error {
	if err != nil {
		return err
	}

	fdata, isdir, err := GenerateFileChecksum(path, d)
	if err != nil {
		return err
	}

	if !isdir {
		f, err := os.OpenFile(Cfg.OutFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		f.WriteString(fdata.toCsvRow(Cfg.CsvSeparator, Cfg.CsvDelimiter, true))
	}
	return nil
}

func main() {

	// we initialize our Users array
	//err := Cfg.init(os.Args[1])

	/*
		if err != nil {
			log.Fatal(err)
		}
	*/

	log.Println(Cfg)

	err := filepath.WalkDir(Cfg.TargetDir, ScanFileCsvOut)

	if err != nil {
		log.Println(err)
	}
}
