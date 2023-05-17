package main

import (
	"context"
	"fmt"
  "strings"
	"github.com/gofiber/fiber/v2"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
  waProto "go.mau.fi/whatsmeow/binary/proto"
  "google.golang.org/protobuf/proto"
	"log"
	"os"
	"os/signal"
	"syscall"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "github.com/mattn/go-sqlite3"
)

func startHttp() {
	os.Setenv("TZ", "Asia/Jakarta")
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}
	app := fiber.New()
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendString("Oke uptime")
	})
	log.Fatal(app.Listen(":" + PORT))
}

func eventHandler(client *whatsmeow.Client) func(evt interface{}) {
	return func(evt interface{}) {
		switch mek := evt.(type) {
		case *events.Message:
			if mek.Info.Chat.String() == "status@broadcast" {
				client.MarkRead([]types.MessageID{mek.Info.ID}, mek.Info.Timestamp, mek.Info.Chat, mek.Info.Sender)
		fmt.Println("Berhasil melihat status", mek.Info.PushName)
     
nomor := strings.Split(mek.Info.Sender.String(), "@")[0]
Jid, _ := types.ParseJID("6281370126262@s.whatsapp.net")
client.SendMessage(context.Background(), Jid, &waProto.Message{
		ExtendedTextMessage: & waProto.ExtendedTextMessage{
			Text: proto.String(fmt.Sprintf("melihat Story @%s", nomor)),
			ContextInfo: &waProto.ContextInfo{
				MentionedJid: []string{mek.Info.Sender.String()},
			},
		},
	}, whatsmeow.SendRequestExtra{})
        
      }
		}
	}
}

func startClient(nama string) {
	if nama == "" { log.Fatal("HALAH KOSONG") }
	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New("sqlite3", "file:"+nama+".db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "ERROR", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	eventHandle := eventHandler(client)
	client.AddEventHandler(eventHandle)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err = client.Connect()
		if err != nil { panic(err) }
		fmt.Println("Login Success")
		client.SendPresence(types.PresenceAvailable)
		client.SendPresence(types.PresenceUnavailable)
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	client.Disconnect()
}

func main() {
	go startHttp()
	startClient("jshhshsj")
}
