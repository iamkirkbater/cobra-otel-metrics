package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// global vars
var (
	// Globally accessible var to determine if the user has opted in or not
	// for metric collection.
	UserHasOptedInForMetrics bool

	// TODO - test this works on linux and when $HOME is not set or the path in XDG_CONFIG_HOME is not relative
	defaultOptInDirectory        = func() string { dir, _ := os.UserConfigDir(); return dir }()
	defaultOptInFilenamePostfix  = "metrics-optin"
	defaultOptOutFilenamePostfix = "metrics-optout"
)

var defaultConsentPrompt = `
#######################################################
#                                                     # 
# This command line utility would like to collect     #
# anonymous metrics on usage patterns to help build   #
# a better tool.                                      #
#                                                     #
# We require users to explicitly opt-in to consent to #
# allow us to collect these metrics.                  #
#                                                     #
#######################################################
`

var defaultConsentQuestion = `
Would you like to share anonymous usage stats with
the developers of this tool? (y|N) `

var defaultOptInMessage = "\nThank you for sharing!"
var defaultOptOutMessage = "\nYou have opted out. We will not collect metrics."
var defaultConsentRetryMessage = "\nInvalid value detected. Only Y and N are allowed..."

func HandleMetricsOptIn(cmd *cobra.Command) error {
	// By default, if we have a non-interactive session let's automatically
	// collect metrics.
	rootCmdName := GetRootCmdName(cmd)
	if !IsTTY() {
		handleDefaultInteractiveOptOut(getDefaultOptOutFilePath(rootCmdName))
		return nil
	}

	filePath := getDefaultOptInFilePath(rootCmdName)
	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			printConsentPrompt()
			userConsent, promptErr := askForConsentAndSave(filePath)
			if promptErr != nil {
				UserHasOptedInForMetrics = false
				return fmt.Errorf("Error prompting for metric collection consent. Metrics will not be collected. %w", promptErr)
			}
			UserHasOptedInForMetrics = userConsent
			return nil
		}
		return fmt.Errorf("Unexpected error opening metric opt-in file. Metrics will not be collected. %w", err)
	}

	optIn := make([]byte, 1)
	_, err = file.Read(optIn)
	if err != nil {
		return fmt.Errorf("Unexpected error reading metric opt-in file. Metrics will not be collected. %w", err)
	}
	if string(optIn) == "0" {
		UserHasOptedInForMetrics = false
		return nil
	}
	if string(optIn) == "1" {
		UserHasOptedInForMetrics = true
		return nil
	}
	// If we get here, we can assume that the file is corrupted or has unknown values.
	// So let's ask the user to opt-in again.
	printConsentPrompt()
	userConsent, promptErr := askForConsentAndSave(filePath)
	if promptErr != nil {
		UserHasOptedInForMetrics = false
		return fmt.Errorf("Error prompting for metric collection consent. Metrics will not be collected. %w", promptErr)
	}
	UserHasOptedInForMetrics = userConsent
	return nil
}

func handleDefaultInteractiveOptOut(filePath string) {
	_, err := os.Open(filePath)
	if errors.Is(err, fs.ErrNotExist) {
		UserHasOptedInForMetrics = true
		return
	}
	UserHasOptedInForMetrics = false
}

func printConsentPrompt() {
	fmt.Fprint(os.Stderr, defaultConsentPrompt)
}

func askForConsentAndSave(filePath string) (bool, error) {
	fmt.Fprint(os.Stderr, defaultConsentQuestion)
	reader := bufio.NewReader(os.Stdin)
	message, _ := reader.ReadString('\n')
	if message == "" {
		message = "n"
	}

	var optInStatus bool
	switch strings.ToLower(message)[0] {
	case 'n':
		optInStatus = false
		fmt.Fprintln(os.Stderr, defaultOptOutMessage)
	case 'y':
		optInStatus = true
		fmt.Fprintln(os.Stderr, defaultOptInMessage)
	default:
		fmt.Fprintln(os.Stderr, defaultConsentRetryMessage)
		return askForConsentAndSave(filePath)
	}

	err := saveOptInStatus(filePath, optInStatus)
	if err != nil {
		return optInStatus, fmt.Errorf("Unable to save opt-in status. User will be asked again for consent for telemetry. %w", err)
	}
	return optInStatus, nil
}

func saveOptInStatus(filePath string, optInStatus bool) error {
	b := "0"
	if optInStatus {
		b = "1"
	}

	err := os.WriteFile(filePath, []byte(b), 0666)
	return err
}

func getDefaultConsentFilePath(cmdName string, postfix string) string {
	optInFilename := "." + cmdName + "-" + postfix
	optInFilepath := filepath.Join(defaultOptInDirectory, optInFilename)
	return optInFilepath
}

func getDefaultOptInFilePath(cmdName string) string {
	return getDefaultConsentFilePath(cmdName, defaultOptInFilenamePostfix)
}

func getDefaultOptOutFilePath(cmdName string) string {
	return getDefaultConsentFilePath(cmdName, defaultOptOutFilenamePostfix)
}
