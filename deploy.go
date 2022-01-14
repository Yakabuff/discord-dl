package main

import(
	"html/template"
	"net/http"
	// "fmt"
	"log"
	// "regexp"
	"time"
	"github.com/go-chi/chi/v5"
	"strconv"
	"strings"
)

func Deploy(){
	log.Println("starting web app")
	r := chi.NewRouter();
	r.Route("/{guild}/{channel}", func(r chi.Router){
		r.Get("/", messageHandler)
		r.Get("/{date}", messageHandlerDate)
		r.Get("/{date}/next", messageHandlerNav)
		r.Get("/{date}/prev", messageHandlerNav)
	})
	// r.Get("/{guild}/{channel}/*", messageHandler)
	http.ListenAndServe(":8000", r) 
}

func messageHandler(w http.ResponseWriter, r *http.Request) {

	guildParam := strings.TrimSpace(chi.URLParam(r, "guild"))
	channelParam := strings.TrimSpace(chi.URLParam(r, "channel"))

	date_unix := int(time.Now().Unix())
	
	tmpl,err := template.ParseFiles("web/channel.html")
	if err != nil{
		log.Println(err)
	}
	log.Println("-----")
	log.Println(guildParam)
	log.Println(channelParam)
	log.Println(date_unix)
	_, msgs := getMessages(db, guildParam, channelParam, date_unix)
	for _, i := range(msgs.Messages){
		log.Println(i.Message_id)
	}
	tmpl.Execute(w, *msgs)

}

func messageHandlerDate(w http.ResponseWriter, r *http.Request){
	guildParam := strings.TrimSpace(chi.URLParam(r, "guild"))
	channelParam := strings.TrimSpace(chi.URLParam(r, "channel"))
	dateParam := strings.TrimSpace(chi.URLParam(r, "*"))

	date_unix, err := strconv.Atoi(dateParam)
	if err != nil{
		date_unix = int(time.Now().Unix())
	}
	tmpl,err := template.ParseFiles("web/channel.html")
	if err != nil{
		log.Println(err)
	}
	log.Println("-----")
	log.Println(guildParam)
	log.Println(channelParam)
	log.Println(date_unix)
	_, msgs := getMessages(db, guildParam, channelParam, date_unix)
	for _, i := range(msgs.Messages){
		log.Println(i.Message_id)
	}
	tmpl.Execute(w, *msgs)
}

func messageHandlerNav(w http.ResponseWriter, r *http.Request){
	guildParam := strings.TrimSpace(chi.URLParam(r, "guild"))
	channelParam := strings.TrimSpace(chi.URLParam(r, "channel"))
	dateParam := strings.TrimSpace(chi.URLParam(r, "date"))

	date_unix, err := strconv.Atoi(dateParam)
	if err != nil{
		date_unix = int(time.Now().Unix())
	}
	tmpl,err := template.ParseFiles("web/channel.html")
	if err != nil{
		log.Println(err)
	}
	log.Println("-----")
	log.Println(guildParam)
	log.Println(channelParam)
	log.Println(date_unix)
	_, msgs := getMessages(db, guildParam, channelParam, date_unix)
	// for _, i := range(msgs.Messages){
	// 	log.Println(i.Message_id)
	// }
	tmpl.ExecuteTemplate(w,"msgs", *msgs)
}

//need to somehow use ajax when going to next or prev page
//when click next page, get date of last message, send via POST, server returns new URL and navigate to that URL. Update prev page button with old URL
//at beginning: next has date of last message. prev has no date 
//next x1: next has date of last message. prev has no date
//next x2: next has date of last message. prev has next x1's next
//next x3: next has date of last message. prev has next x2's next
//when click prev page, navigate to prev page using updated URL

//click next/prev, 