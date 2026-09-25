package internal

import "time"

const (
	ApplianceConfigVersion = "1.2.0"
	ApplianceConfigSource  = "taosetup"
	DefaultWelcomeTimeout  = 30 * time.Second
	DefaultConfigPath      = "/etc/tao-ce/cozy/taosetup.json"
	MediumPath             = "/etc/tao-ce/cozy/medium"
	DefaultImageFormat     = "quay.io/tao-ce/tao-ce:%s"
)
