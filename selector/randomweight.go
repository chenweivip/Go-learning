package main

import "fmt"

// 权重随机
//func RandomWeight(services []*registry.Service) Next {
//	nodes := make([]*registry.Node, 0, len(services))
//
//	for _, service := range services {
//		nodes = append(nodes, service.Nodes...)
//	}
//
//	return func() (*registry.Node, error) {
//		if len(nodes) == 0 {
//			return nil, ErrNoneAvailable
//		}
//
//		i := rand.Int() % len(nodes)
//		return nodes[i], nil
//	}
//}

func main() {
	services := make(map[string]int)
	services["192.168.0.1"] = 1
	services["192.168.0.2"] = 8
	services["192.168.0.3"] = 3
	services["192.168.0.4"] = 6
	services["192.168.0.5"] = 5
	services["192.168.0.6"] = 5
	services["192.168.0.7"] = 4
	services["192.168.0.8"] = 7
	services["192.168.0.9"] = 2
	services["192.168.0.10"] = 9

	ip_list := make([]string, 10)

	// 按照权重复制
	for ip, weight := range services {
		for i := 0; i < weight; i++ {
			ip_list = append(ip_list, ip)
		}
	}
	fmt.Println(len(ip_list))
}