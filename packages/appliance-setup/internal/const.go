package internal

import "time"

const (
	ApplianceConfigVersion = "1.2.0"
	ApplianceConfigSource  = "appliance-setup"
	DefaultWelcomeTimeout  = 30 * time.Second
	DefaultConfigPath      = "/etc/tao-ce/cozy/appliance-setup.json"
	DefaultMedium          = MediumUnknown
	MediumPath             = "/etc/tao-ce/cozy/medium"
	DefaultImageFormat     = "quay.io/tao-ce/tao-ce:%s"
)
const (
	IsVM ModelPredicate = iota
	IsContainer
	IsWiFiHotspotAvailable
	IsRaspberryPi4
	IsRaspberryPi5
	IsRaspberryPi
	IsUnknown
)
const (
	SetupStepWelcome SetupStep = iota
	SetupStepLocale
	SetupStepFeatures
	SetupStepTAOCEConfig
	SetupStepWiFiHotspotSetup
	SetupStepDone
)
const (
	MediumVM           = "vm"
	MediumContainer    = "container"
	MediumRaspberryPi4 = "rpi4"
	MediumRaspberryPi5 = "rpi5"
	MediumUnknown      = "unknown"
)

const (
	EnvWelcomeTimeout = "APPLIANCE_WELCOME_TIMEOUT_SECONDS"
)
