package slack

import (
	"context"
	"encoding/json"
	"strings"
)

type UpdateViewTrigger struct {
	TriggerID string      `json:"trigger_id"` //Required. Must respond within 3 seconds.
	View      ViewPayload `json:"view"`       //Required.
	ViewId    string      `json:"view_id"`
}
type ViewTrigger struct {
	TriggerID string      `json:"trigger_id"` //Required. Must respond within 3 seconds.
	View      ViewPayload `json:"view"`       //Required.
}

type ViewPayloadCallback struct {
	ViewPayload
	Id string `json:"id"`
}

type ViewPayload struct {
	Type       string          `json:"type"`
	Title      TextBlockObject `json:"title"`
	Blocks     []Block         `json:"blocks"`
	Close      TextBlockObject `json:"close"`
	Submit     TextBlockObject `json:"submit"`
	CallbackId string          `json:"callback_id"`
}

// DialogCallback DEPRECATED use InteractionCallback
type ViewCallback InteractionCallback

// ViewOpenResponse response from `dialog.open`
type ViewOpenResponse struct {
	SlackResponse
	ViewResponseMetadata ViewResponseMetadata `json:"response_metadata"`
}

// ViewResponseMetadata lists the error messages
type ViewResponseMetadata struct {
	Messages []string `json:"messages"`
}

// OpenView opens a dialog window where the triggerID originated from.
// EXPERIMENTAL: view functionality is currently experimental, api is not considered stable.
func (api *Client) OpenView(triggerID string, viewPayload ViewPayload) (err error) {
	return api.OpenViewContext(context.Background(), triggerID, viewPayload)
}

// OpenViewgContext opens a dialog window where the triggerId originated from with a custom context
// EXPERIMENTAL: view functionality is currently experimental, api is not considered stable.
func (api *Client) OpenViewContext(ctx context.Context, triggerID string, viewPayload ViewPayload) (err error) {
	return api.ViewContext(ctx, triggerID, viewPayload, "views.open", "")
}
func (api *Client) UpdateViewContext(ctx context.Context, triggerID string, viewPayload ViewPayload, viewId string) (err error) {
	return api.ViewContext(ctx, triggerID, viewPayload, "views.update", viewId)
}

func (api *Client) PushViewContext(ctx context.Context, triggerID string, viewPayload ViewPayload) (err error) {
	return api.ViewContext(ctx, triggerID, viewPayload, "views.push", "")
}

func (api *Client) ViewContext(ctx context.Context, triggerID string, viewPayload ViewPayload, action string, viewId string) (err error) {
	if triggerID == "" {
		return ErrParametersMissing
	}
	encoded := []byte{}
	switch action {
	case "views.open", "views.push":
		req := ViewTrigger{
			TriggerID: triggerID,
			View:      viewPayload,
		}
		encoded, err = json.Marshal(req)
	case "views.update":
		req := UpdateViewTrigger{
			TriggerID: triggerID,
			View:      viewPayload,
			ViewId:    viewId,
		}
		encoded, err = json.Marshal(req)
	}

	if err != nil {
		return err
	}
	//fmt.Println(string(encoded))

	response := &ViewOpenResponse{}

	endpoint := api.endpoint + "views.open"
	if err := postJSON(ctx, api.httpclient, endpoint, api.token, encoded, response, api); err != nil {

		return err
	}

	if len(response.ViewResponseMetadata.Messages) > 0 {
		response.Ok = false
		response.Error += "\n" + strings.Join(response.ViewResponseMetadata.Messages, "\n")
	}

	return response.Err()
}
