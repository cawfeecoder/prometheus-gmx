package collector

import (
	"fmt"
	"github.com/nfrush/prometheus-gmx/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/thoas/go-funk"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

func ScrapeTarget(target string, config *config.Config) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(config.HostPort)

	if err != nil {
		return nil, fmt.Errorf("error connecting to target %s: %s", target, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading payload from target %s: %s", target, err)
	}

	return body, nil
}

func ParsePayload(target string, payload []byte) ([]*fastjson.Value, error) {
	json_body, err := fastjson.ParseBytes(payload)
	if err != nil {
		return nil, fmt.Errorf("error malformed payload from target %s: %s", target, err)
	}
	return json_body.GetArray("beans"), nil
}

func FilterPayload(payload []*fastjson.Value, config *config.Config) ([]*fastjson.Value) {
	filtered := funk.Filter(payload, func(x *fastjson.Value) bool {
		if funk.Contains(config.BlacklistObjectNames, x.Get("name").String()){
			return false
		}
		if funk.Contains(config.WhitelistObjectNames, x.Get("name").String()) {
			return true
		}
		blacklist_regex_excludes := false
		funk.ForEach(config.BlacklistObjectNamesRegexp, func(y *regexp.Regexp){
			if y.Match([]byte(x.Get("name").String())) {
				blacklist_regex_excludes = true
			}
		})
		if blacklist_regex_excludes {
			return false
		}
		whitelist_regex_includes := false
		funk.ForEach(config.WhitelistObjectNamesRegexp, func(y *regexp.Regexp){
			if y.Match([]byte(x.Get("name").String())) {
				whitelist_regex_includes = true
			}
		})
		return whitelist_regex_includes
	})
	return filtered.([]*fastjson.Value)
}

func PayloadToSamples(payload []*fastjson.Value) ([]prometheus.Metric) {
	
}