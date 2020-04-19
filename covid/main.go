package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	fbGraphAPIEndpoint = "https://graph.facebook.com/v6.0"
	statsAPIEndpoint   = "https://api.covid19api.com"
)

var (
	fbVerifyToken = os.Getenv("FB_VERIFY_TOKEN")
	fbAccessToken = os.Getenv("FB_ACCESS_TOKEN")
)

type IncomingPayload struct {
	Object  string `json:"object"`
	Entries []struct {
		Messaging []IncomingMessagingPayload `json:"messaging"`
	} `json:"entry"`
}

type IncomingMessagingPayload struct {
	Sender   IncomingMessageUserPayload `json:"sender"`
	Message  *IncomingMessagePayload    `json:"message"`
	Postback *IncomingPostbackPayload   `json:"postback"`
}

type IncomingMessagePayload struct {
	Text       string `json:"text"`
	QuickReply *struct {
		Payload string `json:"payload"`
	} `json:"quick_reply"`
}

type IncomingPostbackPayload struct {
	Title   string `json:"title"`
	Payload string `json:"payload"`
}

type IncomingMessageUserPayload struct {
	ID string `json:"id"`
}

type OutgoingPayload struct {
	Recipient OutgoingMessageUserPayload `json:"recipient"`
}

type OutgoingMessagePayload struct {
	OutgoingPayload

	Message struct {
		Text string `json:"text"`
	} `json:"message"`
}

type OutgoingPostbackPayload struct {
	OutgoingPayload

	Message struct {
		Attachment struct {
			Type    string `json:"type"`
			Payload struct {
				TemplateType string                          `json:"template_type"`
				Text         string                          `json:"text"`
				Buttons      []OutgoingPostbackButtonPayload `json:"buttons"`
			} `json:"payload"`
		} `json:"attachment"`
	} `json:"message"`
}

type OutgoingPostbackButtonPayload struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Payload string `json:"payload"`
}

type OutgoingQuickRepliesPayload struct {
	OutgoingPayload

	MessagingType string `json:"message_type"`
	Message       struct {
		Text         string                                  `json:"text"`
		QuickReplies []OutgoingQuickRepliesQuickReplyPayload `json:"quick_replies"`
	} `json:"message"`
}

type OutgoingQuickRepliesQuickReplyPayload struct {
	Type     string `json:"content_type"`
	Title    string `json:"title"`
	Payload  string `json:"payload"`
	ImageURL string `json:"image_url"`
}

type OutgoingMessageUserPayload struct {
	ID string `json:"id"`
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handleVerification(w, r)
		return
	}

	if r.Method == http.MethodPost {
		handleEvent(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handleVerification(w http.ResponseWriter, r *http.Request) {
	var (
		q = r.URL.Query()

		challenge   = q.Get("hub.challenge")
		mode        = q.Get("hub.mode")
		verifyToken = q.Get("hub.verify_token")
	)

	if mode != "subscribe" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if verifyToken != fbVerifyToken {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	fmt.Fprint(w, challenge)
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	var p IncomingPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if p.Object != "page" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, e := range p.Entries {
		var err error
		switch {
		case len(e.Messaging) != 0:
			switch m := e.Messaging[0]; {
			case m.Message != nil:
				if m.Message.QuickReply != nil {
					err = handleQuickReplyEvent(m.Sender.ID, *m.Message)
					break
				}
				err = handleMessageEvent(m.Sender.ID, *m.Message)
			case m.Postback != nil:
				err = handlePostbackEvent(m.Sender.ID, *m.Postback)
			}
		default:
			err = fmt.Errorf("unknown event type")
		}

		if err != nil {
			log.Print(err)
			return
		}
	}
}

func handleMessageEvent(senderID string, p IncomingMessagePayload) error {
	log.Printf("%s: %s\n", senderID, p.Text)
	return sendButtonsMessage(senderID, p.Text)
}

func handlePostbackEvent(senderID string, p IncomingPostbackPayload) error {
	log.Printf("%s: %s [%s]\n", senderID, p.Title, p.Payload)

	switch p.Payload {
	case "yes":
		return sendQuickRepliesMessage(senderID)
	case "no":
		return sendTextMessage(senderID, "Ok, maybe next time then =)")
	}

	return nil
}

func handleQuickReplyEvent(senderID string, p IncomingMessagePayload) error {
	log.Printf("%s: %s [%s]\n", senderID, p.Text, p.QuickReply.Payload)
	n, err := getStats(p.QuickReply.Payload)
	if err != nil {
		return err
	}
	return sendTextMessage(senderID, fmt.Sprintf("%s has %d confirmed cases as of this moment", p.Text, n))
}

func sendTextMessage(senderID, text string) error {
	var p OutgoingMessagePayload
	p.Recipient.ID = senderID
	p.Message.Text = text
	return sendMessage(p)
}

func sendButtonsMessage(senderID, text string) error {
	user, err := getUserProfile(senderID)
	if err != nil {
		return err
	}

	var p OutgoingPostbackPayload
	p.Recipient.ID = senderID
	p.Message.Attachment.Type = "template"
	p.Message.Attachment.Payload.TemplateType = "button"
	p.Message.Attachment.Payload.Text = fmt.Sprintf("Hi %s, do you want to check for COVID19 stats?", user.FirstName)
	p.Message.Attachment.Payload.Buttons = []OutgoingPostbackButtonPayload{
		{"postback", "Yes", "yes"},
		{"postback", "No", "no"},
	}
	return sendMessage(p)
}

func sendQuickRepliesMessage(senderID string) error {
	var p OutgoingQuickRepliesPayload
	p.Recipient.ID = senderID
	p.MessagingType = "RESPONSE"
	p.Message.Text = "Which country should I check?"
	p.Message.QuickReplies = []OutgoingQuickRepliesQuickReplyPayload{
		{"text", "Worldwide", "WWW",
			"https://i-meet.com/wp-content/uploads/2018/03/cropped-ModernXP-73-Globe-icon-1-1.png"},
		{"text", "Egypt", "EG",
			"https://www.countryflags.io/eg/flat/64.png"},
		{"text", "Denmark", "DK",
			"https://www.countryflags.io/dk/flat/64.png"},
		{"text", "Germany", "DE",
			"https://www.countryflags.io/de/flat/64.png"},
	}
	return sendMessage(p)
}

func sendMessage(p interface{}) error {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&p); err != nil {
		return err
	}

	res, err := http.Post(
		fbGraphAPIEndpoint+"/me/messages?access_token="+fbAccessToken,
		"application/json",
		&b,
	)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		return fmt.Errorf("got: %d %s", res.StatusCode, string(resBody))
	}
	return nil
}

type UserProfile struct {
	FirstName string `json:"first_name"`
}

func getUserProfile(id string) (UserProfile, error) {
	res, err := http.Get(fbGraphAPIEndpoint + "/" + id + "?fields=first_name&access_token=" + fbAccessToken)
	if err != nil {
		return UserProfile{}, err
	}

	var resBody UserProfile
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return UserProfile{}, err
	}
	defer res.Body.Close()

	return resBody, nil
}

type Stats struct {
	Global struct {
		TotalConfirmed int `json:"TotalConfirmed"`
	} `json:"Global"`
	Countries []struct {
		CountryCode    string `json:"CountryCode"`
		TotalConfirmed int    `json:"TotalConfirmed"`
	} `json:"Countries"`
}

func getStats(countryCode string) (int, error) {
	res, err := http.Get(statsAPIEndpoint + "/summary")
	if err != nil {
		return 0, err
	}

	var resBody Stats
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return 0, err
	}
	defer res.Body.Close()

	switch countryCode {
	case "WWW":
		return resBody.Global.TotalConfirmed, nil
	default:
		for _, c := range resBody.Countries {
			if c.CountryCode != countryCode {
				continue
			}
			return c.TotalConfirmed, nil
		}
		return 0, fmt.Errorf("unknown country: %s", countryCode)
	}
	return 0, nil
}

func main() {
	http.HandleFunc("/webhook", handle)
	log.Print("Listening on 3000...")
	http.ListenAndServe("localhost:3000", nil)
}
