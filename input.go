package main
import(
	"encoding/json"
	"os"
	"log"
)
type config struct {
	before string
	after string
	fast_update bool
	token string
	output string
	guild string
	channel string
	dms bool
	listen bool
	deploy bool
	media_location string
 }
func parse_input(file_name string) error{
	//read file 
	file, err := os.Open(file_name)
	if err != nil{
		return err
	}
	defer file.Close()
	//parse each element of json as an args type
	decoder := json.NewDecoder(file);
	config := config{}
	err = decoder.Decode(&config);
	if err != nil{
		return err
	}
	log.Println(config)
	//execute instructions in order
	return nil;
}