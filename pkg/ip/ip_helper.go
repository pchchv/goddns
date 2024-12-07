package ip

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

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

func (helper *IPHelper) getIPFromMikrotik() string {
	u, err := url.Parse(helper.configuration.Mikrotik.Addr)
	if err != nil {
		log.Fatal("fail to parse mikrotik address: ", err)
		return ""
	}

	u.Path = path.Join(u.Path, "/rest/ip/address")
	q := u.Query()
	q.Add("interface", helper.configuration.Mikrotik.Interface)
	q.Add(".proplist", "address")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	auth := fmt.Sprintf("%s:%s", helper.configuration.Mikrotik.Username, helper.configuration.Mikrotik.Password)
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout:   time.Second * utils.DefaultTimeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	response, err := client.Do(req)
	if err != nil {
		log.Fatal("request mikrotik address failed:", err)
		return ""
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatal("request code failed: ", response.StatusCode)
		return ""
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("read body failed: ", err)
		return ""
	}

	m := []map[string]string{}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Fatal("unmarshal body failed: ", err)
		return ""
	} else if len(m) < 1 {
		log.Fatal("could not get ip from response: ", m)
		return ""
	}

	res := strings.Split(m[0]["address"], "/")
	return res[0]
}

func (helper *IPHelper) getNext() string {
	newIdx := atomic.AddInt64(&helper.idx, 1)
	helper.mutex.RLock()
	defer helper.mutex.RUnlock()
	newIdx %= int64(len(helper.reqURLs))
	return helper.reqURLs[newIdx]
}

// getIPOnline gets public IP from internet.
func (helper *IPHelper) getIPOnline() (onlineIP string) {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
			proto := "tcp"
			if strings.ToUpper(helper.configuration.IPType) == utils.IPV4 {
				// Force the network to "tcp4" to use only IPv4
				proto = "tcp4"
			}

			return (&net.Dialer{
				Timeout:   time.Second * utils.DefaultTimeout,
				KeepAlive: 30 * time.Second,
			}).DialContext(ctx, proto, addr)
		},
	}
	client := &http.Client{
		Timeout:   time.Second * utils.DefaultTimeout,
		Transport: transport,
	}
	for {
		reqURL := helper.getNext()
		req, _ := http.NewRequest("GET", reqURL, nil)
		if helper.configuration.UserAgent != "" {
			req.Header.Set("User-Agent", helper.configuration.UserAgent)
		}

		response, err := client.Do(req)
		if err != nil {
			log.Fatal("Cannot get IP:", err)
			time.Sleep(time.Millisecond * 300)
			continue
		}

		if response.StatusCode != http.StatusOK {
			log.Fatal(fmt.Sprintf("request %v got httpCode:%v", reqURL, response.StatusCode))
			continue
		}

		body, _ := io.ReadAll(response.Body)
		ipReg := regexp.MustCompile(utils.IPPattern)
		onlineIP = ipReg.FindString(string(body))
		if onlineIP == "" {
			log.Fatal(fmt.Sprintf("request:%v failed to get online IP", reqURL))
			continue
		}

		if isIPv4(onlineIP) && strings.ToUpper(helper.configuration.IPType) != utils.IPV4 {
			log.Fatalf("The online IP (%s) from %s is not IPV6, will skip it.", onlineIP, reqURL)
			continue
		} else if strings.ToUpper(helper.configuration.IPType) != utils.IPV6 {
			log.Fatalf("The online IP (%s) from %s is not IPV4, will skip it.", onlineIP, reqURL)
			continue
		}

		log.Printf("Get ip success by: %s, online IP: %s", reqURL, onlineIP)

		if err = response.Body.Close(); err != nil {
			log.Fatal(fmt.Sprintf("request:%v failed to get online IP", reqURL))
			continue
		}

		if onlineIP == "" {
			log.Fatal("fail to get online IP")
		}

		break
	}

	return
}

func isIPv4(ip string) bool {
	return strings.Count(ip, ":") < 2
}
