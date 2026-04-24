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
- a small table is added to the message

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

	card := adaptivecard.Card{
		Type:    adaptivecard.TypeAdaptiveCard,
		Schema:  adaptivecard.AdaptiveCardSchema,
		Version: fmt.Sprintf(adaptivecard.AdaptiveCardVersionTmpl, adaptivecard.AdaptiveCardMaxVersion),
		Body: []adaptivecard.Element{
			{
				Type:              adaptivecard.TypeElementTable,
				GridStyle:         adaptivecard.ContainerStyleAccent,
				ShowGridLines:     func() *bool { show := false; return &show }(),
				FirstRowAsHeaders: func() *bool { show := false; return &show }(),
				Columns: []adaptivecard.Column{
					{
						Type:                           adaptivecard.TypeTableColumnDefinition,
						Width:                          1,
						HorizontalCellContentAlignment: adaptivecard.HorizontalAlignmentLeft,
						VerticalCellContentAlignment:   adaptivecard.VerticalAlignmentBottom,
					},
					{
						Type:                           adaptivecard.TypeTableColumnDefinition,
						Width:                          1,
						HorizontalCellContentAlignment: adaptivecard.HorizontalAlignmentCenter,
						VerticalCellContentAlignment:   adaptivecard.VerticalAlignmentCenter,
					},
					{
						Type:                           adaptivecard.TypeTableColumnDefinition,
						Width:                          1,
						HorizontalCellContentAlignment: adaptivecard.HorizontalAlignmentRight,
						VerticalCellContentAlignment:   adaptivecard.VerticalAlignmentBottom,
					},
				},
				Rows: []adaptivecard.TableRow{
					{
						Type: adaptivecard.TypeTableRow,
						Cells: []adaptivecard.TableCell{
							{
								Type: adaptivecard.TypeTableCell,
								Items: []*adaptivecard.Element{
									{
										Type: adaptivecard.TypeElementTextBlock,
										Wrap: true,
										Text: "Column 1 header",
									},
								},
							},
							{
								Type: adaptivecard.TypeTableCell,
								Items: []*adaptivecard.Element{
									{
										Type: adaptivecard.TypeElementTextBlock,
										Wrap: true,
										Text: "Column 2 header",
									},
								},
							},
							{
								Type: adaptivecard.TypeTableCell,
								Items: []*adaptivecard.Element{
									{
										Type: adaptivecard.TypeElementTextBlock,
										Wrap: true,
										Text: "Column 3 header",
									},
								},
							},
						},
					},
					{
						Type: adaptivecard.TypeTableRow,
						Cells: []adaptivecard.TableCell{
							{
								Type: adaptivecard.TypeTableCell,
								Items: []*adaptivecard.Element{
									{
										Type: adaptivecard.TypeElementTextBlock,
										Wrap: true,
										Text: "Table cell test!",
									},
								},
							},
							{
								Type: adaptivecard.TypeTableCell,
								Items: []*adaptivecard.Element{
									{
										Type: adaptivecard.TypeElementTextBlock,
										Wrap: true,
										Text: "Table cell test!",
									},
								},
							},
							{
								Type: adaptivecard.TypeTableCell,
								Items: []*adaptivecard.Element{
									{
										Type: adaptivecard.TypeElementTextBlock,
										Wrap: true,
										Text: "Table cell test!",
									},
								},
							},
						},
					},
					{
						Type: adaptivecard.TypeTableRow,
						Cells: []adaptivecard.TableCell{
							{
								Type: adaptivecard.TypeTableCell,
								Items: []*adaptivecard.Element{
									{
										Type: adaptivecard.TypeElementTextBlock,
										Wrap: true,
										Text: "Table cell test!",
									},
								},
							},
							{
								Type: adaptivecard.TypeTableCell,
								Items: []*adaptivecard.Element{
									{
										Type: adaptivecard.TypeElementTextBlock,
										Wrap: true,
										Text: "Table cell test!",
									},
								},
							},
							{
								Type: adaptivecard.TypeTableCell,
								Items: []*adaptivecard.Element{
									{
										Type: adaptivecard.TypeElementTextBlock,
										Wrap: true,
										Text: "Table cell test!",
									},
								},
							},
						},
					},
				},
			},
		},
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
