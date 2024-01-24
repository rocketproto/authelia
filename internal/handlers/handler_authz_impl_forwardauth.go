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

// Forward Auth.
func handleAuthzUnauthorizedForwardAuth(ctx *middlewares.AutheliaCtx, authn *Authn, redirectionURL *url.URL) {
	var (
		statusCode int
	)

	ctx.Logger.Infof("Using: FORWARD AUTH METHOD") // REMOVE.

	// Checks if request is expecting html or is from a browser.
	if isRenderingHTML(ctx) {
		statusCode = fasthttp.StatusUnauthorized
	} else {
		statusCode = determineStatusCodeFromAuthn(authn)
	}

	handleSpecialRedirect(ctx, authn, redirectionURL, statusCode)
}

// Forward Auth.
func handleAuthzForbiddenForwardAuth(ctx *middlewares.AutheliaCtx, authn *Authn, redirectionURL *url.URL) {
	var (
		statusCode int
	)

	ctx.Logger.Infof("Using: FORWARD AUTH METHOD") // REMOVE.

	// Checks if request is expecting html or is from a browser.
	if isRenderingHTML(ctx) {
		statusCode = fasthttp.StatusForbidden
	} else {
		statusCode = determineStatusCodeFromAuthn(authn)
	}

	handleSpecialRedirect(ctx, authn, redirectionURL, statusCode)
}
