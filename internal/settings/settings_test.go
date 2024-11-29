package settings

import "testing"

func TestLoadJSONSetting(t *testing.T) {
	var settings Settings
	if err := LoadSettings("../../configs/config_sample.json", &settings); err != nil {
		t.Fatal(err.Error())
	}

	if len(settings.IPUrls) == 0 && settings.IPUrl == "" {
		t.Fatal("neither ip_urls nor ip_url contain valid entries")
	}

	if err := LoadSettings("./file/does/not/exists", &settings); err == nil {
		t.Fatal("file doesn't exist, should return error")
	}
}
