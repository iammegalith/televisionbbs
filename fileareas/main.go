package fileareas

import (
	"database/sql"
	"fmt"
	"televisionbbs/util"
	"time"
)

func fileUpload(db *sql.DB, areaName, fileName, description, uploadedBy string, fileSize int64) error {
	stmt, err := db.Prepare("INSERT INTO fileareas(areaname, filename, description, uploadedby, date, size) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(areaName, fileName, description, uploadedBy, time.Now(), fileSize)
	if err != nil {
		return fmt.Errorf("error inserting file: %v", err)
	}
	fmt.Println("File uploaded successfully.")
	return nil
}

func fileDownload(db *sql.DB, fileID int) ([]byte, error) {
	var fileData []byte

	stmt, err := db.Prepare("SELECT data FROM fileareas WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(fileID).Scan(&fileData)
	if err != nil {
		return nil, fmt.Errorf("error fetching file: %v", err)
	}

	return fileData, nil
}
func listFiles(db *sql.DB, area string, sortCriteria string) ([]util.File, error) {
	var orderBy string
	switch sortCriteria {
	case "date":
		orderBy = "date"
	case "filename":
		orderBy = "filename"
	case "uploader":
		orderBy = "uploadedby"
	default:
		return nil, fmt.Errorf("invalid sort criteria: %s", sortCriteria)
	}

	rows, err := db.Query("SELECT id, filename, description, uploadedby, date, size FROM fileareas WHERE areaname = ? ORDER BY "+orderBy, area)
	if err != nil {
		return nil, fmt.Errorf("error querying file area: %v", err)
	}
	defer rows.Close()

	var files []util.File
	for rows.Next() {
		var f util.File
		if err := rows.Scan(&f.ID, &f.FileName, &f.Description, &f.UploadedBy, &f.Date, &f.Size); err != nil {
			return nil, fmt.Errorf("error scanning file: %v", err)
		}
		files = append(files, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching files: %v", err)
	}
	return files, nil
}

func deleteFile(db *sql.DB, id int) error {
	stmt, err := db.Prepare("DELETE FROM fileareas WHERE id = ?")
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}
	fmt.Println("File deleted successfully.")
	return nil
}

func createFileArea(db *sql.DB, areaname, description string) error {
	stmt, err := db.Prepare("INSERT INTO fileareas(areaname, description) VALUES(?,?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(areaname, description)
	if err != nil {
		return fmt.Errorf("error inserting file area: %v", err)
	}
	fmt.Println("File area created successfully.")
	return nil
}

func deleteFileArea(db *sql.DB, area string) error {
	// prepare statement
	stmt, err := db.Prepare("DELETE FROM fileareas WHERE areaname = ?")
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(area)
	if err != nil {
		return fmt.Errorf("error deleting file area: %v", err)
	}

	return nil
}
