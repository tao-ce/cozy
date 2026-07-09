package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"time"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

//go:embed tzdata.json
var tzdata []byte

var timezoneOptions = make(map[string][]string)

type TAOCEConfig struct {
	FQDN   string `json:"fqdn,omitempty"`
	Flavor string `json:"flavor,omitempty"`
}

type ApplianceFeatures struct {
	ShouldStartDDNS  bool `json:"ddns"`
	StartWiFiHotspot bool `json:"hotspot"`
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

type setupStep int

const welcomeTimeout = 10 * time.Second

const (
	SetupStepWelcome setupStep = iota
	SetupStepLocale
	SetupStepFeatures
	SetupStepTAOCEConfig
	SetupStepWiFiHotspotSetup
	SetupStepDone
)

type modelPredicate int

const (
	IsVM modelPredicate = iota
	IsContainer
	IsWiFiHotspotAvailable
	IsRaspberryPi4
	IsRaspberryPi5
	IsRaspberryPi
	IsUnknown
)

var SetupStepNames = map[setupStep]string{
	SetupStepWelcome:          "Welcome",
	SetupStepFeatures:         "Features",
	SetupStepTAOCEConfig:      "TAO CE Config",
	SetupStepWiFiHotspotSetup: "WiFi Hotspot Setup",
	SetupStepDone:             "Done",
}

var SetupStepHelps = map[setupStep]string{
	SetupStepWelcome: `# TAO CE Appliance Setup

Follow this guide to configure your appliance.

`,
	SetupStepLocale: `# Locale Setup

* **Language**: The language of the appliance.
* **Timezone Region**: The timezone region of the appliance.
* **Timezone City**: The timezone city of the appliance.
`,
	SetupStepFeatures: `# Appliance Features

* **Wifi Hotspot**: When running on compatible hardware, a WiFi hotspot will allow to use the appliance as a WiFi access point.
* **DDNS**: When enabled, the appliance will use a DDNS service to broadcast the IP address of the appliance.
`,
	SetupStepTAOCEConfig: `# TAO CE Configuration

* **FQDN**: The fully qualified domain name of the appliance.
* **Flavor**: The flavor of the appliance.
`,
	SetupStepWiFiHotspotSetup: `# WiFi Hotspot Setup
* **SSID**: The SSID of the WiFi hotspot.
* **Security**: The security of the WiFi hotspot.
* **Password**: The password of the WiFi hotspot.
`,
	SetupStepDone: `# Setup complete
Setup complete
`,
}

var leftPaneStyle = lipgloss.NewStyle().Width(30).Height(22).Border(lipgloss.NormalBorder())
var rightPaneStyle = lipgloss.NewStyle().Width(50).Height(22).Border(lipgloss.NormalBorder())
var helpPaneStyle = lipgloss.NewStyle().Width(80).Height(2).Align(lipgloss.Left)
var helpRenderer, _ = glamour.NewTermRenderer(
	glamour.WithWordWrap(30),
	glamour.WithStandardStyle("dark"),
)

type model struct {
	setup              ApplianceSetup
	setupStep          setupStep
	forms              map[setupStep]*huh.Form
	useDefaultSettings bool
	timer              timer.Model
	configPath         string
	windowWidth        int
	windowHeight       int
}

func (m *model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m *model) toJSON() string {
	m.setup.Metadata = map[string]string{
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
		"source":    "appliance-setup",
	}

	if !m.setup.ApplianceFeatures.StartWiFiHotspot {
		m.setup.WifiHotspotSetup = WiFiHotspotSetup{}
	}

	jsonData, err := json.Marshal(&m.setup)
	if err != nil {
		fmt.Println("Error marshalling setup:", err)
	}
	return string(jsonData)
}

func (m *model) is(predicate modelPredicate) bool {
	switch predicate {
	case IsUnknown:
		return m.setup.Medium == "unknown"
	case IsRaspberryPi4:
		return m.setup.Medium == "rpi4"
	case IsRaspberryPi5:
		return m.setup.Medium == "rpi5"
	case IsRaspberryPi:
		return m.is(IsRaspberryPi4) || m.is(IsRaspberryPi5)
	case IsVM:
		return m.setup.Medium == "vm"
	case IsContainer:
		return m.setup.Medium == "container"
	case IsWiFiHotspotAvailable:
		return !m.is(IsVM) && !m.is(IsContainer) && !m.is(IsUnknown)
	}
	return false
}

func NewModel(medium string) *model {
	m := model{
		setupStep:          SetupStepWelcome,
		useDefaultSettings: false,
		timer:              timer.New(welcomeTimeout, timer.WithInterval(time.Millisecond*10)),
		setup: ApplianceSetup{
			Medium: medium,
			LocaleSetup: LocaleSetup{
				Language:       "en_US",
				TimezoneRegion: "Etc",
				TimezoneCity:   "UTC",
			},
			ApplianceFeatures: ApplianceFeatures{
				ShouldStartDDNS:  true,
				StartWiFiHotspot: true,
			},
			TAOCEConfig: TAOCEConfig{
				FQDN:   "tao-community-edition.local",
				Flavor: "full",
			},
			WifiHotspotSetup: WiFiHotspotSetup{
				SSID:     "TAO Community Edition",
				Security: "wpa2",
				Password: "ChangeMeNow",
			},
		},
	}

	m.setup.ApplianceFeatures.StartWiFiHotspot = m.is(IsWiFiHotspotAvailable)

	m.forms = make(map[setupStep]*huh.Form)
	m.forms[SetupStepWelcome] = m.welcomeScreen()
	m.forms[SetupStepLocale] = m.localeScreen()
	m.forms[SetupStepFeatures] = m.featuresScreen()
	m.forms[SetupStepTAOCEConfig] = m.taoCeConfigScreen()
	m.forms[SetupStepWiFiHotspotSetup] = m.wifiHotspotSetupScreen()
	m.forms[SetupStepDone] = m.doneScreen()

	return &m
}

func (m *model) appView() tea.View {
	help, _ := helpRenderer.Render(SetupStepHelps[m.setupStep])
	timeout := float64(m.timer.Timeout.Seconds()) / float64(welcomeTimeout.Seconds())

	var progressBar, welcomeTitle string
	if m.setupStep == SetupStepWelcome {
		welcomeTitle = "\n" + lipgloss.NewStyle().
			Height(2).
			Align(lipgloss.Center).
			Width(30).
			Padding(0, 6).
			MarginLeft(10).
			Bold(true).
			Italic(true).
			Background(lipgloss.BrightRed).
			Render("Welcome to TAO CE Appliance Setup")
		progressBar = strings.Join([]string{
			"",
			"",
			"Appliance Setup will proceed with default settings in " + time.Duration(m.timer.Timeout).Round(time.Second).String(),
			"",
			progress.New(
				progress.WithoutPercentage(),
				progress.WithWidth(30),
				progress.WithScaled(true),
				progress.WithColors(lipgloss.BrightYellow),
			).ViewAs(timeout),
			"",
		}, "\n")
	} else {
		progressBar = ""
		welcomeTitle = ""
	}

	return tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				leftPaneStyle.Render(
					lipgloss.JoinVertical(
						lipgloss.Left,
						help,
					),
				),
				lipgloss.JoinVertical(
					lipgloss.Center,
					rightPaneStyle.Render(
						lipgloss.JoinVertical(
							lipgloss.Left,
							welcomeTitle,
							lipgloss.NewStyle().Align(lipgloss.Center).Width(30).MarginLeft(10).Render(progressBar),
							"",
							m.forms[m.setupStep].View(),
							"",
						),
					),
				),
			),
			helpPaneStyle.Render(m.forms[m.setupStep].Help().ShortHelpView(m.forms[m.setupStep].KeyBinds())),
		),
	)
}

func (m *model) welcomeScreen() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Affirmative("Default settings").
				Negative("Customize settings").
				Value(&m.useDefaultSettings).
				Title("Do you want to continue with default settings?"),
		),
	).WithShowHelp(false)
}

func (m *model) localeScreen() *huh.Form {
	timezoneRegions := slices.Collect(maps.Keys(timezoneOptions))
	timezoneRegionsOptions := make([]huh.Option[string], len(timezoneRegions))
	for i, region := range timezoneRegions {
		timezoneRegionsOptions[i] = huh.NewOption(region, region)
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Language").
				Value(&m.setup.LocaleSetup.Language).
				Inline(true).
				Options(
					huh.NewOption("English", "en_US.UTF-8"),
					huh.NewOption("Japanese / 日本語", "ja_JP.UTF-8"),
					huh.NewOption("French / Français", "fr_FR.UTF-8"),
					huh.NewOption("German / Deutsch", "de_DE.UTF-8"),
					huh.NewOption("Italian / Italiano", "it_IT.UTF-8"),
					huh.NewOption("Spanish / Español", "es_ES.UTF-8"),
					huh.NewOption("Portuguese / Português", "pt_PT.UTF-8"),
				),
			huh.NewSelect[string]().
				Title("Timezone Region").
				Value(&m.setup.LocaleSetup.TimezoneRegion).
				Inline(true).
				Options(timezoneRegionsOptions...),

			huh.NewSelect[string]().
				Title("Timezone City").
				Value(&m.setup.LocaleSetup.TimezoneCity).
				Inline(true).
				OptionsFunc(func() []huh.Option[string] {
					options := []huh.Option[string]{}
					for _, city := range timezoneOptions[m.setup.LocaleSetup.TimezoneRegion] {
						options = append(options, huh.NewOption(strings.ReplaceAll(city, "_", " "), city))
					}
					return options
				}, &m.setup.LocaleSetup.TimezoneRegion),
		),
	).WithShowHelp(false)
}
func (m *model) doneScreen() *huh.Form {
	configPath := m.configPath
	if configPath == "" {
		configPath = "standard output"
	}
	return huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Description("Setup complete.\n\nConfiguration saved to " + configPath),
		),
	).WithShowHelp(false)
}

func (m *model) featuresScreen() *huh.Form {
	description := "\n"
	if !m.is(IsWiFiHotspotAvailable) {
		description += lipgloss.NewStyle().Foreground(lipgloss.Yellow).Bold(true).Render("Note: WiFi Hotspot is not available on " + m.setup.Medium + " medium.\n")
	}
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("WiFi Hotspot  ").
				Value(&m.setup.ApplianceFeatures.StartWiFiHotspot).
				Inline(true).
				Affirmative("ON").
				Negative("OFF"),
			huh.NewConfirm().
				Title("DDNS          ").
				Value(&m.setup.ApplianceFeatures.ShouldStartDDNS).
				Inline(true).
				Affirmative("ON").
				Negative("OFF"),
		).
			Description(description),
	).WithShowHelp(false)
}

func (m *model) quit() tea.Cmd {
	m.timer.Stop()
	m.setupStep = SetupStepDone

	return tea.Batch(tea.Quit)
}

func (m *model) taoCeConfigScreen() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("TAO CE domain name").
				Value(&m.setup.TAOCEConfig.FQDN),
			huh.NewSelect[string]().
				Title("Flavor").
				Value(&m.setup.TAOCEConfig.Flavor).
				Options(
					huh.NewOption("Full      Complete TAO CE stack", "full"),
					huh.NewOption("Essential Without scoring/proctoring/devkit", "essential"),
					huh.NewOption("Lite      Only portal and delivery", "lite"),
					huh.NewOption("Minimal   Only delivery", "minimal"),
				),
		),
	).WithShowHelp(false)
}

func (m *model) wifiHotspotSetupScreen() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("WiFi hotspot SSID").
				Value(&m.setup.WifiHotspotSetup.SSID),
			huh.NewSelect[string]().
				Title("WiFi hotspot security").
				Value(&m.setup.WifiHotspotSetup.Security).
				Options(
					huh.NewOption("WPA2", "wpa2"),
					huh.NewOption("Open", "open"),
				),
			huh.NewInput().
				Title("WiFi hotspot password").
				Value(&m.setup.WifiHotspotSetup.Password),
		)).WithShowHelp(false)

}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		return m, nil
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case timer.TimeoutMsg:
		if m.setupStep == SetupStepWelcome {
			return m, m.quit()
		}

	case tea.KeyPressMsg:
		switch msg.String() {
		case "pgup":
			if m.setupStep > SetupStepWelcome {
				m.setupStep--
			}
		case "pgdown":
			if m.setupStep < SetupStepDone {
				m.setupStep++
			}
		case "ctrl+c", "esc":
			return m, m.quit()
		}
	}

	var cmds []tea.Cmd

	if m.forms[m.setupStep] != nil {
		form, cmd := m.forms[m.setupStep].Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.forms[m.setupStep] = f
			cmds = append(cmds, cmd)
		}
	}

	if m.forms[m.setupStep].State == huh.StateCompleted {

		if m.useDefaultSettings {
			return m, m.quit()
		}

		m.setupStep++

		if m.setupStep == SetupStepWiFiHotspotSetup && !m.setup.ApplianceFeatures.StartWiFiHotspot {
			m.setupStep++
		}

		if m.setupStep > SetupStepDone {
			return m, m.quit()
		}

		cmds = append(cmds, m.forms[m.setupStep].Init())
	}
	return m, tea.Batch(cmds...)
}

func (m *model) View() tea.View {
	if m.forms[m.setupStep] != nil {
		v := m.appView()
		v.AltScreen = true
		return v
	}

	return tea.NewView("")
}

func main() {
	err := json.Unmarshal(tzdata, &timezoneOptions)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not unmarshal timezone options:", err)
		os.Exit(1)
	}

	mediumBytes, err := os.ReadFile("/etc/tao-ce/cozy/medium")
	medium := "unknown"
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not read medium:", err)
	} else {
		medium = strings.TrimSpace(string(mediumBytes))
	}

	model := NewModel(medium)

	if len(os.Args) > 1 {
		model.configPath = os.Args[1]
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "could not start program:", err)
		os.Exit(1)
	}

	if model.configPath != "" {
		os.WriteFile(model.configPath, []byte(model.toJSON()+"\n"), 0644)
	} else {
		fmt.Fprintln(os.Stdout, model.toJSON())
	}

}
