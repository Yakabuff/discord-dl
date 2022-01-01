package main
import(
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"io"
	"encoding/json"
)
func dm_download(dg *discordgo.Session, a args) error {
	//get all channels from guild into array
	//use endpoint https://discord.com/api/v8/users/@me/channels
	//https://discord.com/channels/@me
	err, channels := get_dms_request(dg, a)
	if err != nil{
		log.Println("Could not find guild")
		return err
	}
	//download messages from every channel
	for _, c := range channels{
		a.channel = c.ID
		log.Printf("Downloading DM channel: %s", a.channel)
		err := channel_download(dg, a)
		if err != nil{
			return err;
		}
	}
	

	return nil;
}

func get_dms_request(dg *discordgo.Session, a args) (error, []discordgo.Channel) {
//send get request headers = {'authorization': token, 'content-type': 'application/json'}
	req, err := http.NewRequest("GET", "https://discord.com/api/v8/users/@me/channels", nil)
	if err != nil{
		log.Println(err)
	}

	req.Header.Set("authorization", a.token)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		log.Println(err)
	}
	bodybytes, err := io.ReadAll(resp.Body)
	if err != nil{
		log.Println(err)
	}
	bodyString := string(bodybytes)

	var dms []DM

	if err := json.Unmarshal([]byte(bodyString), &dms); err != nil {
        panic(err)
    }
	var channels []discordgo.Channel
	for _, i := range dms{
		c, err := dg.Channel(i.Id)
		if err != nil{
			return err, nil
		}
		channels = append(channels, *c)
	}

	return err, channels
}

type DM struct {
	Id string `json:"id"`

}
