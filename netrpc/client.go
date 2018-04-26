package netrpc

import (
	"log"
	"net/rpc"
	"strings"
)

type Client interface {
	Call(serviceMethod string, args interface{}, reply interface{}) error
}

type rpcClient struct {
	address string
	path    string
	*rpc.Client
}

func (c *rpcClient) Call(serviceMethod string, args interface{}, reply interface{}) error {
	if c.Client == nil {
		newClient, err := rpc.DialHTTPPath("tcp", c.address, c.path)
		if err != nil {
			log.Println("net/rpc dial error.", err)
			return err
		}
		c.Client = newClient
	}
	err := c.Client.Call(serviceMethod, args, reply)
	if err == nil {
		return nil
	}

	log.Println("net/rpc call error.", err)

	c.Client = nil

	if err != rpc.ErrShutdown {
		return err
	}

	newClient, err := rpc.DialHTTPPath("tcp", c.address, c.path)
	if err != nil {
		log.Println("net/rpc dial retry error.", err)
		return err
	}

	log.Println("net/rpc dial retry success.")

	c.Client = newClient

	if err := c.Client.Call(serviceMethod, args, reply); err != nil {
		log.Println("net/rpc call retry error.", err)
		c.Client = nil
		return err
	}
	log.Println("net/rpc call retry success.")
	return nil
}

func NewClient(addr string) Client {
	index := strings.Index(addr, "/")

	if index <= 0 {
		return &rpcClient{
			address: addr,
			path:    rpc.DefaultRPCPath,
		}
	}
	return &rpcClient{
		address: addr[:index],
		path:    addr[index:] + rpc.DefaultRPCPath,
	}

}
