package main

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/nfrush/prometheus-gmx/collector"
	config2 "github.com/nfrush/prometheus-gmx/config"
	"log"
	_ "net/http/pprof"
	"regexp"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var name_regex = regexp.MustCompile(`(.+):[^=]+=(.+)`)
var prop_regex = regexp.MustCompile(`=(.+)`)

func bean_to_metric_name(bean_name string) (string) {
	bean_name = bean_name[1:len(bean_name)-1]
	parts := strings.Split(bean_name, ",")
	raw_parent_name := name_regex.FindStringSubmatch(parts[0])
	raw_parent_name = raw_parent_name[1:len(raw_parent_name)]
	for i := 1; i < len(parts); i++ {
		prop_name := prop_regex.FindStringSubmatch(parts[i])
		raw_parent_name = append(raw_parent_name, prop_name[1])
	}
	return strings.Join(raw_parent_name, "_")
}

func main() {
	config, err := config2.LoadFile("jmx.yml")
	if err != nil {
		log.Fatal(err)
	}
	payload, err := collector.ScrapeTarget("127.0.0.1", config)
	if err != nil {
		log.Fatal(err)
	}
	parsed, err := collector.ParsePayload("127.0.0.1", payload)
	if err != nil {
		log.Fatal(err)
	}
	filtered := collector.FilterPayload(parsed, config)
	fmt.Printf("Payload: %v", filtered)
}
