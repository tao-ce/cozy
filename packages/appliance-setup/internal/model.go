package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
	"github.com/tao-ce/cozy/root/usr/libexec/tao-ce/cozy/appliance-setup/internal/i18n"
)

type SetupStep int
type ModelPredicate int

type TAOCEConfig struct {
	FQDN   string `json:"fqdn,omitempty"`
	Flavor string `json:"flavor,omitempty"`
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
type Model struct {
	Setup              ApplianceSetup
	SetupStep          SetupStep
	screens            map[SetupStep]*Screen
	UseDefaultSettings bool
	Restart            bool
	WelcomeTimeout     time.Duration
	Timer              timer.Model
	ConfigPath         string
	WindowWidth        int
	WindowHeight       int
}

func (m *Model) I8nGet(key string) string {
	return i18n.Get(m.Setup.LocaleSetup.Language, key)
}

func (m *Model) Init() tea.Cmd {
	m.Setup.Medium = GetMedium()
	m.WelcomeTimeout = GetWelcomeTimeout()
	m.ConfigPath = GetConfigPath()

	m.SetupStep = SetupStepWelcome
	m.UseDefaultSettings = false
	m.Restart = false
	m.Timer = timer.New(m.WelcomeTimeout, timer.WithInterval(time.Millisecond*50))
	m.WindowWidth = 0
	m.WindowHeight = 0
	m.Setup.ApplianceFeatures.StartWiFiHotspot = m.Is(IsWiFiHotspotAvailable)

	m.screens = make(map[SetupStep]*Screen)
	return m.Timer.Init()
}

func (m *Model) GetScreen() *Screen {
	var screen *Screen
	if _, ok := m.screens[m.SetupStep]; !ok {
		switch m.SetupStep {
		case SetupStepWelcome:
			screen = m.WelcomeScreen()
		case SetupStepLocale:
			screen = m.LocaleScreen()
		case SetupStepFeatures:
			screen = m.FeaturesScreen()
		case SetupStepTAOCEConfig:
			screen = m.TAOCEConfigScreen()
		case SetupStepWiFiHotspotSetup:
			screen = m.WiFiHotspotSetupScreen()
		case SetupStepDone:
			screen = m.DoneScreen()
		default:
			return nil
		}
		m.screens[m.SetupStep] = screen
	}
	return m.screens[m.SetupStep]
}

func (m *Model) Export() {
	m.Setup.Metadata = map[string]string{
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   ApplianceConfigVersion,
		"source":    ApplianceConfigSource,
	}

	if !m.Setup.ApplianceFeatures.StartWiFiHotspot {
		m.Setup.WifiHotspotSetup = WiFiHotspotSetup{}
	}
	if m.Setup.WifiHotspotSetup.Security == "open" {
		m.Setup.WifiHotspotSetup.Password = ""
	}

	jsonData, err := json.Marshal(&m.Setup)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error marshalling setup:", err)
		return
	}
	if m.ConfigPath != "" {
		os.WriteFile(m.ConfigPath, jsonData, 0644)
	} else {
		fmt.Fprintln(os.Stdout, string(jsonData))
	}
}

func (m *Model) Validate(key string) func(string) error {
	return func(value string) error {
		switch key {
		case "taoce.fqdn":
			if value == "" {
				return errors.New(m.I8nGet("error." + key + ".required"))
			}
			if m.Setup.ApplianceFeatures.StartMulticastDNS && !strings.HasSuffix(value, ".local") {
				return errors.New(m.I8nGet("error." + key + ".not_local"))
			}
			if !regexp.MustCompile(`^[a-z0-9-]+(\.[a-z0-9-]+)*$`).MatchString(value) {
				return errors.New(m.I8nGet("error." + key + ".invalid"))
			}
		case "wifihotspot.ssid":
			if value == "" {
				return errors.New(m.I8nGet("error." + key + ".required"))
			}
			if !regexp.MustCompile(`^[ -~]{1,32}$`).MatchString(value) {
				return errors.New(m.I8nGet("error." + key + ".invalid"))
			}
		case "wifihotspot.security":
			if value == "" {
				return errors.New(m.I8nGet("error." + key + ".required"))
			}
		case "wifihotspot.password":
			if value == "" {
				return errors.New(m.I8nGet("error." + key + ".required"))
			}
			if !regexp.MustCompile(`^[ -~]{1,32}$`).MatchString(value) {
				return errors.New(m.I8nGet("error." + key + ".invalid"))
			}
		}
		return nil
	}
}

func (m *Model) Is(predicate ModelPredicate) bool {
	switch predicate {
	case IsUnknown:
		return m.Setup.Medium == MediumUnknown
	case IsRaspberryPi4:
		return m.Setup.Medium == MediumRaspberryPi4
	case IsRaspberryPi5:
		return m.Setup.Medium == MediumRaspberryPi5
	case IsRaspberryPi:
		return m.Is(IsRaspberryPi4) || m.Is(IsRaspberryPi5)
	case IsVM:
		return m.Setup.Medium == MediumVM
	case IsContainer:
		return m.Setup.Medium == MediumContainer
	case IsWiFiHotspotAvailable:
		return !m.Is(IsVM) && !m.Is(IsContainer) && !m.Is(IsUnknown)
	}
	return false
}

func (m *Model) restart() {
	m.SetupStep = SetupStepLocale
	m.UseDefaultSettings = false
	m.Restart = false
	m.screens = make(map[SetupStep]*Screen)
}

func (m *Model) quit() tea.Cmd {
	m.Timer.Stop()
	m.SetupStep = SetupStepDone

	return tea.Batch(tea.Quit)
}
