// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*

This is an example of a client application which uses this library to:

- generate a Microsoft Teams message in MessageCard format
- submit the message using an explicit proxy server URL

Of note:

- message is in MessageCard format
- default timeout
- package-level logging is disabled by default
- validation of known webhook URL prefixes is *enabled*
- simple message submitted to Microsoft Teams consisting of title and
  formatted message body

*/

package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
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

	// Setup message card.
	msgCard := messagecard.NewMessageCard()
	msgCard.Title = "Hello world"
	msgCard.Text = "Here are some examples of formatted stuff like " +
		"<br> * this list itself  <br> * **bold** <br> * *italic* <br> * ***bolditalic***"
	msgCard.ThemeColor = "#DF813D"

	// Send the message with default timeout/retry settings.
	if err := mstClient.Send(webhookUrl, msgCard); err != nil {
		log.Printf("failed to send message: %v", err)
		os.Exit(1)
	}
}
