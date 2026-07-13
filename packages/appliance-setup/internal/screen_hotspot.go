package internal

import "charm.land/huh/v2"

func (m *Model) WiFiHotspotSetupScreen() *Screen {
	return &Screen{
		Help: "wifihotspot.help",
		Form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title(m.I8nGet("wifihotspot.ssid")).
					CharLimit(32).
					Value(&m.Setup.WifiHotspotSetup.SSID).Validate(m.Validate("wifihotspot.ssid")).
					WithWidth(32),
				huh.NewSelect[string]().
					Title(m.I8nGet("wifihotspot.security")).
					Value(&m.Setup.WifiHotspotSetup.Security).
					Options(
						huh.NewOption("WPA2", "wpa2"),
						huh.NewOption(m.I8nGet("wifihotspot.security.open"), "open"),
					).
					WithWidth(32),
				huh.NewInput().
					Title(m.I8nGet("wifihotspot.password")).
					CharLimit(32).
					Value(&m.Setup.WifiHotspotSetup.Password).Validate(m.Validate("wifihotspot.password")).
					WithWidth(32),
			)),
	}
}
