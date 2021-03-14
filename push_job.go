package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/appleboy/gorush/rpc/proto"
	"github.com/gocraft/work"
	"google.golang.org/grpc"
)

func (c *Context) PushToMobile(job *work.Job) error {
	title := job.ArgString("title")
	body := job.ArgString("body")
	var tokens []string

	json.Unmarshal([]byte(job.ArgString("tokens")), &tokens)

	platform, _ := strconv.ParseInt(job.ArgString("platform"), 10, 32)
	if err := job.ArgError(); err != nil {
		return err
	}

	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := proto.NewGorushClient(conn)

	r, err := client.Send(context.Background(), &proto.NotificationRequest{
		Platform: int32(platform),
		Tokens:   tokens,
		Message:  body,
		Badge:    1,
		Alert: &proto.Alert{
			Title: title,
			Body:  body,
		},
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Success: %t\n", r.Success)
	log.Printf("Count: %d\n", r.Counts)

	return nil
}
