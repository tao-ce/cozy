package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetMedium() string {
	mediumBytes, err := os.ReadFile(MediumPath)
	medium := DefaultMedium
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not read medium:", err)
	} else {
		medium = strings.TrimSpace(string(mediumBytes))
	}
	return medium
}

func GetWelcomeTimeout() time.Duration {
	envWelcomeTimeout := os.Getenv(EnvWelcomeTimeout)
	if envWelcomeTimeout != "" {
		welcomeTimeout, err := strconv.Atoi(envWelcomeTimeout)
		if err != nil {
			fmt.Fprintln(os.Stderr, "could not parse welcome timeout:", err)
			os.Exit(1)
		} else {
			return time.Duration(welcomeTimeout) * time.Second
		}
	}
	return DefaultWelcomeTimeout
}

func GetConfigPath() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return ""
}
