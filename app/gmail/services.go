package gmail

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func readEmails() {
	// Connect to Gmail's IMAP server
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	// Login to Gmail account using your email and app password
	err = c.Login("danieltanimanuel@gmail.com", "your-app-password")
	if err != nil {
		log.Fatal(err)
	}

	// Select the mailbox (INBOX)
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the last 5 messages
	from := uint32(1)
	to := mbox.Messages
	if to > 5 {
		to = 5
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	// Fetch the emails
	messages := make(chan *imap.Message, 5)
	go func() {
		err := c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Print the subjects of the emails
	for msg := range messages {
		fmt.Println("Subject:", msg.Envelope.Subject)
	}
}
