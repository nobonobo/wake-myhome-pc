package main

import (
	"embed"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"syscall/js"

	"github.com/joho/godotenv"
	"github.com/rtctunnel/rtctunnel/channels"
	"github.com/rtctunnel/rtctunnel/crypt"
	"github.com/rtctunnel/rtctunnel/peer"
	"github.com/rtctunnel/rtctunnel/signal"
)

//go:embed secret.env
var secret embed.FS

type Client struct {
	*rpc.Client
	closer []func() error
}

func NewClient(keypair crypt.KeyPair, peerPublicKey crypt.Key) *Client {
	c := &Client{}
	pc, err := peer.Open(keypair, peerPublicKey)
	if err != nil {
		log.Fatal(err)
	}
	c.closer = append(c.closer, pc.Close)
	conn, err := pc.Open(8080)
	if err != nil {
		log.Fatal(err)
	}
	c.closer = append(c.closer, conn.Close)
	c.Client = jsonrpc.NewClient(conn)
	return c
}

func (c *Client) Close() error {
	for i := len(c.closer) - 1; i >= 0; i-- {
		c.closer[i]()
	}
	return nil
}

func main() {
	fp, err := secret.Open("secret.env")
	if err != nil {
		log.Print(err)
	}
	d, err := godotenv.Parse(fp)
	if err != nil {
		log.Print(err)
	}
	signal.SetDefaultOptions(
		signal.WithChannel(channels.Must(channels.Get("operator://operator.irieda.com"))),
	)
	pub, err := crypt.NewKey(d["PUBLIC_KEY"])
	if err != nil {
		log.Fatal(err)
	}
	priv, err := crypt.NewKey(d["PRIVATE_KEY"])
	if err != nil {
		log.Fatal(err)
	}
	keyPair := crypt.KeyPair{Public: pub, Private: priv}
	peerPub, err := crypt.NewKey(d["PEER_PUBLIC_KEY"])
	if err != nil {
		log.Fatal(err)
	}
	client := NewClient(keyPair, peerPub)
	js.Global().Set("WakeOnLan", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go client.Call("Node.WakeOnLan", "08:BF:B8:8D:32:43", nil)
		return nil
	}))
	defer client.Close()
	select {}
}
