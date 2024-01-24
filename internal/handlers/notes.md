StatusForbidden 403 just responds 403 text
StatusUnauthorized 401 redirects to Auth?

Return a 403 status message with a location and the proxy is in charge of the redirect. Proxy has to change 403 to 302 and make sure the header is included
