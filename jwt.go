package main

import (
	oidc "github.com/coreos/go-oidc"
	_ "github.com/lestrrat-go/jwx/jwk"
)

// Token represents a JWT
type Token struct {
	token *oidc.IDToken
}


func (t Token) Roles(client string) []string {
	var m map[string]interface{}
	okk := t.token.Claims(&m)
	var rawResourceAccess = m["resource_access"]
	if okk != nil {
		return nil
	}
	resourceAccess, ok := rawResourceAccess.(map[string]interface{})
	if !ok {
		return nil
	}
	rawClientResource, ok := resourceAccess[client]
	if !ok {
		return nil
	}
	clientResource, ok := rawClientResource.(map[string]interface{})
	if !ok {
		return nil
	}
	rawRoles, ok := clientResource["roles"]
	if !ok {
		return nil
	}
	roles, ok := rawRoles.([]interface{})
	if !ok {
		return nil
	}
	var roleStrings []string
	for _, r := range roles {
		rs, ok := r.(string)
		if ok {
			roleStrings = append(roleStrings, rs)
		}
	}
	return roleStrings
}

func Find(slice []string, value string) (int, bool) {
	for i, item := range slice {
		if item == value {
			return i, true
		}
	}
	return -1, false
}

