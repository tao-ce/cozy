package internal

import (
	_ "embed"
	"encoding/json"
	"maps"
	"slices"
	"strings"

	"charm.land/huh/v2"
	"github.com/tao-ce/cozy/root/usr/libexec/tao-ce/cozy/appliance-setup/internal/i18n"
)

//go:embed tzdata.json
var tzdata []byte
var timezoneOptions = make(map[string][]string)

func init() {
	err := json.Unmarshal(tzdata, &timezoneOptions)
	if err != nil {
		panic(err)
	}
}

func (m *Model) LocaleScreen() *Screen {
	timezoneRegions := slices.Collect(maps.Keys(timezoneOptions))
	timezoneRegionsOptions := make([]huh.Option[string], len(timezoneRegions))
	for i, region := range timezoneRegions {
		timezoneRegionsOptions[i] = huh.NewOption(region, region)
	}
	langs := i18n.GetLanguagesCodes()
	langsOptions := make([]huh.Option[string], len(langs))
	for i, lang := range langs {
		langsOptions[i] = huh.NewOption(strings.Join([]string{i18n.GetLangName(lang), i18n.Get(lang, "_langName")}, " · "), lang)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				TitleFunc(func() string { return m.I8nGet("locale.language") }, &m.Setup.LocaleSetup.Language).
				Value(&m.Setup.LocaleSetup.Language).
				Inline(true).
				Options(langsOptions...).
				WithWidth(30),
			huh.NewSelect[string]().
				TitleFunc(func() string { return m.I8nGet("locale.timezone.region") }, &m.Setup.LocaleSetup.Language).
				Value(&m.Setup.LocaleSetup.TimezoneRegion).
				Inline(true).
				Options(timezoneRegionsOptions...).
				WithWidth(30),

			huh.NewSelect[string]().
				TitleFunc(func() string { return m.I8nGet("locale.timezone.city") }, &m.Setup.LocaleSetup.Language).
				Value(&m.Setup.LocaleSetup.TimezoneCity).
				Inline(true).
				OptionsFunc(func() []huh.Option[string] {
					options := []huh.Option[string]{}
					for _, city := range timezoneOptions[m.Setup.LocaleSetup.TimezoneRegion] {
						options = append(options, huh.NewOption(strings.ReplaceAll(city, "_", " "), city))
					}
					return options
				}, &m.Setup.LocaleSetup.TimezoneRegion).
				WithWidth(30),
		),
	)

	return &Screen{
		Help: "locale.help",
		Form: form,
	}
}
