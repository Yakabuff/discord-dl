package main

import (
	"fmt"
	"flag"
   "os"
   "strings"
   "strconv"
   "github.com/bwmarrin/discordgo"
   "log"
   "database/sql"
   // "time"
   "os/signal"
	"syscall"
)

type Mode int;

const(
   NONE Mode = iota
   INPUT
   GUILD
   CHANNEL
   DMS
   LISTEN
   DEPLOY
)

var db *sql.DB
var err error

func main(){
	var a = init_cli();

   db, err = init_db();
   if a.mode == INPUT{
      fmt.Println("Selected input mode")
      parse_input(a.input);
   }else{
      dg, err := discordgo.New(a.token);

      if err != nil{
         fmt.Println("Error creating discord session");
         return;
      }
      err = dg.Open();

      if err != nil{
         fmt.Println(err.Error());
         return
      }

      u, err := dg.User("@me");

      log.Printf("discord-dl has succesfully logged into %s#%s %s\n", u.Username, u.Discriminator, u.ID);

      if a.listen{
         dg.Identify.Intents = discordgo.IntentsGuildMessages
         log.Println("Listening...")
         dg.AddHandler(messageListen)
         dg.AddHandler(messageUpdateListen)
      }

      switch a.mode{
      case GUILD:
         fmt.Println("Archiving guild")
         err := guild_download(dg, *a)
         if err != nil{
            log.Println(err)
         }
      case CHANNEL:
         fmt.Println("Selected channel mode")
         err := channel_download(dg, *a)
         if err != nil{
            log.Println(err)
         }
      case DMS:
         fmt.Println("Selected DM mode")
         err := dm_download(dg, *a)
         if err != nil{
            log.Println(err)
         }
      }

      if a.listen || a.deploy{
         //If deploy launch new goroutine
         go Deploy()
         // Wait here until CTRL-C or other term signal is received.
         fmt.Println("Bot is now running.  Press CTRL-C to exit.")
         sc := make(chan os.Signal, 1)
         signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
         <-sc
      }
      dg.Close();
   }
}

type args struct {
   mode Mode
   progress bool
   before string
   after string
   fast_update bool
   download_media bool
   token string
   output string
   input string
   guild string
   channel string
   dms bool
   listen bool
   deploy bool
}

func init_cli() *args{
   progress := flag.Bool("progress", false, "Displays progress of task. Cannot be used with listen mode");
   before := flag.String("before", "", "Retrieves all messages before this date or message id");
   after := flag.String("after", "", "Retrieves all messages after this date or message id");
   fast_update := flag.Bool("fast-update", false, "Retrieves all message after the last downloaded message");
   download_media := flag.Bool("download-media", false, "downloads embedded images and files from message");
   token := flag.String("t", "", "Sets user or bot token");
   output := flag.String("o", "", "Sets output db path");

   input := flag.String("i", "", "Input mode. Gets config from input file");
   guild := flag.String("guild", "", "Guild mode. Retrieves messages and channels from selected guild");
   channel := flag.String("channel", "", "Retrieves messages from selected channel");
   dms := flag.Bool("dms", false, "DM mode. Retrieves all DM conversations"); 
   listen := flag.Bool("listen", false, "Listens for new messages/events and archives in real time.  Can only be used with a bot account");
   deploy := flag.Bool("deploy",false, "Deploys webapp")
   flag.Parse();

   mode := check_flag_mode(*input, *guild, *channel, *dms, *listen, *deploy);
   if mode == NONE{
      fmt.Fprintln(os.Stderr,"Invalid flags");
      os.Exit(1);
   }

   a := args{
      progress: *progress,
      before: *before,
       after: *after,
 fast_update: *fast_update,
download_media: *download_media,
       token: *token,
      output: *output,
       input: *input,
       guild: *guild,
     channel: *channel,
         dms: *dms,
      listen: *listen,
      deploy: *deploy,
        mode: mode}

   if(*input != "" && len(os.Args) > 3){
      fmt.Fprintln(os.Stderr,"Option --i cannot be used in conjunction with other flags");
      os.Exit(1);
   }

   if(*guild != "" && *channel != ""){
      fmt.Fprintln(os.Stderr,"Cannot use --guild and --channel together");
      os.Exit(1);
   }

   if((*before != "" || *after != "") && *fast_update != false){
      fmt.Fprintln(os.Stderr,"Cannot have before/after flags with fast-update");
      os.Exit(1);
   }

   if(*before != "" && *after != ""){
      if(strings.Contains(*before, "-") && !strings.Contains(*after, "-") || !strings.Contains(*before, "-") && strings.Contains(*after, "-")){
         fmt.Fprintln(os.Stderr,"Before and after flags must be in the same format");
         os.Exit(1);
      }

      if(!strings.Contains(*before, "-") && !strings.Contains(*after, "-")){
         i, _ := strconv.Atoi(*before);
         fmt.Println(i)
         j, _ := strconv.Atoi(*after);
         fmt.Println(j)
         // if(i <= j || err != nil){
         //    fmt.Fprintln(os.Stderr,"Before value must not be greater or equal to after");
         //    os.Exit(1);
         // }
      }
   }

   return &a;
}

func check_flag_mode(input string, guild string, channel string, dms bool, listen bool, deploy bool) Mode{
   var count int;
   var mode Mode;
   if input != ""{
      count++;
      mode = INPUT;
   }
   if guild != ""{
      count++
      mode = GUILD;
   }
   if channel != ""{
      count++
      mode = CHANNEL;
   }
   if dms != false{
      count++
      mode = DMS;
   }
   if input != ""{
      count ++
      mode = INPUT;
   }
   if count == 1{
      return mode;
   }
   //option to run these flags by itself
   if count == 0 && listen {
      mode = LISTEN
      return mode
   }
   if count == 0 && deploy{
      mode = DEPLOY
      return mode
   }
   mode = NONE;
   return mode;
}
