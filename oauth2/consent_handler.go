/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package oauth2

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/pkg/errors"
)

const (
	ConsentRequestAccepted = "accepted"
	ConsentRequestRejected = "rejected"

	ConsentRequestPath = "/oauth2/consent/requests"

	ConsentResource = "oauth2:consent:requests:%s"
	ConsentScope    = "hydra.consent"
)

type ConsentSessionHandler struct {
	H herodot.Writer
	M ConsentRequestManager

	ResourcePrefix string
}

func (h *ConsentSessionHandler) PrefixResource(resource string) string {
	if h.ResourcePrefix == "" {
		h.ResourcePrefix = "rn:hydra"
	}

	if h.ResourcePrefix[len(h.ResourcePrefix)-1] == ':' {
		h.ResourcePrefix = h.ResourcePrefix[:len(h.ResourcePrefix)-1]
	}

	return h.ResourcePrefix + ":" + resource
}

func (h *ConsentSessionHandler) SetRoutes(r *httprouter.Router) {
	r.GET(ConsentRequestPath+"/:id", h.FetchConsentRequest)
	r.PATCH(ConsentRequestPath+"/:id/reject", h.RejectConsentRequestHandler)
	r.PATCH(ConsentRequestPath+"/:id/accept", h.AcceptConsentRequestHandler)
}

// swagger:route GET /oauth2/consent/requests/{id} oAuth2 getOAuth2ConsentRequest
//
// Receive consent request information
//
// Call this endpoint to receive information on consent requests. The consent request id is usually transmitted via the URL query `consent`.
// For example: `http://consent-app.mydomain.com/?consent=1234abcd`
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: oAuth2ConsentRequest
//       401: genericError
//       500: genericError
func (h *ConsentSessionHandler) FetchConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if session, err := h.M.GetConsentRequest(ps.ByName("id")); err != nil {
		h.H.WriteError(w, r, err)
		return
	} else {
		h.H.Write(w, r, session)
	}
}

// swagger:route PATCH /oauth2/consent/requests/{id}/reject oAuth2 rejectOAuth2ConsentRequest
//
// Reject a consent request
//
// Call this endpoint to reject a consent request. This usually happens when a user denies access rights to an
// application.
//
//
// The consent request id is usually transmitted via the URL query `consent`.
// For example: `http://consent-app.mydomain.com/?consent=1234abcd`
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       204: emptyResponse
//       401: genericError
//       500: genericError
func (h *ConsentSessionHandler) RejectConsentRequestHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var payload RejectConsentRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if err := h.M.RejectConsentRequest(ps.ByName("id"), &payload); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route PATCH /oauth2/consent/requests/{id}/accept oAuth2 acceptOAuth2ConsentRequest
//
// Accept a consent request
//
// Call this endpoint to accept a consent request. This usually happens when a user agrees to give access rights to
// an application.
//
//
// The consent request id is usually transmitted via the URL query `consent`.
// For example: `http://consent-app.mydomain.com/?consent=1234abcd`
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       204: emptyResponse
//       401: genericError
//       500: genericError
func (h *ConsentSessionHandler) AcceptConsentRequestHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var payload AcceptConsentRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.H.WriteError(w, r, errors.WithStack(err))
		return
	}

	if err := h.M.AcceptConsentRequest(ps.ByName("id"), &payload); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
