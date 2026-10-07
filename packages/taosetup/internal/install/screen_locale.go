package install

import (
	"maps"
	"slices"
	"strings"

	"charm.land/huh/v2"
	"github.com/tao-ce/cozy/taosetup/internal"
	"github.com/tao-ce/cozy/taosetup/internal/i18n"
)

func (m *Model) LocaleScreen() *Screen {
	timezoneRegions := slices.Collect(maps.Keys(internal.GetTimezoneOptions()))
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
					for _, city := range internal.GetTimezoneOptions()[m.Setup.LocaleSetup.TimezoneRegion] {
						label := strings.ReplaceAll(city.Name, "_", " ")
						if city.Country.Code != "" {
							label = strings.Join([]string{label, city.Country.Code}, " · ")
						}
						options = append(options, huh.NewOption(label, city.Name))
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
