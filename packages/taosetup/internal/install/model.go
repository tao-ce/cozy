package install

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/tao-ce/cozy/taosetup/internal"
	"github.com/tao-ce/cozy/taosetup/internal/i18n"
)

type SetupStep int
type ModelPredicate int

type Model struct {
	Setup              *internal.ApplianceSetup
	SetupStep          SetupStep
	screens            map[SetupStep]*Screen
	imageInput         *huh.Input
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

	m.SetupStep = SetupStepWelcome
	m.UseDefaultSettings = false
	m.Restart = false
	m.Timer = timer.New(m.WelcomeTimeout, timer.WithInterval(time.Millisecond*50))
	m.WindowWidth = 0
	m.WindowHeight = 0
	m.Setup.ApplianceFeatures.StartWiFiHotspot = m.Is(IsWiFiHotspotAvailable)
	m.Setup.TAOCEConfig.Flavor = "full"
	if m.Is(IsRaspberryPi) {
		m.Setup.TAOCEConfig.Flavor = "essential"
	}
	m.Setup.TAOCEConfig.Image = m.DefaultImage()

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

func (m *Model) DefaultImage() string {
	return fmt.Sprintf(DefaultImageFormat, m.Setup.TAOCEConfig.Flavor)
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
