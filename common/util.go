package common

import (
	"bytes"
	"errors"
	"time"

	// "github.com/bwmarrin/discordgo"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const DiscordEpoch = 1420070400000

//160724514639

//Turns a date in form yyyy-mm-dd into discord unix time (epoch at 2015-01-01)
func DateToDiscordTime(date string) (int64, error) {
	split := strings.Split(date, "-")
	//["2021", "02", "23"]
	year, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return 0, err
	}
	month, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return 0, err
	}
	day, err := strconv.ParseInt(split[2], 10, 64)
	if err != nil {
		return 0, err
	}
	zone, _ := time.Now().Zone()
	loc, _ := time.LoadLocation(zone)
	//make unix time and convert to discord time
	t := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, loc)

	unix := t.Unix()*1000 - DiscordEpoch
	return unix, nil

}

//Turns a date in form yyyy-mm-dd HH:MM:SS TIMEZONE into time
func DateToTime(date string) (time.Time, error) {
	//[yyyy-mm-dd, HH:MM:SS, TIMEZONE]
	split := strings.Split(date, " ")
	//["2021", "02", "23"]
	dates := strings.Split(split[0], "-")
	//["12", "11", "42"]
	times := strings.Split(split[1], ":")
	//EST
	timeZone := split[2]

	year, err := strconv.ParseInt(dates[0], 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	month, err := strconv.ParseInt(dates[1], 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	day, err := strconv.ParseInt(dates[2], 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	hour, err := strconv.ParseInt(times[0], 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	minute, err := strconv.ParseInt(times[1], 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	second, err := strconv.ParseInt(times[2], 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	loc, err2 := time.LoadLocation(timeZone)

	t := time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, loc)

	return t, err2
}

//if embed or attachment
//GET file url's response.body.
//get bytes of response via ReadAll()
//get hash of file via bytes  sha256.Sum256([]byte("hello world\n"))
//write bytes into file with hash as name into specified path
//return sha256 sum
func DownloadFile(url string, channel_id string, path string, save bool) (string, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	newpath := filepath.Join(".", path, channel_id)
	os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	//read response stream into byte array
	body, err := io.ReadAll(resp.Body)
	//Duplicate body stream as you cannot read it twice
	body2 := bytes.NewReader(body)
	if err != nil {
		return "", err
	}
	//hash byte array
	sum := fmt.Sprintf("%x", sha256.Sum256(body))

	if !save {
		return sum, nil
	}

	//create file with hash as file name
	newpath = filepath.Join(".", newpath, sum)
	_, errExist := os.Stat(newpath)
	if errExist == nil {
		//If exist, return hash and do not download file
		log.Println("File exists")
		return sum, nil
	}
	if errors.Is(errExist, os.ErrNotExist) {
		//If file does not exist, download file and return sum
		out, err := os.Create(newpath)
		if err != nil {
			return "", err
		}
		defer out.Close()
		// Write the body to file
		_, err = io.Copy(out, body2)
		if err != nil {
			return sum, err
		}
		return sum, nil
	}
	//If error is not errNotExist, return error and sum. Something is wrong
	return sum, errExist
}
