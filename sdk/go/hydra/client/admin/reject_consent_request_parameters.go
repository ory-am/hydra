// Code generated by go-swagger; DO NOT EDIT.

package admin

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

	models "github.com/ory/hydra/sdk/go/hydra/models"
)

// NewRejectConsentRequestParams creates a new RejectConsentRequestParams object
// with the default values initialized.
func NewRejectConsentRequestParams() *RejectConsentRequestParams {
	var ()
	return &RejectConsentRequestParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewRejectConsentRequestParamsWithTimeout creates a new RejectConsentRequestParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewRejectConsentRequestParamsWithTimeout(timeout time.Duration) *RejectConsentRequestParams {
	var ()
	return &RejectConsentRequestParams{

		timeout: timeout,
	}
}

// NewRejectConsentRequestParamsWithContext creates a new RejectConsentRequestParams object
// with the default values initialized, and the ability to set a context for a request
func NewRejectConsentRequestParamsWithContext(ctx context.Context) *RejectConsentRequestParams {
	var ()
	return &RejectConsentRequestParams{

		Context: ctx,
	}
}

// NewRejectConsentRequestParamsWithHTTPClient creates a new RejectConsentRequestParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewRejectConsentRequestParamsWithHTTPClient(client *http.Client) *RejectConsentRequestParams {
	var ()
	return &RejectConsentRequestParams{
		HTTPClient: client,
	}
}

/*RejectConsentRequestParams contains all the parameters to send to the API endpoint
for the reject consent request operation typically these are written to a http.Request
*/
type RejectConsentRequestParams struct {

	/*Body*/
	Body *models.RequestDeniedError
	/*Challenge*/
	Challenge string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the reject consent request params
func (o *RejectConsentRequestParams) WithTimeout(timeout time.Duration) *RejectConsentRequestParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the reject consent request params
func (o *RejectConsentRequestParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the reject consent request params
func (o *RejectConsentRequestParams) WithContext(ctx context.Context) *RejectConsentRequestParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the reject consent request params
func (o *RejectConsentRequestParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the reject consent request params
func (o *RejectConsentRequestParams) WithHTTPClient(client *http.Client) *RejectConsentRequestParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the reject consent request params
func (o *RejectConsentRequestParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the reject consent request params
func (o *RejectConsentRequestParams) WithBody(body *models.RequestDeniedError) *RejectConsentRequestParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the reject consent request params
func (o *RejectConsentRequestParams) SetBody(body *models.RequestDeniedError) {
	o.Body = body
}

// WithChallenge adds the challenge to the reject consent request params
func (o *RejectConsentRequestParams) WithChallenge(challenge string) *RejectConsentRequestParams {
	o.SetChallenge(challenge)
	return o
}

// SetChallenge adds the challenge to the reject consent request params
func (o *RejectConsentRequestParams) SetChallenge(challenge string) {
	o.Challenge = challenge
}

// WriteToRequest writes these params to a swagger request
func (o *RejectConsentRequestParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// query param challenge
	qrChallenge := o.Challenge
	qChallenge := qrChallenge
	if qChallenge != "" {
		if err := r.SetQueryParam("challenge", qChallenge); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
