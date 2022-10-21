package process

import (
	"log"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

var client *redis.Client
var sub *redis.PubSub

func init() {
	log.Println("connecting to Redis...")
	client = redis.NewClient(&redis.Options{Addr: "redis:6379"})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal("failed to connect to redis", err)
	}
	log.Println("connected to redis")
	startSubscriber()
}

const channel = "chat"

func startSubscriber() {
	/*
		this goroutine exits when the application shuts down. When the pusub connection is closed,
		the channel range loop terminates, hence terminating the goroutine
	*/
	go func() {
		log.Println("starting subscriber...")
		sub = client.Subscribe(channel)
		messages := sub.Channel()
		for message := range messages {
			from := strings.Split(message.Payload, ":")[0]
			//send to all websocket sessions/peers
			for user, peer := range Peers {
				if from != user { //don't recieve your own messages
					peer.WriteMessage(websocket.TextMessage, []byte(message.Payload))
				}
			}
		}
	}()
}

func SendToChannel(msg string) {
	err := client.Publish(channel, msg).Err()
	if err != nil {
		log.Println("could not publish to channel", err)
	}
}

const users = "users"
func CheckUserExists(user string) (bool, error) {
	usernameTaken, err := client.SIsMember(users, user).Result()
	if err != nil {
		return false, err
	}
	return usernameTaken, nil
}

func CreateUser(user string) error {
	err := client.SAdd(users, user).Err()
	if err != nil {
		return err
	}
	return nil
}

func RemoveUser(user string) {
	err := client.SRem(users, user).Err()
	if err != nil {
		log.Println("failed to remove user:", user)
		return
	}
	log.Println("removed user from redis:", user)
}

func Cleanup() {
	for user, peer := range Peers {
		client.SRem(users, user)
		peer.Close()
	}
	log.Println("cleaned up users and sessions...")
	err := sub.Unsubscribe(channel)
	if err != nil {
		log.Println("failed to unsubscribe redis channel subscription:", err)
	}
	err = sub.Close()
	if err != nil {
		log.Println("failed to close redis channel subscription:", err)
	}

	err = client.Close()
	if err != nil {
		log.Println("failed to close redis connection: ", err)
		return
	}
}