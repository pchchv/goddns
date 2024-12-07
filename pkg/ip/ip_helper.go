package ip

import (
	"log"
	"strings"
	"sync"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

type IPHelper struct {
	reqURLs       []string
	currentIP     string
	mutex         sync.RWMutex
	configuration *settings.Settings
	idx           int64
}

func (helper *IPHelper) UpdateConfiguration(conf *settings.Settings) {
	helper.mutex.Lock()
	defer helper.mutex.Unlock()

	// clear urls
	helper.reqURLs = helper.reqURLs[:0]
	// reset the index
	helper.idx = -1

	if conf.IPType == "" || strings.ToUpper(conf.IPType) == utils.IPV4 {
		// filter empty urls
		for _, url := range conf.IPUrls {
			if url != "" {
				helper.reqURLs = append(helper.reqURLs, url)
			}
		}

		if conf.IPUrl != "" {
			helper.reqURLs = append(helper.reqURLs, conf.IPUrl)
		}
	} else {
		// filter empty urls
		for _, url := range conf.IPV6Urls {
			if url != "" {
				helper.reqURLs = append(helper.reqURLs, url)
			}
		}

		if conf.IPV6Url != "" {
			helper.reqURLs = append(helper.reqURLs, conf.IPV6Url)
		}
	}

	log.Printf("Update ip helper configuration, urls: %v", helper.reqURLs)
}
