// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*
This is an example of a client application which uses this library to generate
a message with a table within a specific Microsoft Teams channel.

Of note:

- default timeout
- package-level logging is disabled by default
- validation of known webhook URL prefixes is *enabled*
- message is in Adaptive Card format
- text is unformatted
- a table in "grid" format is added to the message

See https://docs.microsoft.com/en-us/adaptive-cards/authoring-cards/text-features
for the list of supported Adaptive Card text formatting options.
*/
package main

import (
	"fmt"
	"log"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
)

func main() {

	// Initialize a new Microsoft Teams client.
	mstClient := goteamsnotify.NewTeamsClient()

	// Set webhook url.
	webhookUrl := "https://outlook.office.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

	vals := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	items := make([]interface{}, len(vals))

	// Create a slice of interface values from our input values.
	for i := range vals {
		items[i] = vals[i]
	}

	tableCells, err := adaptivecard.NewTableCellsWithTextBlock(items)
	if err != nil {
		log.Printf(
			"failed to create table cells: %v",
			err,
		)
		os.Exit(1)
	}

	card := adaptivecard.NewCard()

	// Attach several tables to our card, each with a separator, each showing
	// a custom number of cells per row.
	cellsPerRowExamples := []int{2, 4, 6, 8, 10}
	for _, perRow := range cellsPerRowExamples {
		table, err := adaptivecard.NewTableWithGridFromTableCells(tableCells, perRow)
		if err != nil {
			log.Printf(
				"failed to create table: %v",
				err,
			)
			os.Exit(1)
		}

		table.Separator = true

		card.Body = append(card.Body, table)
	}

	msg := &adaptivecard.Message{
		Type: adaptivecard.TypeMessage,
	}

	msg.Attach(card)

	if err := msg.Prepare(); err != nil {
		log.Printf(
			"failed to prepare message payload: %v",
			err,
		)
		os.Exit(1)
	}

	fmt.Println(msg.PrettyPrint())

	// Having this here makes it easy to comment out the mstClient.Send block.
	_ = mstClient
	_ = webhookUrl

	// Send the message with default timeout/retry settings.
	if err := mstClient.Send(webhookUrl, msg); err != nil {
		log.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}
}
