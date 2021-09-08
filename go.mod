module github.com/dequbed/terraform-provider-ldap/v2

go 1.14

require (
	github.com/go-ldap/ldap/v3 v3.2.4
	github.com/hashicorp/terraform-plugin-docs v0.4.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.5.0
	github.com/jcmturner/gokrb5/v8 v8.4.2 // indirect
)

replace github.com/go-ldap/ldap/v3 => github.com/dequbed/ldap/v3 v3.4.2
