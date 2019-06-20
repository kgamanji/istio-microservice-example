package main

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	proto "github.com/AlexsJones/golang-microservice-example/protocolbuffers"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"strings"
	"time"

)

type server struct{}

func (*server)SendMessage(c context.Context,r *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {

	if r.Message == "" {
		return nil,errors.New("bad message")
	}
	log.Println(r.Message)
	response := "Nada"
	str := strings.Split(r.Message,":")
	if len(str) > 1 {
		response = "Pong number " + str[1]
	}

	return &proto.SendMessageResponse{Response:response},nil
}

func serverStart(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Warn("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterMessageServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Warn("failed to serve: %v", err)
	}
}

func client(address string, message string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := proto.NewMessageClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SendMessage(ctx, &proto.SendMessageRequest{Message: message})
	if err != nil {
		return err
	}
	log.Printf("Response: %s", r.Response)
	return nil
}
func clientPulse() {

	count := 0
	for {
		time.Sleep(time.Second * 1)
		err := client(Options.TargetAddress,fmt.Sprintf("Sending ping:%d",count))
		if err != nil {
			log.Warn(err.Error())
		}
		count++
	}
}
var Options struct {
	TargetAddress string `short:"t" long:"targetAddress e.g. localhost:12701" required:"true"`
	ServerPort string `short:"s" long:"serverPort e.g.0.0.0.0:9000" required:"true"`
}
func main() {
	// Set up a connection to the server.
	_, err := flags.ParseArgs(&Options,os.Args)
	if err != nil {
		panic(err)
	}

	 go clientPulse()

	serverStart(Options.ServerPort)
}