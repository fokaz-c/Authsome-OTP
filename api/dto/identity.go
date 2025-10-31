package dto

type Identity_t struct {
	Identity_type     string `json:"identity_type"`
	Identity_source   string `json:"identity_source"`
	Identity_password string `json:"identity_password"`
	Tenant_username   string `json:"tenant_username"`
}
