package provider

import (
	"fmt"

	"github.com/elastic-infra/terraform-provider-ldap/internal/helper/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider creates a new LDAP provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"ldap_host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_HOST", nil),
				Description: "The LDAP server to connect to.",
			},
			"ldap_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_PORT", 389),
				Description: "The LDAP protocol port (default: 389).",
			},
			"bind_user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_BIND_USER", nil),
				Description: "Bind user to be used for authenticating on the LDAP server.",
			},
			"bind_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_BIND_PASSWORD", nil),
				Description: "Password to authenticate the Bind user.",
			},
			"start_tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_START_TLS", false),
				Description: "Upgrade TLS to secure the connection (default: false).",
			},
			"tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_TLS", false),
				Description: "Enable TLS encryption for LDAP (LDAPS) (default: false).",
			},
			"tls_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_TLS_INSECURE", false),
				Description: "Don't verify server TLS certificate (default: false).",
			},
			"use_gssapi": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_USE_GSSAPI", false),
				Description: "Use GSSAPI authentication instead of simple (default: false).",
			},
			"ccache": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KRB5CCNAME", nil),
				Description: "Kerberos CCache override location",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"ldap_object": resourceLDAPObject(),
		},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := &client.Config{
		LDAPHost:     d.Get("ldap_host").(string),
		LDAPPort:     d.Get("ldap_port").(int),
		BindUser:     d.Get("bind_user").(string),
		BindPassword: d.Get("bind_password").(string),
		CCache:       d.Get("ccache").(string),
		StartTLS:     d.Get("start_tls").(bool),
		TLS:          d.Get("tls").(bool),
		TLSInsecure:  d.Get("tls_insecure").(bool),
		UseGSSAPI:    d.Get("use_gssapi").(bool),
	}

	if config.UseGSSAPI && (config.CCache == "") {
		return nil, fmt.Errorf("When using GSSAPI setting `ccache` or $KRB5CCNAME is required.")
	}

	connection, err := client.DialAndBind(config)
	if err != nil {
		return nil, err
	}

	return connection, nil
}
