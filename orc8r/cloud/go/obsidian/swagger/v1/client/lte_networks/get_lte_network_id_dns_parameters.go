// Code generated by go-swagger; DO NOT EDIT.

package lte_networks

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

// NewGetLTENetworkIDDNSParams creates a new GetLTENetworkIDDNSParams object
// with the default values initialized.
func NewGetLTENetworkIDDNSParams() *GetLTENetworkIDDNSParams {
	var ()
	return &GetLTENetworkIDDNSParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetLTENetworkIDDNSParamsWithTimeout creates a new GetLTENetworkIDDNSParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetLTENetworkIDDNSParamsWithTimeout(timeout time.Duration) *GetLTENetworkIDDNSParams {
	var ()
	return &GetLTENetworkIDDNSParams{

		timeout: timeout,
	}
}

// NewGetLTENetworkIDDNSParamsWithContext creates a new GetLTENetworkIDDNSParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetLTENetworkIDDNSParamsWithContext(ctx context.Context) *GetLTENetworkIDDNSParams {
	var ()
	return &GetLTENetworkIDDNSParams{

		Context: ctx,
	}
}

// NewGetLTENetworkIDDNSParamsWithHTTPClient creates a new GetLTENetworkIDDNSParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetLTENetworkIDDNSParamsWithHTTPClient(client *http.Client) *GetLTENetworkIDDNSParams {
	var ()
	return &GetLTENetworkIDDNSParams{
		HTTPClient: client,
	}
}

/*GetLTENetworkIDDNSParams contains all the parameters to send to the API endpoint
for the get LTE network ID DNS operation typically these are written to a http.Request
*/
type GetLTENetworkIDDNSParams struct {

	/*NetworkID
	  Network ID

	*/
	NetworkID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get LTE network ID DNS params
func (o *GetLTENetworkIDDNSParams) WithTimeout(timeout time.Duration) *GetLTENetworkIDDNSParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get LTE network ID DNS params
func (o *GetLTENetworkIDDNSParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get LTE network ID DNS params
func (o *GetLTENetworkIDDNSParams) WithContext(ctx context.Context) *GetLTENetworkIDDNSParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get LTE network ID DNS params
func (o *GetLTENetworkIDDNSParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get LTE network ID DNS params
func (o *GetLTENetworkIDDNSParams) WithHTTPClient(client *http.Client) *GetLTENetworkIDDNSParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get LTE network ID DNS params
func (o *GetLTENetworkIDDNSParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithNetworkID adds the networkID to the get LTE network ID DNS params
func (o *GetLTENetworkIDDNSParams) WithNetworkID(networkID string) *GetLTENetworkIDDNSParams {
	o.SetNetworkID(networkID)
	return o
}

// SetNetworkID adds the networkId to the get LTE network ID DNS params
func (o *GetLTENetworkIDDNSParams) SetNetworkID(networkID string) {
	o.NetworkID = networkID
}

// WriteToRequest writes these params to a swagger request
func (o *GetLTENetworkIDDNSParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param network_id
	if err := r.SetPathParam("network_id", o.NetworkID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
