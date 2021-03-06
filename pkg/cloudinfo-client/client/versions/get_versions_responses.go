// Code generated by go-swagger; DO NOT EDIT.

package versions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/banzaicloud/cloudinfo/pkg/cloudinfo-client/models"
)

// GetVersionsReader is a Reader for the GetVersions structure.
type GetVersionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetVersionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetVersionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetVersionsOK creates a GetVersionsOK with default headers values
func NewGetVersionsOK() *GetVersionsOK {
	return &GetVersionsOK{}
}

/*GetVersionsOK handles this case with default header values.

VersionsResponse
*/
type GetVersionsOK struct {
	Payload *models.VersionsResponse
}

func (o *GetVersionsOK) Error() string {
	return fmt.Sprintf("[GET /providers/{provider}/services/{service}/regions/{region}/versions][%d] getVersionsOK  %+v", 200, o.Payload)
}

func (o *GetVersionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.VersionsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
