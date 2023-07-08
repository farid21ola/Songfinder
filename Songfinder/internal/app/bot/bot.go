package bot

import (
	"Songfinder/internal/app/bot/vkApi"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"os"
	"regexp"
	"strings"
	"time"
)

var msgRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG) #\w+(?: :(.*))?$`)
var cmdRegex *regexp.Regexp = regexp.MustCompile(`^!(\w+)\s?(\w+)?`)

const PSTFormat = "Jan 2 15:04:05 PST"

type OAuthCred struct {
	Password string `json:"password,omitempty"`
}

type BasicBot struct {
	ChannelName string
	conn        net.Conn
	Credentials *OAuthCred
	MsgRate     time.Duration
	Name        string
	Port        string
	PrivatePath string
	Server      string
	startTime   time.Time
}

type TwitchBot interface {
	Connect()
	Disconnect()
	HandleChat() error
	JoinChannel()
	ReadCredentials() error
	Say(msg string) error
	Start()
}

func (bb *BasicBot) Connect() {
	var err error
	fmt.Printf("[%s] Connecting to %s...\n", timeStamp(), bb.Server)

	bb.conn, err = net.Dial("tcp", bb.Server+":"+bb.Port)
	if nil != err {
		fmt.Printf("[%s] Cannot connect to %s, retrying.\n", timeStamp(), bb.Server)
		bb.Connect()
		return
	}
	fmt.Printf("[%s] Connected to %s!\n", timeStamp(), bb.Server)
	bb.startTime = time.Now()
}

func (bb *BasicBot) Disconnect() {
	bb.conn.Close()
	upTime := time.Now().Sub(bb.startTime).Seconds()
	fmt.Printf("[%s] Closed connection from %s! | Live for: %fs\n", timeStamp(), bb.Server, upTime)
}

func (bb *BasicBot) HandleChat() error {
	fmt.Printf("[%s] Watching #%s...\n", timeStamp(), bb.ChannelName)

	tp := textproto.NewReader(bufio.NewReader(bb.conn))

	for {
		line, err := tp.ReadLine()
		if nil != err {

			bb.Disconnect()

			return errors.New("bb.Bot.HandleChat: Failed to read line from channel. Disconnected.")
		}

		//fmt.Printf("[%s] %s\n", timeStamp(), line)

		if "PING :tmi.twitch.tv" == line {

			bb.conn.Write([]byte("PONG :tmi.twitch.tv\r\n"))
			continue
		} else {

			matches := msgRegex.FindStringSubmatch(line)
			if nil != matches {
				userName := matches[1]
				msgType := matches[2]

				switch msgType {
				case "PRIVMSG":
					msg := matches[3]
					fmt.Printf("[%s] %s: %s\n", timeStamp(), userName, msg)

					cmdMatches := cmdRegex.FindStringSubmatch(msg)
					if nil != cmdMatches {
						cmd := cmdMatches[1]

						switch cmd {
						case "tbdown":
							bb.tbdown(userName)
							return nil
						case "music":
							if e := bb.music(); e != nil {
								fmt.Printf("~~~~[%s] can't request song: %s \n", timeStamp(), e.Error())
							}
						default:

						}
					}
				default:

				}
			}
		}
		time.Sleep(bb.MsgRate)
	}
}

func (bb *BasicBot) JoinChannel() {
	fmt.Printf("[%s] Joining #%s...\n", timeStamp(), bb.ChannelName)
	bb.conn.Write([]byte("PASS " + bb.Credentials.Password + "\r\n"))
	bb.conn.Write([]byte("NICK " + bb.Name + "\r\n"))
	bb.conn.Write([]byte("JOIN #" + bb.ChannelName + "\r\n"))

	fmt.Printf("[%s] Joined #%s as @%s!\n", timeStamp(), bb.ChannelName, bb.Name)
}

func (bb *BasicBot) ReadCredentials() error {

	credFile, err := os.ReadFile(bb.PrivatePath)
	if nil != err {
		return err
	}

	bb.Credentials = &OAuthCred{}

	dec := json.NewDecoder(strings.NewReader(string(credFile)))
	if err = dec.Decode(bb.Credentials); nil != err && io.EOF != err {
		return err
	}

	return nil
}

func (bb *BasicBot) Say(msg string) error {
	if "" == msg {
		return errors.New("BasicBot.Say: msg was empty.")
	}
	_, err := bb.conn.Write([]byte(fmt.Sprintf("PRIVMSG #%s :%s\r\n", bb.ChannelName, msg)))
	if nil != err {
		return err
	}
	return nil
}
func (bb *BasicBot) tbdown(userName string) {
	if userName == bb.ChannelName {
		fmt.Printf(
			"[%s] Shutdown command received. Shutting down now...\n",
			timeStamp(),
		)
		bb.Disconnect()
	}
}

func (bb *BasicBot) music() error {
	str, e := vkApi.SearchMusic()
	if e != nil {
		return e
	}
	err := bb.Say(str)
	if nil != err {
		return err
	} else {
		fmt.Printf(
			"[%s] Sent message: %s \n",
			timeStamp(),
			str,
		)

		return nil
	}
}

func (bb *BasicBot) Start() {
	err := bb.ReadCredentials()
	if nil != err {
		fmt.Println(err)
		fmt.Println("Aborting...")
		return
	}

	for {
		bb.Connect()
		bb.JoinChannel()
		err = bb.HandleChat()
		if nil != err {
			// attempts to reconnect upon unexpected chat error
			time.Sleep(1000 * time.Millisecond)
			fmt.Println(err)
			fmt.Println("Starting bot again...")
		} else {
			return
		}
	}
}

func timeStamp() string {
	return TimeStamp(PSTFormat)
}

func TimeStamp(format string) string {
	return time.Now().Format(format)
}
