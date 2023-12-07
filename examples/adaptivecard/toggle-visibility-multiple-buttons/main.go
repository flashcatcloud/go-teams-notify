// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*
This example uses three action "buttons" on the main card body to illustrate
toggling visibility states for multiple elements.

While this example aims to showcase one or more specific features it may not
illustrate overall best practices.

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

	// First text block that we'll use as our header.
	headerTextBlock := adaptivecard.NewTitleTextBlock("Press the buttons to toggle visibility", false)

	// This element is intended to remain visible so we skip setting an ID
	// value. If we did want to change its visibility we would need to set a
	// unique ID value as shown below.
	//
	// headerTextBlock.ID = "headerBlock"

	textBlock1 := adaptivecard.NewHiddenTextBlock("Text Block 1", true)
	textBlock1.ID = "textBlock1"

	textBlock2 := adaptivecard.NewHiddenTextBlock("Text Block 2", true)
	textBlock2.ID = "textBlock2"

	textBlock3 := adaptivecard.NewHiddenTextBlock("Text Block 3", true)
	textBlock3.ID = "textBlock3"

	// This grouping is used for convenience.
	allTextBlocks := []adaptivecard.Element{
		headerTextBlock,
		textBlock1,
		textBlock2,
		textBlock3,
	}

	toggleTargets := []adaptivecard.Element{
		textBlock1,
		textBlock2,
		textBlock3,
	}

	if err := card.AddElement(true, allTextBlocks...); err != nil {
		log.Printf(
			"failed to add text blocks to card: %v",
			err,
		)
		os.Exit(1)
	}

	toggleButton := adaptivecard.NewActionToggleVisibility("Toggle!")
	if err := toggleButton.AddTargetElement(nil, toggleTargets...); err != nil {
		log.Printf(
			"failed to add element IDs to toggle button: %v",
			err,
		)
		os.Exit(1)
	}

	showButton := adaptivecard.NewActionToggleVisibility("Show!")
	if err := showButton.AddVisibleTargetElement(toggleTargets...); err != nil {
		log.Printf(
			"failed to add element IDs to show button: %v",
			err,
		)
		os.Exit(1)
	}

	hideButton := adaptivecard.NewActionToggleVisibility("Hide!")
	if err := hideButton.AddHiddenTargetElement(toggleTargets...); err != nil {
		log.Printf(
			"failed to add element IDs to hide button: %v",
			err,
		)
		os.Exit(1)
	}

	if err := card.AddAction(true, toggleButton); err != nil {
		log.Printf(
			"failed to add toggle button action to card: %v",
			err,
		)
		os.Exit(1)
	}

	if err := card.AddAction(false, showButton); err != nil {
		log.Printf(
			"failed to add show button action to card: %v",
			err,
		)
		os.Exit(1)
	}

	if err := card.AddAction(false, hideButton); err != nil {
		log.Printf(
			"failed to add hide button action to card: %v",
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
