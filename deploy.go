package main

import(
	"html/template"
	"net/http"
	"fmt"
	"log"
	"regexp"
	"time"
)
//simple rest api
// localhost:8080/<guild_id>/<channel_id>/date


func Deploy(){
	log.Println("starting web app")

	http.HandleFunc("/", messageHandler)
	http.ListenAndServe(":8000", nil) 
}
var urlExp = regexp.MustCompile(`/(?P<guild>\d+)/(?P<channel>\d+)/(?P<date>\d*)`)

func messageHandler(w http.ResponseWriter, r *http.Request) {

	match := urlExp.FindStringSubmatch(r.URL.Path)
	if len(match) > 0 {
        result := make(map[string]string)
        for i, name := range urlExp.SubexpNames() {
            if i != 0 && name != "" {
                result[name] = match[i]
            }
        }

		var date_unix int
		date, err2 := DateToTime(result["date"])
		if err2 != nil{
			date_unix = int(time.Now().Unix())
		}else{
			date_unix = int(date.Unix())
		}
		log.Println(date_unix)
		tmpl,err := template.ParseFiles("web/channel.html")
		if err != nil{
			log.Println(err)
		}
		_, msgs := getMessages(db, result["guild"], result["channel"], date_unix)
		for _, i := range(msgs.Messages){
			log.Println(i.Message_id)
		}
		tmpl.Execute(w, *msgs)
		
    } else {
        fmt.Fprintf(w, "Wrong url\n")
    }
}