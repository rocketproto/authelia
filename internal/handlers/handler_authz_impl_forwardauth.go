package handlers

import (
	"fmt"
	"net/url"

	"github.com/valyala/fasthttp"

	"github.com/authelia/authelia/v4/internal/authorization"
	"github.com/authelia/authelia/v4/internal/middlewares"
)

func handleAuthzGetObjectForwardAuth(ctx *middlewares.AutheliaCtx) (object authorization.Object, err error) {
	protocol, host, uri := ctx.XForwardedProto(), ctx.XForwardedHost(), ctx.XForwardedURI()

	var (
		targetURL *url.URL
		method    []byte
	)

	if targetURL, err = getRequestURIFromForwardedHeaders(protocol, host, uri); err != nil {
		return object, fmt.Errorf("failed to get target URL: %w", err)
	}

	if method = ctx.XForwardedMethod(); len(method) == 0 {
		return object, fmt.Errorf("header 'X-Forwarded-Method' is empty")
	}

	if hasInvalidMethodCharacters(method) {
		return object, fmt.Errorf("header 'X-Forwarded-Method' with value '%s' has invalid characters", method)
	}

	return authorization.NewObjectRaw(targetURL, method), nil
}

// Forward Auth
func handleAuthzUnauthorizedForwardAuth(ctx *middlewares.AutheliaCtx, authn *Authn, redirectionURL *url.URL) {
	var (
		statusCode int
	)

	ctx.Logger.Infof("Using: FORWARD AUTH METHOD") // REMOVE

	// Checks if request is expecting html or is from a browser
	if isRenderingHTML(ctx) {
		statusCode = fasthttp.StatusUnauthorized
	} else {
		statusCode = determineStatusCodeFromAuthn(authn)
	}

	// Checks if request is expecting html or is from a browser
	// switch {
	// case ctx.IsXHR() || !ctx.AcceptsMIME("text/html"):
	// 	statusCode = fasthttp.StatusUnauthorized
	// default:
	// 	// determine redirect type
	// 	switch authn.Object.Method {
	// 	case fasthttp.MethodGet, fasthttp.MethodOptions, fasthttp.MethodHead:
	// 		statusCode = fasthttp.StatusFound
	// 	default:
	// 		statusCode = fasthttp.StatusSeeOther
	// 	}
	// }

	handleSpecialRedirect(ctx, authn,redirectionURL, statusCode)

	// ctx.Logger.Infof(logFmtAuthzRedirect, authn.Object.String(), authn.Method, authn.Username, statusCode, redirectionURL)

	// // NOTE :) 401 Redirects

	// // Special Redirect Handling
	// switch authn.Object.Method {
	// case fasthttp.MethodHead:
	// 	ctx.SpecialRedirectNoBody(redirectionURL.String(), statusCode)
	// default:
	// 	ctx.SpecialRedirect(redirectionURL.String(), statusCode)
	// }
}
