package notification

import (
	"bytes"
	"log"
	"text/template"
)

func buildTemplate(currentIP, domain string, tplsrc string) string {
	t := template.New("notification template")
	if _, err := t.Parse(tplsrc); err != nil {
		log.Fatalf("Failed to parse template: %e", err)
		return ""
	}

	var tpl bytes.Buffer
	data := struct {
		CurrentIP string
		Domain    string
	}{
		currentIP,
		domain,
	}
	if err := t.Execute(&tpl, data); err != nil {
		log.Fatal(err)
		return ""
	}

	return tpl.String()
}
