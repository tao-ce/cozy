package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type TAOCEConfig struct {
	FQDN   string `json:"fqdn,omitempty"`
	Flavor string `json:"flavor,omitempty"`
	Image  string `json:"image,omitempty"`
}

type ApplianceFeatures struct {
	StartMulticastDNS bool `json:"mdns"`
	StartWiFiHotspot  bool `json:"hotspot"`
}

type WiFiHotspotSetup struct {
	SSID     string `json:"ssid,omitempty"`
	Security string `json:"security,omitempty"`
	Password string `json:"password,omitempty"`
}

type LocaleSetup struct {
	Language       string `json:"language,omitempty"`
	TimezoneRegion string `json:"timezoneRegion,omitempty"`
	TimezoneCity   string `json:"timezoneCity,omitempty"`
}

type ApplianceSetup struct {
	ApplianceFeatures ApplianceFeatures `json:"features,omitempty"`
	LocaleSetup       LocaleSetup       `json:"locale,omitempty"`
	TAOCEConfig       TAOCEConfig       `json:"taoCe,omitempty"`
	WifiHotspotSetup  WiFiHotspotSetup  `json:"hotspot,omitempty"`
	Metadata          map[string]string `json:"metadata"`
	Medium            string            `json:"medium,omitempty"`
}

func (c *ApplianceSetup) Export(configPath string) {
	c.Metadata = map[string]string{
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   ApplianceConfigVersion,
		"source":    ApplianceConfigSource,
	}

	if !c.ApplianceFeatures.StartWiFiHotspot {
		c.WifiHotspotSetup = WiFiHotspotSetup{}
	}
	if c.WifiHotspotSetup.Security == "open" {
		c.WifiHotspotSetup.Password = ""
	}

	jsonData, err := json.Marshal(&c)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error marshalling setup:", err)
		return
	}
	if configPath != "" {
		os.WriteFile(configPath, jsonData, 0644)
	} else {
		fmt.Fprintln(os.Stdout, string(jsonData))
	}
}

func (c *ApplianceSetup) Import(configPath string) {
	jsonData, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading setup:", err)
		return
	}
	err = json.Unmarshal(jsonData, &c)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error unmarshalling setup:", err)
		return
	}
}
