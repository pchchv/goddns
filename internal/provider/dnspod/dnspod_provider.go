package dnspod

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

const providerURL = "https://dnsapi.cn"

type DNSProvider struct {
	configuration *settings.Settings
}

func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}

func (provider *DNSProvider) UpdateIP(domainName, subdomainName, ip string) error {
	domainID := provider.getDomain(domainName)
	if domainID == -1 {
		return errors.New("domain ID not found")
	}

	subdomainID, currentIP := provider.getSubDomain(domainID, subdomainName)
	if subdomainID == "" || currentIP == "" {
		return fmt.Errorf("domain or subdomain not configured yet. domain: %s.%s subDomainID: %s ip: %s", subdomainName, domainName, subdomainID, ip)
	}

	log.Printf("%s.%s Start to update record IP...", subdomainName, domainName)
	return provider.updateIP(domainID, subdomainID, subdomainName, ip)
}

// generateHeader generates the request header for DNSPod API.
func (provider *DNSProvider) generateHeader(content url.Values) url.Values {
	header := url.Values{}
	if provider.configuration.LoginToken != "" {
		header.Add("login_token", provider.configuration.LoginToken)
	}

	header.Add("format", "json")
	header.Add("lang", "en")
	header.Add("error_on_empty", "no")

	for k := range content {
		header.Add(k, content.Get(k))
	}

	return header
}

// postData post data and invoke DNSPod API.
func (provider *DNSProvider) postData(url string, content url.Values) (string, error) {
	client := utils.GetHTTPClient(provider.configuration)

	if client == nil {
		return "", errors.New("failed to create HTTP client")
	}

	values := provider.generateHeader(content)
	req, _ := http.NewRequest("POST", providerURL+url, strings.NewReader(values.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("GoDNS/0.1 (%s)", ""))

	response, err := client.Do(req)
	if err != nil {
		log.Fatal("Post failed:", err)
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Failed to close body:", err)
		}
	}(response.Body)

	resp, _ := io.ReadAll(response.Body)

	return string(resp), nil
}

// updateIP update subdomain with current IP.
func (provider *DNSProvider) updateIP(domainID int64, subDomainID string, subDomainName string, ip string) error {
	value := url.Values{}
	value.Add("domain_id", strconv.FormatInt(domainID, 10))
	value.Add("record_id", subDomainID)
	value.Add("sub_domain", subDomainName)

	if strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		value.Add("record_type", utils.IPTypeA)
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		value.Add("record_type", utils.IPTypeAAAA)
	} else {
		log.Fatal("Must specify ip_type in config for DNSPod.")
		return errors.New("must specify ip_type in config for DNSPod")
	}

	value.Add("record_line", "默认")
	value.Add("value", ip)

	response, err := provider.postData("/Record.Modify", value)
	if err != nil {
		log.Fatal("Failed to update record to new IP:", err)
		return err
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))
	if parseErr != nil {
		log.Fatal(parseErr)
		return err
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		log.Printf("New IP updated: %s", ip)
	} else {
		log.Fatalf("Failed to update IP record: %s", sjson.Get("status").Get("message").MustString())
		return fmt.Errorf("failed to update IP record: %s", sjson.Get("status").Get("message").MustString())
	}

	return nil
}

// getSubDomain returns subdomain by domain id.
func (provider *DNSProvider) getSubDomain(domainID int64, name string) (string, string) {
	var ret, ip string
	value := url.Values{}
	value.Add("domain_id", strconv.FormatInt(domainID, 10))
	value.Add("offset", "0")
	value.Add("length", "1")
	value.Add("sub_domain", name)
	if provider.configuration.IPType == "" || strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		value.Add("record_type", "A")
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		value.Add("record_type", "AAAA")
	} else {
		log.Fatal("Error: must specify \"ip_type\" in config for DNSPod.")
		return "", ""
	}

	response, err := provider.postData("/Record.List", value)
	if err != nil {
		log.Fatal("Failed to get domain list:", err)
		return "", ""
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))
	if parseErr != nil {
		log.Fatal(parseErr)
		return "", ""
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		records, _ := sjson.Get("records").Array()
		for _, d := range records {
			m := d.(map[string]interface{})
			if m["name"] == name {
				ret = m["id"].(string)
				ip = m["value"].(string)
				break
			}
		}

		if len(records) == 0 {
			log.Print("records slice is empty.")
		}
	} else {
		log.Printf("get_subdomain:status code: %s", sjson.Get("status").Get("code").MustString())
	}

	return ret, ip
}

// getDomain returns specific domain by name.
func (provider *DNSProvider) getDomain(name string) (ret int64) {
	values := url.Values{}
	values.Add("type", "all")
	values.Add("offset", "0")
	values.Add("length", "20")
	response, err := provider.postData("/Domain.List", values)
	if err != nil {
		log.Fatal("Failed to get domain list:", err)
		return -1
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))
	if parseErr != nil {
		log.Fatal(parseErr)
		return -1
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		domains, _ := sjson.Get("domains").Array()
		for _, d := range domains {
			m := d.(map[string]interface{})
			if m["name"] == name {
				id := m["id"]
				switch t := id.(type) {
				case json.Number:
					ret, _ = t.Int64()
				}
				break
			}
		}

		if len(domains) == 0 {
			log.Print("domains slice is empty.")
		}
	} else {
		log.Printf("get_domain:status code: %s", sjson.Get("status").Get("code").MustString())
	}

	return
}
