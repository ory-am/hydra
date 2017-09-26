package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/firewall"
	"github.com/pkg/errors"
)

const (
	ConsentRequestAccepted = "accepted"
	ConsentRequestRejected = "rejected"

	ConsentSessionPath = "/oauth2/consent/requests"

	ConsentResource = "rn:hydra:oauth2:consent:requests:%s"
	ConsentScope    = "hydra.consent"
)

type ConsentSessionHandler struct {
	H herodot.Writer
	M ConsentRequestManager
	W firewall.Firewall
}

func (h *ConsentSessionHandler) SetRoutes(r *httprouter.Router) {
	r.GET(ConsentSessionPath+"/:id", h.FetchConsentRequest)
	r.PATCH(ConsentSessionPath+"/:id/reject", h.RejectConsentRequestHandler)
	r.PATCH(ConsentSessionPath+"/:id/accept", h.AcceptConsentRequestHandler)
}

// swagger:route GET /oauth2/consent/requests/{id} oauth2 consent getConsentRequest
//
// Receive consent request information
//
// Call this endpoint to receive information on consent requests.
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:oauth2:consent:requests:<request-id>"],
//    "actions": ["get"],
//    "effect": "allow"
//  }
//  ```
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2: hydra.consent
//
//     Responses:
//       200: oauthConsentRequest
//       401: genericError
//       500: genericError
func (h *ConsentSessionHandler) FetchConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, err := h.W.TokenAllowed(r.Context(), h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(ConsentResource, ps.ByName("id")),
		Action:   "get",
	}, ConsentScope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

	if session, err := h.M.GetConsentRequest(ps.ByName("id")); err != nil {
		h.H.WriteError(w, r, err)
		return
	} else {
		h.H.Write(w, r, session)
	}
}

// swagger:route PATCH /oauth2/consent/requests/{id}/reject oauth2 consent rejectConsentRequest
//
// Reject a consent request
//
// Call this endpoint to reject a consent request. This usually happens when a user denies access rights to an
// application.
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:oauth2:consent:requests:<request-id>"],
//    "actions": ["reject"],
//    "effect": "allow"
//  }
//  ```
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2: hydra.consent
//
//     Responses:
//       204: emptyResponse
//       401: genericError
//       500: genericError
func (h *ConsentSessionHandler) RejectConsentRequestHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, err := h.W.TokenAllowed(r.Context(), h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(ConsentResource, ps.ByName("id")),
		Action:   "reject",
	}, ConsentScope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

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

// swagger:route PATCH /oauth2/consent/requests/{id}/accept oauth2 consent acceptConsentRequest
//
// Accept a consent request
//
// Call this endpoint to accept a consent request. This usually happens when a user agrees to give access rights to
// an application.
//
// The subject making the request needs to be assigned to a policy containing:
//
//  ```
//  {
//    "resources": ["rn:hydra:oauth2:consent:requests:<request-id>"],
//    "actions": ["accept"],
//    "effect": "allow"
//  }
//  ```
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       oauth2: hydra.consent
//
//     Responses:
//       204: emptyResponse
//       401: genericError
//       500: genericError
func (h *ConsentSessionHandler) AcceptConsentRequestHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, err := h.W.TokenAllowed(r.Context(), h.W.TokenFromRequest(r), &firewall.TokenAccessRequest{
		Resource: fmt.Sprintf(ConsentResource, ps.ByName("id")),
		Action:   "accept",
	}, ConsentScope); err != nil {
		h.H.WriteError(w, r, err)
		return
	}

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
