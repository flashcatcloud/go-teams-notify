// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*
This example illustrates how to toggle visibility for a text block using a
container's select action.

Of note:

- default timeout
- package-level logging is disabled by default
- validation of known webhook URL prefixes is *enabled*

See https://docs.microsoft.com/en-us/adaptive-cards/authoring-cards/text-features
for the list of supported Adaptive Card text formatting options.
*/
package main

import (
	"log"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
)

func main() {

	// Initialize a new Microsoft Teams client.
	mstClient := goteamsnotify.NewTeamsClient()

	// Set webhook url.
	//
	// NOTE: This is for illustration purposes only. Best practice is to NOT
	// hardcode credentials of any kind.
	webhookUrl := "https://outlook.office.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

	// Allow specifying webhook URL via environment variable, fall-back to
	// hard-coded value in this example file.
	expectedEnvVar := "WEBHOOK_URL"
	envWebhookURL := os.Getenv(expectedEnvVar)
	switch {
	case envWebhookURL != "":
		log.Printf(
			"Using webhook URL %q from environment variable %q\n\n",
			envWebhookURL,
			expectedEnvVar,
		)
		webhookUrl = envWebhookURL
	default:
		log.Println(expectedEnvVar, "environment variable not set.")
		log.Printf("Using hardcoded value %q as fallback\n\n", webhookUrl)
	}

	// Create blank card that we'll manually fill in.
	card := adaptivecard.NewCard()

	headerTextBlock := adaptivecard.NewTitleTextBlock("Press the link text to show details", false)

	// Details that we'll hide by default but allow toggling visibility for.
	detailsMessageBlock := adaptivecard.NewHiddenTextBlock("Details text block content here", true)
	detailsMessageBlock.ID = "details"

	cardBodyElements := []adaptivecard.Element{
		headerTextBlock,
		detailsMessageBlock,
	}

	if err := card.AddElement(true, cardBodyElements...); err != nil {
		log.Printf(
			"failed to add card body text blocks: %v",
			err,
		)
		os.Exit(1)
	}

	showDetailsTextBlock := adaptivecard.NewTextBlock("Show details", false)
	showDetailsTextBlock.ID = "showDetails"

	hideDetailsTextBlock := adaptivecard.NewHiddenTextBlock("Hide details", false)
	hideDetailsTextBlock.ID = "hideDetails"

	showHideLinkContainer := adaptivecard.NewContainer()

	if err := showHideLinkContainer.AddElement(true, showDetailsTextBlock); err != nil {
		log.Printf(
			"failed to add text block to container: %v",
			err,
		)
		os.Exit(1)
	}

	if err := showHideLinkContainer.AddElement(false, hideDetailsTextBlock); err != nil {
		log.Printf(
			"failed to add text block to container: %v",
			err,
		)
		os.Exit(1)
	}

	toggleTargets := []adaptivecard.Element{
		detailsMessageBlock,
		showDetailsTextBlock,
		hideDetailsTextBlock,
	}

	detailsDisplayAction := adaptivecard.NewActionToggleVisibility("")
	if err := detailsDisplayAction.AddTargetElement(nil, toggleTargets...); err != nil {
		log.Printf(
			"failed to add element IDs to toggle action: %v",
			err,
		)
		os.Exit(1)
	}

	if err := showHideLinkContainer.AddSelectAction(detailsDisplayAction); err != nil {
		log.Printf(
			"failed to add action to container: %v",
			err,
		)
		os.Exit(1)
	}

	if err := card.AddContainer(false, showHideLinkContainer); err != nil {
		log.Printf(
			"failed to add button container to card body: %v",
			err,
		)
		os.Exit(1)
	}

	// Create new Message using Card as input.
	msg, err := adaptivecard.NewMessageFromCard(card)
	if err != nil {
		log.Printf("failed to create message from card: %v", err)
		os.Exit(1)
	}

	// We explicitly prepare the message for transmission ahead of calling
	// mstClient.Send so that we can print the JSON payload in human readable
	// format for review. If we do not explicitly prepare the message then the
	// mstClient.Send call will handle that for us (which is how this is
	// usually handled).
	{
		if err := msg.Prepare(); err != nil {
			log.Printf(
				"failed to prepare message: %v",
				err,
			)
			os.Exit(1)
		}

		log.Println(msg.PrettyPrint())
	}

	// Send the message with default timeout/retry settings.
	if err := mstClient.Send(webhookUrl, msg); err != nil {
		log.Printf("failed to send message: %v", err)
		os.Exit(1)
	}

}
