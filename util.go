package main

import(
	"time"
	// "github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

const DiscordEpoch = 1420070400000

                     //160724514639

//Turns a date in form yyyy-mm-dd into discord unix time (epoch at 2015-01-01)
func DateToDiscordTime(date string) (int64, error){
	split := strings.Split(date, "-");
	//["2021", "02", "23"]
	year, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil{
		return 0, err;
	}
	month, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil{
		return 0, err;
	}
	day, err := strconv.ParseInt(split[2], 10, 64)
	if err != nil{
		return 0, err;
	}
	zone, _ := time.Now().Zone()
	loc, _ := time.LoadLocation(zone)
	//make unix time and convert to discord time
	t := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, loc)
    
	unix:= t.Unix()*1000 - DiscordEpoch;
	return unix, nil;
	
}


func DateToTime(date string)(time.Time, error){
	split := strings.Split(date, "-");
	//["2021", "02", "23"]
	year, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil{
		return time.Time{}, err;
	}
	month, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil{
		return time.Time{}, err;
	}
	day, err := strconv.ParseInt(split[2], 10, 64)
	if err != nil{
		return time.Time{}, err;
	}
	zone, _ := time.Now().Zone()
	loc, _ := time.LoadLocation(zone)
	//make unix time and convert to discord time
	t := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, loc)
    
	return t, nil;
}
