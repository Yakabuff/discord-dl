package main
import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"crypto/sha256"
	"log"
	"errors"
	"fmt"
)

//channel_id/sha256hash
//download file in channel_id folder
//calculate hash
//rename file to hash
func DownloadFile(url string, channel_id string) (error, string) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err, ""
	}
	defer resp.Body.Close()

	newpath := filepath.Join(".", "media", channel_id)
	os.MkdirAll(newpath, os.ModePerm)
	newpath = filepath.Join(".", "media", channel_id, "tmp")
	// Create the file
	out, err := os.Create(newpath)
	if err != nil {
		return err, ""
	}
	
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	out.Close()

	f, err := os.Open(newpath)
	if err != nil {
		return err, ""
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
	  log.Println("could not hash file")
	  return err, ""
	}
	//Note: try hashing resp.Body????
	hash := fmt.Sprintf("%x",h.Sum(nil))

	hashpath := filepath.Join(".", "media", channel_id, hash)
	if _, err := os.Stat(hashpath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist. rename tmp to hash
		log.Println("New hash found: "+ hash)
		err := os.Rename(newpath, hashpath)
		if err != nil{
			log.Println(err)
			return err, ""
		}
	}else{
		//delete tmp
		log.Println("Hash already exists: " + hash)
		err := os.Remove(newpath)
		if err != nil {
			log.Println(err)
			return err, "";
		}
	}

	return err, hash
}
