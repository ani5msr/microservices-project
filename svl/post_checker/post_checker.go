package main

import (
	"encoding/json"
	"errors"
	"fmt"

	om "github.com/ani5msr/microservices-project/pkg/object_model"
	"github.com/ani5msr/microservices-project/pkg/post_checker"
	"github.com/ani5msr/microservices-project/pkg/post_checker_events"
	"github.com/nuclio/nuclio-sdk-go"
)

const natsUrl = "nats-cluster.default.svc.cluster.local:4222"

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	r := nuclio.Response{
		StatusCode:  200,
		ContentType: "application/text",
	}

	body := event.GetBody()
	var e om.CheckPostRequest
	err := json.Unmarshal(body, &e)
	if err != nil {
		msg := fmt.Sprintf("failed to unmarshal body: %v", body)
		context.Logger.Error(msg)

		r.StatusCode = 400
		r.Body = []byte(fmt.Sprintf(msg))
		return r, errors.New(msg)

	}

	username := e.Username
	url := e.Url
	if username == "" || url == "" {
		msg := fmt.Sprintf("missing USERNAME ('%s') and/or URL ('%s')", username, url)
		context.Logger.Error(msg)

		r.StatusCode = 400
		r.Body = []byte(msg)
		return r, errors.New(msg)
	}

	status := om.PostStatusValid
	err = post_checker.CheckPost(url)
	if err != nil {
		status = om.PostStatusInvalid
	}

	sender, err := post_checker_events.NewEventSender(natsUrl)
	if err != nil {
		context.Logger.Error(err.Error())

		r.StatusCode = 500
		r.Body = []byte(err.Error())
		return r, err
	}

	sender.OnPostChecked(username, url, status)
	return r, nil
}

func main() {

}
