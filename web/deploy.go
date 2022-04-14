package web

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yakabuff/discord-dl/db"
	"github.com/yakabuff/discord-dl/models"
)

type Web struct {
	db            db.Db
	port          int
	mediaLocation string
}

func NewWeb(db db.Db, port int, mediaLocation string) Web {
	web := Web{}
	web.db = db
	web.port = port
	web.mediaLocation = mediaLocation
	return web
}

func (web Web) Deploy(Db db.Db) {
	go func(Db db.Db) {
		log.Println("starting web app on port: " + strconv.Itoa(web.port))
		r := chi.NewRouter()
		r.Route("/index", func(r chi.Router) {
			r.Get("/", web.guildHandler)
			r.Get("/{guild}", web.channelHandler)
		})

		r.Route("/{guild}/{channel}", func(r chi.Router) {
			//First 100 messages
			r.Get("/", web.messageHandler)
			//100 messages after specified date
			r.Get("/{date}", web.messageHandlerDate)
			//fetch next 100 messages ( date + 1) or fetch previous 100 messages ( date -1)
			r.Get("/{date}/{nav}", web.messageHandlerNav)
		})

		r.Get("/media/{channel}/{hash}", mediaHandler)
		http.ListenAndServe(":"+strconv.Itoa(web.port), r)
	}(Db)
}

func mediaHandler(w http.ResponseWriter, r *http.Request) {
	channelParam := strings.TrimSpace(chi.URLParam(r, "channel"))
	hashParam := strings.TrimSpace(chi.URLParam(r, "hash"))
	path := filepath.FromSlash("media/" + channelParam + "/" + hashParam)
	http.ServeFile(w, r, path)
}

func (web Web) guildHandler(w http.ResponseWriter, r *http.Request) {

	guilds, err := web.db.GetAllGuilds()

	if err != nil {
		log.Println(err)
	}
	web.addGuildMetadataResourceLink(guilds)

	g := models.Guilds{Guilds: guilds}
	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Println(err)
	}
	tmpl.Execute(w, g)
}

func (web Web) channelHandler(w http.ResponseWriter, r *http.Request) {
	guildParam := strings.TrimSpace(chi.URLParam(r, "guild"))
	channels, err := web.db.GetChannelsFromGuild(guildParam)
	if err != nil {
		log.Println(err)
	}
	c := models.Channels{Channels: channels}
	tmpl, err := template.ParseFiles("static/channels.html")
	if err != nil {
		log.Println(err)
	}
	tmpl.Execute(w, c)
}

func (web Web) messageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + " " + r.URL.Path)
	guildParam := strings.TrimSpace(chi.URLParam(r, "guild"))
	channelParam := strings.TrimSpace(chi.URLParam(r, "channel"))

	date_unix := int(time.Now().Unix())

	tmpl, err := template.ParseFiles("static/channel.html")
	if err != nil {
		log.Println(err)
	}

	err, msgs := web.db.GetMessages(guildParam, channelParam, date_unix, false)
	if err != nil {
		log.Println(err)
	}
	for j, i := range msgs.Messages {
		web.addEmbedResourceLink(i.Embeds, channelParam)
		web.addAttachmentResourceLink(i.Attachments, channelParam)
		if msgs.Messages[j].ThreadId != "" {
			msgs.Messages[j].ThreadPath = filepath.FromSlash("/" + i.GuildId + "/" + i.ThreadId + "/")
		}
	}

	tmpl.Execute(w, *msgs)

}

func (web Web) messageHandlerDate(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + " " + r.URL.Path)
	guildParam := strings.TrimSpace(chi.URLParam(r, "guild"))
	channelParam := strings.TrimSpace(chi.URLParam(r, "channel"))
	dateParam := strings.TrimSpace(chi.URLParam(r, "*"))

	date_unix, err := strconv.Atoi(dateParam)
	if err != nil {
		date_unix = int(time.Now().Unix())
	}
	tmpl, err := template.ParseFiles("static/channel.html")
	if err != nil {
		log.Println(err)
	}

	err, msgs := web.db.GetMessages(guildParam, channelParam, date_unix, false)
	if err != nil {
		log.Println(err)
	}
	for j, i := range msgs.Messages {
		web.addEmbedResourceLink(i.Embeds, channelParam)
		web.addAttachmentResourceLink(i.Attachments, channelParam)
		if msgs.Messages[j].ThreadId != "" {
			msgs.Messages[j].ThreadPath = filepath.FromSlash("/" + i.GuildId + "/" + i.ThreadId + "/")
		}
	}
	tmpl.Execute(w, *msgs)
}

func (web Web) messageHandlerNav(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + " " + r.URL.Path)
	guildParam := strings.TrimSpace(chi.URLParam(r, "guild"))
	channelParam := strings.TrimSpace(chi.URLParam(r, "channel"))
	dateParam := strings.TrimSpace(chi.URLParam(r, "date"))
	afterParam := strings.TrimSpace(chi.URLParam(r, "nav"))

	date_unix, err := strconv.Atoi(dateParam)
	if err != nil {
		date_unix = int(time.Now().Unix())
	}
	tmpl, err := template.ParseFiles("static/channel.html")
	if err != nil {
		log.Println(err)
	}

	var msgs *models.Messages
	if afterParam == "next" {
		err, msgs = web.db.GetMessages(guildParam, channelParam, date_unix, true)
		if err != nil {
			log.Println(err)
		}
	} else if afterParam == "prev" {
		err, msgs = web.db.GetMessages(guildParam, channelParam, date_unix, false)
		if err != nil {
			log.Println(err)
		}
	}
	if len(msgs.Messages) != 0 {
		for j, i := range msgs.Messages {
			web.addEmbedResourceLink(i.Embeds, channelParam)
			web.addAttachmentResourceLink(i.Attachments, channelParam)
			if msgs.Messages[j].ThreadId != "" {
				msgs.Messages[j].ThreadPath = filepath.FromSlash("/" + i.GuildId + "/" + i.ThreadId + "/")
			}
		}
	}

	tmpl.ExecuteTemplate(w, "msgs", *msgs)
}

func (web Web) addEmbedResourceLink(embeds []models.EmbedOut, channel_id string) {
	for i, _ := range embeds {
		if embeds[i].EmbedImageHash != "" {
			embeds[i].ResourcePathImage = filepath.FromSlash("/" + web.mediaLocation + "/" + channel_id + "/" + embeds[i].EmbedImageHash)
		}
		if embeds[i].EmbedThumbnailHash != "" {
			embeds[i].ResourcePathThumbnail = filepath.FromSlash("/" + web.mediaLocation + "/" + channel_id + "/" + embeds[i].EmbedThumbnailHash)
		}
		if embeds[i].EmbedVideoHash != "" {
			embeds[i].ResourcePathVideo = filepath.FromSlash("/" + web.mediaLocation + "/" + channel_id + "/" + embeds[i].EmbedVideoHash)
		}
	}
}

func (web Web) addGuildMetadataResourceLink(guilds []models.GuildOut) {
	for i := range guilds {
		guilds[i].GuildBannerResourcePath = filepath.FromSlash("/" + web.mediaLocation + "/" + guilds[i].GuildID + "/" + guilds[i].BannerHash)
		guilds[i].GuildIconResourcePath = filepath.FromSlash("/" + web.mediaLocation + "/" + guilds[i].GuildID + "/" + guilds[i].IconHash)
	}

}

func (web Web) addAttachmentResourceLink(attachments []models.AttachmentOut, channel_id string) {
	for i, _ := range attachments {
		attachments[i].ResourcePath = filepath.FromSlash("/" + web.mediaLocation + "/" + channel_id + "/" + attachments[i].AttachmentHash)
		s := strings.ToLower(strings.Split(attachments[i].AttachmentFilename, ".")[1])
		if s == "png" || s == "jpg" || s == "jpeg" || s == "gif" {
			attachments[i].ResourceType = "IMAGE"
		} else if s == "mp4" || s == "mov" || s == "wmv" || s == "avi" || s == "flv" || s == "swf" || s == "mkv" || s == "webm" {
			attachments[i].ResourceType = "VIDEO"
		} else {
			attachments[i].ResourceType = "FILE"
		}
	}
}

//need to somehow use ajax when going to next or prev page
//when click next page, get date of last message, send via POST, server returns new URL and navigate to that URL. Update prev page button with old URL
//at beginning: next has date of last message. prev has no date
//next x1: next has date of last message. prev has no date
//next x2: next has date of last message. prev has next x1's next
//next x3: next has date of last message. prev has next x2's next
//when click prev page, navigate to prev page using updated URL

//click next/prev,
