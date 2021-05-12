// Code generated by go-swagger; DO NOT EDIT.

package subscribers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "magma/orc8r/cloud/go/obsidian/swagger/v1/models"
)

// GetLTENetworkIDMsisdnsMsisdnReader is a Reader for the GetLTENetworkIDMsisdnsMsisdn structure.
type GetLTENetworkIDMsisdnsMsisdnReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetLTENetworkIDMsisdnsMsisdnReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetLTENetworkIDMsisdnsMsisdnOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetLTENetworkIDMsisdnsMsisdnDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetLTENetworkIDMsisdnsMsisdnOK creates a GetLTENetworkIDMsisdnsMsisdnOK with default headers values
func NewGetLTENetworkIDMsisdnsMsisdnOK() *GetLTENetworkIDMsisdnsMsisdnOK {
	return &GetLTENetworkIDMsisdnsMsisdnOK{}
}

/*GetLTENetworkIDMsisdnsMsisdnOK handles this case with default header values.

Subscriber ID
*/
type GetLTENetworkIDMsisdnsMsisdnOK struct {
	Payload models.SubscriberID
}

func (o *GetLTENetworkIDMsisdnsMsisdnOK) Error() string {
	return fmt.Sprintf("[GET /lte/{network_id}/msisdns/{msisdn}][%d] getLteNetworkIdMsisdnsMsisdnOK  %+v", 200, o.Payload)
}

func (o *GetLTENetworkIDMsisdnsMsisdnOK) GetPayload() models.SubscriberID {
	return o.Payload
}

func (o *GetLTENetworkIDMsisdnsMsisdnOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetLTENetworkIDMsisdnsMsisdnDefault creates a GetLTENetworkIDMsisdnsMsisdnDefault with default headers values
func NewGetLTENetworkIDMsisdnsMsisdnDefault(code int) *GetLTENetworkIDMsisdnsMsisdnDefault {
	return &GetLTENetworkIDMsisdnsMsisdnDefault{
		_statusCode: code,
	}
}

/*GetLTENetworkIDMsisdnsMsisdnDefault handles this case with default header values.

Unexpected Error
*/
type GetLTENetworkIDMsisdnsMsisdnDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get LTE network ID msisdns msisdn default response
func (o *GetLTENetworkIDMsisdnsMsisdnDefault) Code() int {
	return o._statusCode
}

func (o *GetLTENetworkIDMsisdnsMsisdnDefault) Error() string {
	return fmt.Sprintf("[GET /lte/{network_id}/msisdns/{msisdn}][%d] GetLTENetworkIDMsisdnsMsisdn default  %+v", o._statusCode, o.Payload)
}

func (o *GetLTENetworkIDMsisdnsMsisdnDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetLTENetworkIDMsisdnsMsisdnDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
