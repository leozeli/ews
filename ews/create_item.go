package ews

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

// https://msdn.microsoft.com/en-us/library/office/aa563009(v=exchg.140).aspx

type CreateItemReq struct {
	XMLName            struct{}          `xml:"m:CreateItem"`
	MessageDisposition string            `xml:"MessageDisposition,attr"`
	SavedItemFolderId  SavedItemFolderId `xml:"m:SavedItemFolderId"`
	Items              Messages          `xml:"m:Items"`
}

type Messages struct {
	Message []Message `xml:"t:Message"`
}

type SavedItemFolderId struct {
	DistinguishedFolderId DistinguishedFolderId `xml:"t:DistinguishedFolderId"`
}

type DistinguishedFolderId struct {
	Id string `xml:"Id,attr"`
}

type Message struct {
	ItemClass    string     `xml:"t:ItemClass"`
	Subject      string     `xml:"t:Subject"`
	Body         Body       `xml:"t:Body"`
	Sender       OneMailbox `xml:"t:Sender"`
	ToRecipients XMailbox   `xml:"t:ToRecipients"`
}

type Body struct {
	BodyType string `xml:"BodyType,attr"`
	Body     []byte `xml:",chardata"`
}

type OneMailbox struct {
	Mailbox Mailbox `xml:"t:Mailbox"`
}

type XMailbox struct {
	Mailbox []Mailbox `xml:"t:Mailbox"`
}

type Mailbox struct {
	EmailAddress string `xml:"t:EmailAddress"`
}

func CreateItem(c *Client, from string, to []string, subject string, body []byte) error {

	cReq := &CreateItemReq{
		MessageDisposition: "SendAndSaveCopy",
		SavedItemFolderId:  SavedItemFolderId{DistinguishedFolderId{Id: "sentitems"}},
	}
	m := &Message{
		ItemClass: "IPM.Note",
		Subject:   subject,
		Body: Body{
			BodyType: "Text",
			Body:     body,
		},
		Sender: OneMailbox{
			Mailbox: Mailbox{
				EmailAddress: from,
			},
		},
	}
	mb := make([]Mailbox, len(to))
	for i, addr := range to {
		mb[i].EmailAddress = addr
	}
	m.ToRecipients.Mailbox = append(m.ToRecipients.Mailbox, mb...)
	cReq.Items.Message = append(cReq.Items.Message, *m)

	reqBytes, err := xml.MarshalIndent(cReq, "", "  ")
	if err != nil {
		return err
	}

	resp, err := c.sendAndReceive(reqBytes)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		bbs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(string(bbs))
	}

	return nil
}