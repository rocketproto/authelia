package handlers

import (
	"fmt"
	"net/url"

	"github.com/valyala/fasthttp"

	"github.com/authelia/authelia/v4/internal/authorization"
	"github.com/authelia/authelia/v4/internal/middlewares"
)

func handleAuthzGetObjectLegacy(ctx *middlewares.AutheliaCtx) (object authorization.Object, err error) {
	var (
		targetURL *url.URL
		method    []byte
	)

	if targetURL, err = ctx.GetXOriginalURLOrXForwardedURL(); err != nil {
		return object, fmt.Errorf("failed to get target URL: %w", err)
	}

	if method = ctx.XForwardedMethod(); len(method) == 0 {
		method = ctx.Method()
	}

	if hasInvalidMethodCharacters(method) {
		return object, fmt.Errorf("header 'X-Forwarded-Method' with value '%s' has invalid characters", method)
	}

	return authorization.NewObjectRaw(targetURL, method), nil
}

// Legacy Auth.
func handleAuthzUnauthorizedLegacy(ctx *middlewares.AutheliaCtx, authn *Authn, redirectionURL *url.URL) {
	var (
		statusCode int
	)

	ctx.Logger.Infof("Using: LEGACY AUTH METHOD") // REMOVE.

	if authn.Type == AuthnTypeAuthorization {
		handleAuthzUnauthorizedAuthorizationBasic(ctx, authn)

		return
	}

	// Checks if request is expecting html or is from a browser WITH URL CHECK.
	switch {
	case isRenderingHTML(ctx) || redirectionURL == nil:
		statusCode = fasthttp.StatusUnauthorized
	default:
		// Could this be more like :
		// `statusCode = determineStatusCodeFromAuthn(authn, "", fasthttp.StatusFound)`
		// and that be an addition to the switch case?
		if authn.Object.Method == "" {
			statusCode = fasthttp.StatusFound
		} else {
			statusCode = determineStatusCodeFromAuthn(authn)
		}
	}

	// Checks if request is expecting html or is from a browser WITH URL CHECK
	// switch {
	// case ctx.IsXHR() || !ctx.AcceptsMIME("text/html") || redirectionURL == nil:
	// 	statusCode = fasthttp.StatusUnauthorized
	// default:
	// 	// determine redirect type
	// 	switch authn.Object.Method {
	// 	case fasthttp.MethodGet, fasthttp.MethodOptions, fasthttp.MethodHead, "": // Noting Empty Method, unlike others
	// 		statusCode = fasthttp.StatusFound
	// 	default:
	// 		statusCode = fasthttp.StatusSeeOther
	// 	}
	// }.

	if redirectionURL != nil {
		ctx.Logger.Infof(logFmtAuthzRedirect, authn.Object.URL.String(), authn.Method, authn.Username, statusCode, redirectionURL)

		handleSpecialRedirect(ctx, authn, redirectionURL, statusCode)
		// NOTE :) 401 Redirects.

		// Special Redirect Handling
		// switch authn.Object.Method {
		// case fasthttp.MethodHead:
		// 	ctx.SpecialRedirectNoBody(redirectionURL.String(), statusCode)
		// default:
		// 	ctx.SpecialRedirect(redirectionURL.String(), statusCode)
		// }.
	} else {
		ctx.Logger.Infof("Access to %s (method %s) is not authorized to user %s, responding with status code %d", authn.Object.URL.String(), authn.Method, authn.Username, statusCode)
		ctx.ReplyUnauthorized()
	}
}

// Legacy Auth.
func handleAuthzForbiddenLegacy(ctx *middlewares.AutheliaCtx, authn *Authn, redirectionURL *url.URL) {
	var (
		statusCode int
	)

	ctx.Logger.Infof("Using: LEGACY AUTH METHOD") // REMOVE.

	if authn.Type == AuthnTypeAuthorization {
		// handleAuthzUnauthorizedAuthorizationBasic(ctx, authn).

		// Correct?
		ctx.ReplyForbidden()
		return
	}

	// Checks if request is expecting html or is from a browser WITH URL CHECK.
	switch {
	case isRenderingHTML(ctx) || redirectionURL == nil:
		statusCode = fasthttp.StatusUnauthorized
	default:
		// Could this be more like :
		// `statusCode = determineStatusCodeFromAuthn(authn, "", fasthttp.StatusFound)`
		// and that be an addition to the switch case?
		if authn.Object.Method == "" {
			statusCode = fasthttp.StatusFound
		} else {
			statusCode = determineStatusCodeFromAuthn(authn)
		}
	}

	if redirectionURL != nil {
		handleSpecialRedirect(ctx, authn, redirectionURL, statusCode)
	} else {
		ctx.Logger.Infof("Access to %s (method %s) is forbidden for user %s, responding with status code %d", authn.Object.URL.String(), authn.Method, authn.Username, statusCode)
		ctx.ReplyForbidden()
	}
}
