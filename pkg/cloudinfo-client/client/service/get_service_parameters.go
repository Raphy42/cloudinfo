// Code generated by go-swagger; DO NOT EDIT.

package service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetServiceParams creates a new GetServiceParams object
// with the default values initialized.
func NewGetServiceParams() *GetServiceParams {
	var ()
	return &GetServiceParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetServiceParamsWithTimeout creates a new GetServiceParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetServiceParamsWithTimeout(timeout time.Duration) *GetServiceParams {
	var ()
	return &GetServiceParams{

		timeout: timeout,
	}
}

// NewGetServiceParamsWithContext creates a new GetServiceParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetServiceParamsWithContext(ctx context.Context) *GetServiceParams {
	var ()
	return &GetServiceParams{

		Context: ctx,
	}
}

// NewGetServiceParamsWithHTTPClient creates a new GetServiceParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetServiceParamsWithHTTPClient(client *http.Client) *GetServiceParams {
	var ()
	return &GetServiceParams{
		HTTPClient: client,
	}
}

/*GetServiceParams contains all the parameters to send to the API endpoint
for the get service operation typically these are written to a http.Request
*/
type GetServiceParams struct {

	/*Provider*/
	Provider string
	/*Service*/
	Service string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get service params
func (o *GetServiceParams) WithTimeout(timeout time.Duration) *GetServiceParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get service params
func (o *GetServiceParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get service params
func (o *GetServiceParams) WithContext(ctx context.Context) *GetServiceParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get service params
func (o *GetServiceParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get service params
func (o *GetServiceParams) WithHTTPClient(client *http.Client) *GetServiceParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get service params
func (o *GetServiceParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithProvider adds the provider to the get service params
func (o *GetServiceParams) WithProvider(provider string) *GetServiceParams {
	o.SetProvider(provider)
	return o
}

// SetProvider adds the provider to the get service params
func (o *GetServiceParams) SetProvider(provider string) {
	o.Provider = provider
}

// WithService adds the service to the get service params
func (o *GetServiceParams) WithService(service string) *GetServiceParams {
	o.SetService(service)
	return o
}

// SetService adds the service to the get service params
func (o *GetServiceParams) SetService(service string) {
	o.Service = service
}

// WriteToRequest writes these params to a swagger request
func (o *GetServiceParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param provider
	if err := r.SetPathParam("provider", o.Provider); err != nil {
		return err
	}

	// path param service
	if err := r.SetPathParam("service", o.Service); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
