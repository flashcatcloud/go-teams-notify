// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*

This is an example of a client application which uses this library to:

- generate a basic Microsoft Teams message in Adaptive Card format
- submit the message using an explicit proxy server URL

Of note:

- message is in Adaptive Card format
- default timeout
- package-level logging is disabled by default
- validation of known webhook URL prefixes is *enabled*
- simple message submitted to Microsoft Teams consisting of title and
  formatted message body

See https://docs.microsoft.com/en-us/adaptive-cards/authoring-cards/text-features
for the list of supported Adaptive Card text formatting options.

*/

package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
)

func main() {

	// Initialize a new Microsoft Teams client.
	mstClient := goteamsnotify.NewTeamsClient()

	proxyURLString := "http://proxy.example.com:3128"
	proxyUrl, err := url.Parse(proxyURLString)

	if err != nil {
		log.Printf(
			"failed to parse proxy URL %q: %v",
			proxyURLString,
			err,
		)
		os.Exit(1)
	}

	httpClient := mstClient.HTTPClient()
	httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

	// Set webhook url.
	webhookUrl := "https://outlook.office.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

	// The title for message (first TextBlock element).
	msgTitle := "Hello world"

	// Formatted message body.
	msgText := "Here are some examples of formatted stuff like " +
		"\n * this list itself  \n * **bold** \n * *italic* \n * ***bolditalic***"

	// Create message using provided formatted title and text.
	msg, err := adaptivecard.NewSimpleMessage(msgText, msgTitle, true)
	if err != nil {
		log.Printf(
			"failed to create message: %v",
			err,
		)
		os.Exit(1)
	}

	// Send the message with default timeout/retry settings.
	if err := mstClient.Send(webhookUrl, msg); err != nil {
		log.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}
}
