package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/unpoller/unifi"
)

func main() {
	// Construct a new API object
	api, err := cloudflare.New(os.Getenv("CLOUDFLARE_API_KEY"), os.Getenv("CLOUDFLARE_API_EMAIL"))
	if err != nil {
		log.Fatalln("Can't get env vars:", err)
	}

	// make unifi object
	c := unifi.Config{
		User: os.Getenv("UNIFI_USER"),
		Pass: os.Getenv("UNIFI_PWD"),
		URL:  os.Getenv("UNIFI_URL"),
		// Log with log.Printf or make your own interface that accepts (msg, fmt)
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}
	uni, err := unifi.NewUnifi(&c)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	lbs := getLBs(api)
	ips := getUnifiWANIPs(uni)

	// see if they match up
	wan1 := os.Getenv("WAN1_NAME")
	wan2 := os.Getenv("WAN2_NAME")

	// save updates
	var updateLBs []cloudflare.LoadBalancerPool

	for _, lb := range lbs {
		if strings.EqualFold(lb.Name, wan1) {
			if strings.EqualFold(lb.Origins[0].Address, ips[0]) {
				fmt.Println("wan1 ip matches, no change needed")
			} else {
				lb.Origins[0].Address = ips[0]
				updateLBs = append(updateLBs, lb)
				fmt.Printf("wan1 ip does not match, updating to %s\n", ips[0])
			}
		} else {
			if strings.EqualFold(lb.Name, wan2) {
				if strings.EqualFold(lb.Origins[0].Address, ips[1]) {
					fmt.Println("wan2 ip matches, no change needed")
				} else {
					lb.Origins[0].Address = ips[1]
					updateLBs = append(updateLBs, lb)
					fmt.Printf("wan1 ip does not match, updating to %s\n", ips[1])
				}
			}
		}
	}

	for _, lb := range updateLBs {
		api.ModifyLoadBalancerPool(context.Background(), lb)
	}
}

func getLBs(api *cloudflare.API) []cloudflare.LoadBalancerPool {
	// List LB pools
	lbPoolList, err := api.ListLoadBalancerPools(context.Background())
	if err != nil {
		log.Println("Can't get lb pools")
		log.Fatal(err)
	}

	for _, pool := range lbPoolList {
		fmt.Printf("pool: %s, address: %s\n", pool.Name, pool.Origins[0].Address)
	}
	return lbPoolList
}

// returns the wan1 and wan2 ip in that order
func getUnifiWANIPs(uni *unifi.Unifi) (ips [2]string) {
	sites, err := uni.GetSites()
	if err != nil {
		log.Fatalln("Can't get sites:", err)
	}

	devices, err := uni.GetDevices(sites)
	if err != nil {
		log.Fatalln("Can't get devices: ", err)
	}
	for _, usg := range devices.USGs {
		ips[0] = usg.Wan1.IP
		fmt.Printf("WAN1: %s\n", usg.Wan1.IP)
		ips[1] = usg.Wan2.IP
		fmt.Printf("WAN2: %s\n", usg.Wan2.IP)
	}
	return

	// clients, err := uni.GetClients(sites)
	// if err != nil {
	// 	log.Fatalln("Error:", err)
	// }
	// devices, err := uni.GetDevices(sites)
	// if err != nil {
	// 	log.Fatalln("Error:", err)
	// }

	// log.Println(len(sites), "Unifi Sites Found: ", sites)
	// log.Println(len(clients), "Clients connected:")
	// for i, client := range clients {
	// 	log.Println(i+1, client.ID, client.Hostname, client.IP, client.Name, client.LastSeen)
	// }

	// log.Println(len(devices.USWs), "Unifi Switches Found")
	// log.Println(len(devices.USGs), "Unifi Gateways Found")

	// log.Println(len(devices.UAPs), "Unifi Wireless APs Found:")
	// for i, uap := range devices.UAPs {
	// 	log.Println(i+1, uap.Name, uap.IP)
	// }
}
