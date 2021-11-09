package main

import (
	"crypto/tls"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/ldap.v2"
)

//Default Attrs to migrate from zimbra Account
var accDefaultAttrs = "" +
	"displayName," +
	"zimbraAccountstatus," +
	"givenName," +
	"sn," +
	"zimbraIsAdminAccount," +
	"zimbraPrefMailForwardingAddress," +
	"zimbraPrefOutOfOfficeCacheDuration," +
	"zimbraPrefOutOfOfficeDirectAddress," +
	"zimbraPrefOutOfOfficeFromDate," +
	"zimbraPrefOutOfOfficeReply," +
	"zimbraPrefOutOfOfficeReplyEnabled," +
	"zimbraPrefOutOfOfficeUntilDate," +
	"zimbraPrefHtmlEditorDefaultFontColor," +
	"zimbraPrefHtmlEditorDefaultFontFamily," +
	"zimbraPrefHtmlEditorDefaultFontSize," +
	"zimbraPrefMessageViewHtmlPreferred," +
	"zimbraMailSieveScript," +
	"zimbraPrefComposeFormat," +
	"zimbraPrefGroupMailBy," +
	"zimbraSignatureName," +
	"zimbraSignatureId," +
	"zimbraPrefMailSignatureHTML," +
	"zimbraPrefMailSignature," +
	"zimbraPrefForwardReplySignatureId," +
	"zimbraPrefDefaultSignatureId," +
	"userPassword"

var attrs = flag.String("attrs", accDefaultAttrs, "Add a custom list of attributes separatted by comma.")
var domain = flag.String("domain", "*", "Domain to export. Leave blank will export all.")
var host = flag.String("ldap_host", "", "Zimbra LDAP Host to use.")
var port = flag.String("ldap_port", "389", "LDAP BIND PORT.")
var user = flag.String("ldap_user", "uid=zimbra,cn=admins,cn=zimbra", "Zimbra ldap user dn.")
var password = flag.String("ldap_password", "", "Password to authenticate in ldap server.")
var localconfig_xml = flag.String("localconfig", "/opt/zimbra/conf/localconfig.xml", "Where to find the localconfig.xml file.")

//Zimbra Localconfig structs
type Localconfig struct {
	XMLName     xml.Name `xml:"localconfig"`
	Localconfig []Key    `xml:"key"`
}

type Key struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value"`
}

func getLdapCredentials(host, port, user, password *string) (string, string, string, string) {
	var ldap_host, ldap_user, ldap_password, ldap_port string
	var localconfig Localconfig

	//Open zimbra localconfig
	xmlFile, err := os.Open(*localconfig_xml)
	//Ask for credentials if localconfig.xml not found
	if err != nil {
		log.Println("localconfig.xml not found, getting credentials from CLI")
		ldap_host = *host
		ldap_port = *port
		ldap_user = *user
		ldap_password = *password
		if ldap_host == "" {
			log.Println("please set ldap_host")
			flag.PrintDefaults()
			log.Fatal("Exitting... set ldap_host and try again")
		}
		if ldap_password == "" {
			log.Println("please set ldap_password")
			flag.PrintDefaults()
			log.Fatal("Exitting... set ldap_host and try again")
		}
	} else {
		// defer the closing of our xmlFile so that we can parse it later on
		defer xmlFile.Close()
		byteValue, _ := ioutil.ReadAll(xmlFile)

		xml.Unmarshal(byteValue, &localconfig)

		//Get the keys we need from zmlocalconfig
		log.Println("localconfig data loaded successfully. Ignoring ldap credentials from CLI.")
		for _, s := range localconfig.Localconfig {

			if s.Name == "ldap_host" {
				ldap_host = s.Value
			}
			if s.Name == "ldap_port" {
				ldap_port = s.Value
			}
			if s.Name == "zimbra_ldap_userdn" {
				ldap_user = s.Value
			}
			if s.Name == "zimbra_ldap_password" {
				ldap_password = s.Value
			}
		}
	}
	return ldap_host, ldap_port, ldap_user, ldap_password

}

func main() {
	//Init parser
	flag.Parse()
	// Init Credentials
	ldap_host, ldap_port, ldap_user, ldap_password := getLdapCredentials(host, port, user, password)
	//Split zimbra attrs into multiple fields separated by comma
	attrs := strings.Split(*attrs, ",")
	//Start LDAP communication
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", ldap_host, ldap_port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// Reconnect with TLS
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Fatal(err)
	}

	// First bind with a read only user
	err = l.Bind(ldap_user, ldap_password)
	if err != nil {
		log.Fatal(err)
	}

	//search for accounts
	zimbraAttrs := []string{"mail"}
	for _, item := range attrs {
		zimbraAttrs = append(zimbraAttrs, item)
	}
	searchRequest := ldap.NewSearchRequest(
		"",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=zimbraAccount)(mail=*@%s)(!(|(zimbraIsSystemAccount=TRUE)(objectClass=zimbraCalendarResource))))", *domain),
		zimbraAttrs,
		nil,
	)

	accResult, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}
	for _, acc := range accResult.Entries {
		mail := acc.GetAttributeValue("mail")
		fmt.Printf("createAccount %s '' ", mail)
		for _, item := range attrs {
			attr := acc.GetAttributeValue(item)
			if attr != "" {
				attr = strings.Replace(attr, "\n", "\\n", -1)
				attr = strings.Replace(attr, "'", "\\'", -1)
				fmt.Printf("%s '%s' ", item, attr)
			}
		}
		fmt.Printf("\n")
	}

	//search for lists
	zimbraAttrs = []string{"mail", "zimbraMailForwardingAddress"}
	searchRequest = ldap.NewSearchRequest(
		"",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=zimbraDistributionList)(mail=*@%s))", *domain),
		zimbraAttrs,
		nil,
	)
	listResult, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}
	for _, list := range listResult.Entries {
		mail := list.GetAttributeValue("mail")
		members := list.GetAttributeValues("zimbraMailForwardingAddress")
		fmt.Printf("createDistributionList %s\n", mail)
		for _, member := range members {
			fmt.Printf("addDistributionListMember %s %s\n", mail, member)
		}
	}

}
