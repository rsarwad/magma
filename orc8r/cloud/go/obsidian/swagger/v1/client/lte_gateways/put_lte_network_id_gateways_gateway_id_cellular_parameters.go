// Code generated by go-swagger; DO NOT EDIT.

package lte_gateways

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

	models "magma/orc8r/cloud/go/obsidian/swagger/v1/models"
)

// NewPutLTENetworkIDGatewaysGatewayIDCellularParams creates a new PutLTENetworkIDGatewaysGatewayIDCellularParams object
// with the default values initialized.
func NewPutLTENetworkIDGatewaysGatewayIDCellularParams() *PutLTENetworkIDGatewaysGatewayIDCellularParams {
	var ()
	return &PutLTENetworkIDGatewaysGatewayIDCellularParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPutLTENetworkIDGatewaysGatewayIDCellularParamsWithTimeout creates a new PutLTENetworkIDGatewaysGatewayIDCellularParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPutLTENetworkIDGatewaysGatewayIDCellularParamsWithTimeout(timeout time.Duration) *PutLTENetworkIDGatewaysGatewayIDCellularParams {
	var ()
	return &PutLTENetworkIDGatewaysGatewayIDCellularParams{

		timeout: timeout,
	}
}

// NewPutLTENetworkIDGatewaysGatewayIDCellularParamsWithContext creates a new PutLTENetworkIDGatewaysGatewayIDCellularParams object
// with the default values initialized, and the ability to set a context for a request
func NewPutLTENetworkIDGatewaysGatewayIDCellularParamsWithContext(ctx context.Context) *PutLTENetworkIDGatewaysGatewayIDCellularParams {
	var ()
	return &PutLTENetworkIDGatewaysGatewayIDCellularParams{

		Context: ctx,
	}
}

// NewPutLTENetworkIDGatewaysGatewayIDCellularParamsWithHTTPClient creates a new PutLTENetworkIDGatewaysGatewayIDCellularParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPutLTENetworkIDGatewaysGatewayIDCellularParamsWithHTTPClient(client *http.Client) *PutLTENetworkIDGatewaysGatewayIDCellularParams {
	var ()
	return &PutLTENetworkIDGatewaysGatewayIDCellularParams{
		HTTPClient: client,
	}
}

/*PutLTENetworkIDGatewaysGatewayIDCellularParams contains all the parameters to send to the API endpoint
for the put LTE network ID gateways gateway ID cellular operation typically these are written to a http.Request
*/
type PutLTENetworkIDGatewaysGatewayIDCellularParams struct {

	/*Config
	  New cellular configuration

	*/
	Config *models.GatewayCellularConfigs
	/*GatewayID
	  Gateway ID

	*/
	GatewayID string
	/*NetworkID
	  Network ID

	*/
	NetworkID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the put LTE network ID gateways gateway ID cellular params
func (o *PutLTENetworkIDGatewaysGatewayIDCellularParams) WithTimeout(timeout time.Duration) *PutLTENetworkIDGatewaysGatewayIDCellularParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the put LTE network ID gateways gateway ID cellular params
func (o *PutLTENetworkIDGatewaysGatewayIDCellularParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the put LTE network ID gateways gateway ID cellular params
func (o *PutLTENetworkIDGatewaysGatewayIDCellularParams) WithContext(ctx context.Context) *PutLTENetworkIDGatewaysGatewayIDCellularParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the put LTE network ID gateways gateway ID cellular params
func (o *PutLTENetworkIDGatewaysGatewayIDCellularParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the put LTE network ID gateways gateway ID cellular params
func (o *PutLTENetworkIDGatewaysGatewayIDCellularParams) WithHTTPClient(client *http.Client) *PutLTENetworkIDGatewaysGatewayIDCellularParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the put LTE network ID gateways gateway ID cellular params
func (o *PutLTENetworkIDGatewaysGatewayIDCellularParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithConfig adds the config to the put LTE network ID gateways gateway ID cellular params
func (o *PutLTENetworkIDGatewaysGatewayIDCellularParams) WithConfig(config *models.GatewayCellularConfigs) *PutLTENetworkIDGatewaysGatewayIDCellularParams {
	o.SetConfig(config)
	return o
}

// SetConfig adds the config to the put LTE network ID gateways gateway ID cellular params
func (o *PutLTENetworkIDGatewaysGatewayIDCellularParams) SetConfig(config *models.GatewayCellularConfigs) {
	o.Config = config
}

// WithGatewayID adds the gatewayID to the put LTE network ID gateways gateway ID cellular params
func (o *PutLTENetworkIDGatewaysGatewayIDCellularParams) WithGatewayID(gatewayID string) *PutLTENetworkIDGatewaysGatewayIDCellularParams {
	o.SetGatewayID(gatewayID)
	return o
}

// SetGatewayID adds the gatewayId to the put LTE network ID gateways gateway ID cellular params
func (o *PutLTENetworkIDGatewaysGatewayIDCellularParams) SetGatewayID(gatewayID string) {
	o.GatewayID = gatewayID
}

// WithNetworkID adds the networkID to the put LTE network ID gateways gateway ID cellular params
func (o *PutLTENetworkIDGatewaysGatewayIDCellularParams) WithNetworkID(networkID string) *PutLTENetworkIDGatewaysGatewayIDCellularParams {
	o.SetNetworkID(networkID)
	return o
}

// SetNetworkID adds the networkId to the put LTE network ID gateways gateway ID cellular params
func (o *PutLTENetworkIDGatewaysGatewayIDCellularParams) SetNetworkID(networkID string) {
	o.NetworkID = networkID
}

// WriteToRequest writes these params to a swagger request
func (o *PutLTENetworkIDGatewaysGatewayIDCellularParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Config != nil {
		if err := r.SetBodyParam(o.Config); err != nil {
			return err
		}
	}

	// path param gateway_id
	if err := r.SetPathParam("gateway_id", o.GatewayID); err != nil {
		return err
	}

	// path param network_id
	if err := r.SetPathParam("network_id", o.NetworkID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
