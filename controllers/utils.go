package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/facs95/orderwell-backend/entity"
)

func handleTenantFailure(tenant entity.Tenant) {
	authProvider.DeleteTenant(tenant.OauthId)
	dbService.Tenant().Delete(tenant.ID)
}

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

func errorHttpResponse(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func successHttpResponseEncoder(w http.ResponseWriter, response interface{}) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
