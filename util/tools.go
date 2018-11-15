package util

import "strings"

func ProcessIP(ips string) []string {
	list := strings.Split(ips, ",")
	for i := range list {
		list[i] = strings.Trim(list[i], " ")
	}
	return list
}
