package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	ini "gopkg.in/ini.v1"

	"github.com/outlawlabs/awsctl/pkg/logger"
)

const (
	credentialsFile = "~/.aws/credentials"
	configFile      = "~/.aws/config"

	keyAccessKeyID     = "aws_access_key_id"
	keySecretAccessKey = "aws_secret_access_key"
	keySessionToken    = "aws_session_token"
	keyMFASerial       = "mfa_serial"
	keyRegion          = "region"

	awsCLIHelp = "https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html"
)

var (
	version    = "DEV"
	timestamp  = ""
	commitHash = ""
)

func main() {

	// Set the global ini constants to not require pretty formatting.
	ini.PrettyFormat = false
	ini.PrettyEqual = true

	credentialsFile, err := homedir.Expand(credentialsFile)
	if err != nil {
		logger.Critical("Failed to expand ~/.aws/credentials file: %s.", err)
		os.Exit(1)
	}

	configFile, err := homedir.Expand(configFile)
	if err != nil {
		logger.Critical("Failed to expand ~/.aws/config file: %s.", err)
		os.Exit(1)
	}

	app := kingpin.New("awsctl", "CLI tool to help manage multiple AWS profiles with MFA enabled.").
		Author("github.com/outlawlabs").
		Version(version)

	// TODO: add timestamp checks for auth so no need for redundant session creation.

	configureAuthCommand(app, configFile, credentialsFile)
	configureListCommand(app, configFile, credentialsFile)
	configureNewCommand(app, configFile, credentialsFile)
	configureRemoveCommand(app, configFile, credentialsFile)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

// askForConfirmation asks the user for confirmation. This will not return until
// there is a valid response from the user.
func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		logger.Ask(fmt.Sprintf("%s [y/n]: ", s))

		response, err := reader.ReadString('\n')
		if err != nil {
			logger.Critical(err.Error())
			os.Exit(1)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
