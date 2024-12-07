package ip

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
	"github.com/pchchv/goddns/pkg/safe"
)

var (
	helperInstance *IPHelper
	helperOnce     sync.Once
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

func (helper *IPHelper) GetCurrentIP() string {
	// first load
	if helper.currentIP == "" {
		helper.getCurrentIP()
	}

	helper.mutex.RLock()
	defer helper.mutex.RUnlock()

	return helper.currentIP
}

func GetIPHelperInstance(conf *settings.Settings) *IPHelper {
	helperOnce.Do(func() {
		helperInstance = &IPHelper{
			configuration: conf,
			idx:           -1,
		}

		safe.SafeGo(func() {
			for {
				helperInstance.getCurrentIP()
				time.Sleep(time.Second * time.Duration(conf.Interval))
			}
		})
	})

	return helperInstance
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

// getIPFromInterface gets IP address from the specific interface.
func (helper *IPHelper) getIPFromInterface() (string, error) {
	ifaces, err := net.InterfaceByName(helper.configuration.IPInterface)
	if err != nil {
		log.Fatal("Can't get network device "+helper.configuration.IPInterface+":", err)
		return "", err
	}

	addrs, err := ifaces.Addrs()
	if err != nil {
		log.Fatal("Can't get address from "+helper.configuration.IPInterface+":", err)
		return "", err
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip == nil {
			continue
		}

		if ip.IsPrivate() {
			continue
		}

		if isIPv4(ip.String()) && strings.ToUpper(helper.configuration.IPType) != utils.IPV4 {
			continue
		} else if strings.ToUpper(helper.configuration.IPType) != utils.IPV6 {
			continue
		}

		if ip.String() != "" {
			log.Printf("Get ip success from network interface by: %s, IP: %s", helper.configuration.IPInterface, ip.String())
			return ip.String(), nil
		}
	}

	return "", errors.New("can't get a valid address from " + helper.configuration.IPInterface)
}

// getCurrentIP gets an IP from either internet or specific interface, depending on configuration.
func (helper *IPHelper) getCurrentIP() {
	var err error
	var ip string
	if helper.configuration.Mikrotik.Enabled {
		if ip = helper.getIPFromMikrotik(); ip == "" {
			log.Fatal("get ip from mikrotik failed. Fallback to get ip from onlinke if possible.")
		} else {
			helper.setCurrentIP(ip)
			return
		}
	}

	if len(helper.reqURLs) > 0 {
		if ip = helper.getIPOnline(); ip == "" {
			log.Fatal("get ip online failed. Fallback to get ip from interface if possible.")
		} else {
			helper.setCurrentIP(ip)
			return
		}
	}

	if helper.configuration.IPInterface != "" {
		if ip, err = helper.getIPFromInterface(); err != nil {
			log.Fatal("get ip from interface failed. There is no more ways to try.")
		} else {
			helper.setCurrentIP(ip)
			return
		}
	}
}

func (helper *IPHelper) setCurrentIP(ip string) {
	helper.mutex.Lock()
	defer helper.mutex.Unlock()

	helper.currentIP = ip
}

func isIPv4(ip string) bool {
	return strings.Count(ip, ":") < 2
}
