package authentication

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func standardizeCompanyName(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(s)), "-")
}

//Get subdomain from host name and handle respective errors
func getSubDomain(r *http.Request) (subDomain string, err error) {
	baseURl := r.Referer()
	basestring := string(baseURl[:])
	u, err := url.Parse(basestring)
	if err != nil {
		return "", err
	}
	domains := strings.Split(u.Hostname(), ".")
	if len(domains) <= 1 {
		return "", errors.New("Tenant not found on request")
	}
	//Not going to check if domain is longer just return the first
	// if len(domains) > 2 {
	// 	return "", errors.New("Domain not recognized")
	// }
	return domains[0], nil
}

func unauthorizedResponse(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
