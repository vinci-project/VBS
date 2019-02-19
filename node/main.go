package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
)

type storageNode struct {
	HardwareInfo  HardwareInfo
	NodeId        string
	WalletAddress string
	StoragePath   string
}

var storagePath string
var walletAddress string
var node storageNode
var runCheck bool

func main() {

	app := cli.NewApp()

	// Initialization
	app.Name = "VBS Storage Node"
	app.Version = "0.0.1"
	app.Compiled = time.Now()

	app.Usage = "The decentralized backup service"
	app.UsageText = "VBS Node - The decentralized backup service"

	// Read In Flag
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "storagePath",
			Value:       "/var/lib/vbs",
			Usage:       "Storage Path",
			Destination: &storagePath,
		},
		cli.StringFlag{
			Name:        "walletAddress",
			Value:       "",
			Usage:       "Wallet Address",
			Destination: &walletAddress,
		},
		cli.BoolFlag{
			Name:        "hardware",
			Usage:       "Run Hardware Check",
			Destination: &runCheck,
		},
	}

	app.Action = func(c *cli.Context) error {
		err := runApp()
		return err
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runApp() error {
	fmt.Println("Initializing.....")
	fmt.Printf("Wallet Address: %s\n", walletAddress)
	fmt.Printf("Storage path: %s\n", storagePath)
	fmt.Printf("-------------------------------------------------------\n")

	node = storageNode{
		HardwareInfo:  HardwareInfo{},
		NodeId:        "",
		WalletAddress: walletAddress,
		StoragePath:   storagePath,
	}

	if runCheck == true {
		node.getHardwareInformation()
		node.checkRequirement()
	}

	return nil
}
