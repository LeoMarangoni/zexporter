# zexporter
Zimbra Provisioning Exporter

Tool for exporting provisioning information(like accounts, distributionlists, cos, etc) in zimbra zmprov format.

### Implemented Features
- [X] Accounts
- [X] Distribution Lists
- [X] Distribution List Members
- [ ] Class of service
- [ ] Resources


```sh
Usage of zexporter:
  -attrs string
    	Add a custom list of attributes separatted by comma. (default "displayName,zimbraAccountstatus,givenName,sn,zimbraIsAdminAccount,zimbraPrefMailForwardingAddress,zimbraPrefOutOfOfficeCacheDuration,zimbraPrefOutOfOfficeDirectAddress,zimbraPrefOutOfOfficeFromDate,zimbraPrefOutOfOfficeReply,zimbraPrefOutOfOfficeReplyEnabled,zimbraPrefOutOfOfficeUntilDate,zimbraPrefHtmlEditorDefaultFontColor,zimbraPrefHtmlEditorDefaultFontFamily,zimbraPrefHtmlEditorDefaultFontSize,zimbraPrefMessageViewHtmlPreferred,zimbraMailSieveScript,zimbraPrefComposeFormat,zimbraPrefGroupMailBy,zimbraSignatureName,zimbraSignatureId,zimbraPrefMailSignatureHTML,zimbraPrefMailSignature,zimbraPrefForwardReplySignatureId,zimbraPrefDefaultSignatureId,userPassword")
  -domain string
    	Domain to export. Leave blank will export all. (default "*")
  -ldap_host string
    	Zimbra LDAP Host to use.
  -ldap_password string
    	Password to authenticate in ldap server.
  -ldap_port string
    	LDAP BIND PORT. (default "389")
  -ldap_user string
    	Zimbra ldap user dn. (default "uid=zimbra,cn=admins,cn=zimbra")
  -localconfig string
    	Where to find the localconfig.xml file. (default "/opt/zimbra/conf/localconfig.xml")
```
