package clientCalculix

// Database for saving input and output information
// main folder: dbDatFolder
// indide main folder:
// * db.txt - database file in format of line:
// e2c569be17396eca2a2e3c11578123ed 5
// md5 + space(" ") + name of folder
// secondary folder name : 1, 2, 3, ...
// inside secondary folder:
// * model.md5 - MD5 hash algorithm of input file
// * model.inp - input inp file
// * model.dat - output dat file

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	dbDatabaseFolder = "dbDat"
	md5File          = "model.md5"
	inpFile          = "model.inp"
	datFile          = "model.dat"
)

type row struct {
	md5line    [16]byte
	folderName string
}

type database struct {
	rows []row
}

func newDb() (db *database, err error) {
	db = new(database)

	// read file
	err = db.readDatabase()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *database) readDatabase() (err error) {
	dir := "." + string(filepath.Separator) + dbDatabaseFolder
	// create db folder if not exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dbDatabaseFolder, 0755)
		if err != nil {
			return err
		}
		return nil
	}

	data, _ := ioutil.ReadDir(dir)
	for _, f := range data {
		if !f.IsDir() {
			continue
		}
		modelDir := dir + string(filepath.Separator) + f.Name()
		filename := modelDir + string(filepath.Separator) + md5File
		md5line, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		var md5byte [16]byte
		if len(md5line) != 16 {
			continue
		}
		for i := 0; i < 16; i++ {
			md5byte[i] = md5line[i]
		}
		db.rows = append(db.rows, row{md5line: md5byte, folderName: modelDir})
	}

	return nil
}

func (db *database) get(inp string) (datBody string, err error) {
	md5line := md5.Sum([]byte(inp))
	for _, row := range db.rows {
		if row.md5line == md5line {
			// compare inp file with inp file
			filename := row.folderName + string(filepath.Separator) + inpFile
			bytes, err := ioutil.ReadFile(filename)
			if err != nil {
				return "", fmt.Errorf("Cannot read inp file :%v", err)
			}
			if inp != string(bytes) {
				continue
			}
			// return result
			filename = row.folderName + string(filepath.Separator) + datFile
			bytes, err = ioutil.ReadFile(filename)
			if err != nil {
				return "", fmt.Errorf("Cannot read dat file: %v", err)
			}
			return string(bytes), nil
		}
	}
	return datBody, fmt.Errorf("Cannot found in db")
}

func (db *database) write(inp string, dat string) (err error) {
	md5line := md5.Sum([]byte(inp))

	// check result is exist
	for _, row := range db.rows {
		if row.md5line == md5line {
			return nil
		}
	}

	// create new result folder
	dir, err := db.createNewDir()
	if err != nil {
		return err
	}

	// write results
	err = ioutil.WriteFile(dir+string(filepath.Separator)+inpFile, []byte(inp), 0755)
	if err != nil {
		return fmt.Errorf("Cannot write inp file : %v", err)
	}
	err = ioutil.WriteFile(dir+string(filepath.Separator)+datFile, []byte(dat), 0755)
	if err != nil {
		return fmt.Errorf("Cannot write dat file : %v", err)
	}
	var m []byte
	for _, b := range md5line {
		m = append(m, b)
	}
	err = ioutil.WriteFile(dir+string(filepath.Separator)+md5File, m, 0755)
	if err != nil {
		return fmt.Errorf("Cannot write dat file : %v", err)
	}

	// write database
	db.rows = append(db.rows, row{md5line: md5line, folderName: dir})

	return err
}

func (db *database) createNewDir() (dir string, err error) {
	for i := 0; i < 1000000; i++ {
		dir = string(".") + string(filepath.Separator) + dbDatabaseFolder + string(filepath.Separator) + fmt.Sprintf("Row(%v)", i)
		if _, err = os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, 0755)
			if err != nil {
				continue
			}
			return dir, nil
		}
	}
	return "", fmt.Errorf("Cannot create temp folder: %v", err)
}
