// Code generated by go-swagger; DO NOT EDIT.

package events

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetEventsNetworkIDStreamNameParams creates a new GetEventsNetworkIDStreamNameParams object
// with the default values initialized.
func NewGetEventsNetworkIDStreamNameParams() *GetEventsNetworkIDStreamNameParams {
	var ()
	return &GetEventsNetworkIDStreamNameParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetEventsNetworkIDStreamNameParamsWithTimeout creates a new GetEventsNetworkIDStreamNameParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetEventsNetworkIDStreamNameParamsWithTimeout(timeout time.Duration) *GetEventsNetworkIDStreamNameParams {
	var ()
	return &GetEventsNetworkIDStreamNameParams{

		timeout: timeout,
	}
}

// NewGetEventsNetworkIDStreamNameParamsWithContext creates a new GetEventsNetworkIDStreamNameParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetEventsNetworkIDStreamNameParamsWithContext(ctx context.Context) *GetEventsNetworkIDStreamNameParams {
	var ()
	return &GetEventsNetworkIDStreamNameParams{

		Context: ctx,
	}
}

// NewGetEventsNetworkIDStreamNameParamsWithHTTPClient creates a new GetEventsNetworkIDStreamNameParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetEventsNetworkIDStreamNameParamsWithHTTPClient(client *http.Client) *GetEventsNetworkIDStreamNameParams {
	var ()
	return &GetEventsNetworkIDStreamNameParams{
		HTTPClient: client,
	}
}

/*GetEventsNetworkIDStreamNameParams contains all the parameters to send to the API endpoint
for the get events network ID stream name operation typically these are written to a http.Request
*/
type GetEventsNetworkIDStreamNameParams struct {

	/*EventType
	  The type of event to filter the query with.

	*/
	EventType *string
	/*HardwareID
	  The hardware ID to filter the query with.

	*/
	HardwareID *string
	/*NetworkID
	  Network ID

	*/
	NetworkID string
	/*StreamName
	  The user-specified string to categorize events

	*/
	StreamName string
	/*Tag
	  The event tag to filter the query with.

	*/
	Tag *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) WithTimeout(timeout time.Duration) *GetEventsNetworkIDStreamNameParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) WithContext(ctx context.Context) *GetEventsNetworkIDStreamNameParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) WithHTTPClient(client *http.Client) *GetEventsNetworkIDStreamNameParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithEventType adds the eventType to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) WithEventType(eventType *string) *GetEventsNetworkIDStreamNameParams {
	o.SetEventType(eventType)
	return o
}

// SetEventType adds the eventType to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) SetEventType(eventType *string) {
	o.EventType = eventType
}

// WithHardwareID adds the hardwareID to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) WithHardwareID(hardwareID *string) *GetEventsNetworkIDStreamNameParams {
	o.SetHardwareID(hardwareID)
	return o
}

// SetHardwareID adds the hardwareId to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) SetHardwareID(hardwareID *string) {
	o.HardwareID = hardwareID
}

// WithNetworkID adds the networkID to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) WithNetworkID(networkID string) *GetEventsNetworkIDStreamNameParams {
	o.SetNetworkID(networkID)
	return o
}

// SetNetworkID adds the networkId to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) SetNetworkID(networkID string) {
	o.NetworkID = networkID
}

// WithStreamName adds the streamName to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) WithStreamName(streamName string) *GetEventsNetworkIDStreamNameParams {
	o.SetStreamName(streamName)
	return o
}

// SetStreamName adds the streamName to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) SetStreamName(streamName string) {
	o.StreamName = streamName
}

// WithTag adds the tag to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) WithTag(tag *string) *GetEventsNetworkIDStreamNameParams {
	o.SetTag(tag)
	return o
}

// SetTag adds the tag to the get events network ID stream name params
func (o *GetEventsNetworkIDStreamNameParams) SetTag(tag *string) {
	o.Tag = tag
}

// WriteToRequest writes these params to a swagger request
func (o *GetEventsNetworkIDStreamNameParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.EventType != nil {

		// query param event_type
		var qrEventType string
		if o.EventType != nil {
			qrEventType = *o.EventType
		}
		qEventType := qrEventType
		if qEventType != "" {
			if err := r.SetQueryParam("event_type", qEventType); err != nil {
				return err
			}
		}

	}

	if o.HardwareID != nil {

		// query param hardware_id
		var qrHardwareID string
		if o.HardwareID != nil {
			qrHardwareID = *o.HardwareID
		}
		qHardwareID := qrHardwareID
		if qHardwareID != "" {
			if err := r.SetQueryParam("hardware_id", qHardwareID); err != nil {
				return err
			}
		}

	}

	// path param network_id
	if err := r.SetPathParam("network_id", o.NetworkID); err != nil {
		return err
	}

	// path param stream_name
	if err := r.SetPathParam("stream_name", o.StreamName); err != nil {
		return err
	}

	if o.Tag != nil {

		// query param tag
		var qrTag string
		if o.Tag != nil {
			qrTag = *o.Tag
		}
		qTag := qrTag
		if qTag != "" {
			if err := r.SetQueryParam("tag", qTag); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
