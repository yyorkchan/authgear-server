package validation

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFormatPhone(t *testing.T) {
	f := FormatPhone{}.CheckFormat

	Convey("FormatPhone", t, func() {
		So(f(1), ShouldBeNil)
		So(f("+85298765432"), ShouldBeNil)
		So(f(""), ShouldBeError, "not in E.164 format")
		So(f("foobar"), ShouldBeError, "not in E.164 format")
	})
}

func TestFormatEmail(t *testing.T) {
	Convey("FormatEmail", t, func() {
		f := FormatEmail{}.CheckFormat

		So(f(1), ShouldBeNil)
		So(f("user@example.com"), ShouldBeNil)
		So(f(""), ShouldBeError, "mail: no address")
		So(f("John Doe <user@example.com>"), ShouldBeError, "input email must not have name")
		So(f("foobar"), ShouldBeError, "mail: missing '@' or angle-addr")
	})

	Convey("FormatEmail with name", t, func() {
		f := FormatEmail{AllowName: true}.CheckFormat

		So(f(1), ShouldBeNil)
		So(f("user@example.com"), ShouldBeNil)
		So(f("John Doe <user@example.com>"), ShouldBeNil)
		So(f(""), ShouldBeError, "mail: no address")
		So(f("foobar"), ShouldBeError, "mail: missing '@' or angle-addr")
	})
}

func TestFormatURI(t *testing.T) {
	Convey("FormatURI", t, func() {
		f := FormatURI{}.CheckFormat

		So(f(1), ShouldBeNil)
		So(f(""), ShouldBeError, "input URL must be absolute")
		So(f("/a"), ShouldBeError, "input URL must be absolute")
		So(f("#a"), ShouldBeError, "input URL must be absolute")
		So(f("?a"), ShouldBeError, "input URL must be absolute")

		So(f("http://example.com/../.."), ShouldBeError, "invalid path: /../..")

		So(f("http://example.com"), ShouldBeNil)
		So(f("http://example.com/"), ShouldBeNil)
		So(f("http://example.com/a"), ShouldBeNil)
		So(f("http://example.com/a/"), ShouldBeNil)

		So(f("http://example.com?"), ShouldBeNil)
		So(f("http://example.com?a"), ShouldBeNil)
		So(f("http://example.com?a=b"), ShouldBeNil)

		So(f("http://example.com#"), ShouldBeNil)
		So(f("http://example.com#a"), ShouldBeNil)
	})
}

func TestFormatPicture(t *testing.T) {
	Convey("FormatPicture", t, func() {
		f := FormatPicture{}.CheckFormat

		So(f(1), ShouldBeNil)
		So(f(""), ShouldBeError, "invalid scheme: ")
		So(f("foobar:"), ShouldBeError, "invalid scheme: foobar")

		So(f("http://example.com"), ShouldBeNil)
		So(f("http://example.com/"), ShouldBeNil)
		So(f("http://example.com/a"), ShouldBeNil)
		So(f("http://example.com/a/"), ShouldBeNil)

		So(f("http://example.com?"), ShouldBeNil)
		So(f("http://example.com?a"), ShouldBeNil)
		So(f("http://example.com?a=b"), ShouldBeNil)

		So(f("http://example.com#"), ShouldBeNil)
		So(f("http://example.com#a"), ShouldBeNil)

		So(f("authgearimages:///app/object"), ShouldBeNil)
		So(f("authgearimages:///../app/object"), ShouldBeError, "invalid path: /../app/object")
		So(f("authgearimages://host/"), ShouldBeError, "authgearimages URI does not have host")
	})
}

func TestFormatHTTPOrigin(t *testing.T) {
	Convey("FormatHTTPOrigin", t, func() {
		f := FormatHTTPOrigin{}.CheckFormat

		So(f(1), ShouldBeNil)
		So(f("http://example.com"), ShouldBeNil)
		So(f("http://example.com?"), ShouldBeNil)
		So(f("http://example.com#"), ShouldBeNil)

		So(f(""), ShouldBeError, "expect input URL with scheme http / https")
		So(f("http://user:password@example.com"), ShouldBeError, "expect input URL without user info, path, query and fragment")
		So(f("http://example.com/"), ShouldBeError, "expect input URL without user info, path, query and fragment")
		So(f("http://example.com?a"), ShouldBeError, "expect input URL without user info, path, query and fragment")
		So(f("http://example.com#a"), ShouldBeError, "expect input URL without user info, path, query and fragment")
	})
}

func TestFormatHTTPOriginSpec(t *testing.T) {
	Convey("FormatHTTPOriginSpec", t, func() {
		f := FormatHTTPOriginSpec{}.CheckFormat

		So(f(1), ShouldBeNil)
		So(f(""), ShouldBeNil)
		So(f("127.0.0.1"), ShouldBeNil)
		So(f("127.0.0.1/"), ShouldBeError, "127.0.0.1/ is not strict")
	})
}

func TestFormatLDAPURL(t *testing.T) {
	Convey("FormatLDAPURL", t, func() {
		f := FormatLDAPURL{}.CheckFormat

		So(f(1), ShouldBeNil)
		So(f("ldap://example.com"), ShouldBeNil)
		So(f("ldap://localhost:389"), ShouldBeNil)
		So(f("ldaps://example.com"), ShouldBeNil)

		So(f(""), ShouldBeError, "expect input URL with scheme ldap / ldaps")
		So(f("http://example.com"), ShouldBeError, "expect input URL with scheme ldap / ldaps")
		So(f("ldap://user:password@example.com"), ShouldBeError, "expect input URL without user info, path, query and fragment")
		So(f("ldap://example.com/"), ShouldBeError, "expect input URL without user info, path, query and fragment")
		So(f("ldap://example.com?a"), ShouldBeError, "expect input URL without user info, path, query and fragment")
		So(f("ldap://example.com#a"), ShouldBeError, "expect input URL without user info, path, query and fragment")
	})
}

func TestFormatLDAPDN(t *testing.T) {
	Convey("FormatLDAPDN", t, func() {
		f := FormatLDAPDN{}.CheckFormat

		So(f(1), ShouldBeNil)
		So(f("dc=example,dc=com"), ShouldBeNil)
		So(f("cn=admin,dc=example,dc=org"), ShouldBeNil)
		So(f("ou="), ShouldBeError, "invalid DN")
		So(f("asbbalskjedkbwk"), ShouldBeError, "invalid DN")
		So(f(""), ShouldBeError, "expect non-empty base DN")
	})
}

func TestFormatLDAPSearchFilterTemplate(t *testing.T) {
	Convey("FormatLDAPSearchFilterTemplate", t, func() {
		f := FormatLDAPSearchFilterTemplate{}.CheckFormat

		So(f(1), ShouldBeNil)
		So(f("(uid=%s)"), ShouldBeNil)
		So(f("((uid=%s)(cn=%s))"), ShouldBeError, "invalid search filter")
		So(f("(uid=%s)(cn=%s"), ShouldBeError, "invalid search filter")
		So(f(""), ShouldBeError, "invalid search filter")

		test_template := `{{if eq $.Username "test@test.com"}}(mail={{$.Username}}){{else if eq $.Username "+852"}}(telephoneNumber={{$.Username}}){{else}}(uid={{$.Username}}){{end}}`
		So(f(test_template), ShouldBeNil)

		test_template = `
		{{if eq .Username "test@test.com"}}
									(mail={{$.Username}})
		{{end}}

`
		So(f(test_template), ShouldBeNil)

		wrong_template := `{{if eq .Username "test@test.com"}}(mail={{.Username}})`
		So(f(wrong_template), ShouldBeError, "invalid template")
		wrong_filter := `{{if eq .Username "test@test.com"}}(mail={{.Username}}{{end}}`
		So(f(wrong_filter), ShouldBeError, "invalid search filter")
	})
}

func TestFormatLDAPOID(t *testing.T) {
	Convey("FormatLDAPOID", t, func() {
		f := FormatLDAPOID{}.CheckFormat

		So(f(1), ShouldBeNil)
		So(f("1.1"), ShouldBeNil)
		So(f("1.1.1"), ShouldBeNil)
		So(f("1"), ShouldBeNil)
		So(f(""), ShouldBeError, "expect non-empty OID")
		So(f(".1.1"), ShouldBeError, "invalid OID")
		So(f("1.1."), ShouldBeError, "invalid OID")
		So(f("1.1.1.a.2"), ShouldBeError, "invalid OID")
	})
}

func TestFormatWeChatAccountID(t *testing.T) {
	Convey("TestFormatWeChatAccountID", t, func() {
		f := FormatWeChatAccountID{}.CheckFormat

		So(f(1), ShouldBeNil)
		So(f("foobar"), ShouldBeError, "expect WeChat account id start with gh_")
		So(f("gh_foobar"), ShouldBeNil)
	})
}

func TestFormatBCP47(t *testing.T) {
	f := FormatBCP47{}.CheckFormat

	Convey("FormatBCP47", t, func() {
		So(f(1), ShouldBeNil)

		So(f(""), ShouldBeError, "invalid BCP 47 tag: language: tag is not well-formed")
		So(f("a"), ShouldBeError, "invalid BCP 47 tag: language: tag is not well-formed")
		So(f("foobar"), ShouldBeError, "invalid BCP 47 tag: language: tag is not well-formed")

		So(f("en"), ShouldBeNil)
		So(f("zh-TW"), ShouldBeNil)
		So(f("und"), ShouldBeNil)

		So(f("zh_TW"), ShouldBeError, "non-canonical BCP 47 tag: zh_TW != zh-TW")
	})
}

func TestFormatTimezone(t *testing.T) {
	f := FormatTimezone{}.CheckFormat

	Convey("FormatTimezone", t, func() {
		So(f(1), ShouldBeNil)

		So(f(""), ShouldBeError, `valid timezone name has at least 1 slash: ""`)
		So(f("UTC"), ShouldBeError, `valid timezone name has at least 1 slash: "UTC"`)

		So(f("Asia/Hong_Kong"), ShouldBeNil)
		// It seems that support for uppercase name is machine specific.
		// So(f("ASIA/HONG_KONG"), ShouldBeNil)
		So(f("Etc/UTC"), ShouldBeNil)
	})
}

func TestFormatBirthdate(t *testing.T) {
	f := FormatBirthdate{}.CheckFormat

	Convey("FormatBirthdate", t, func() {
		So(f(1), ShouldBeNil)

		So(f(""), ShouldBeError, `invalid birthdate: ""`)
		So(f("foobar"), ShouldBeError, `invalid birthdate: "foobar"`)

		So(f("2021"), ShouldBeNil)
		So(f("--01-01"), ShouldBeNil)
		So(f("0000-01-01"), ShouldBeNil)
		So(f("2021-01-01"), ShouldBeNil)
	})
}

func TestFormatAlpha2(t *testing.T) {
	f := FormatAlpha2{}.CheckFormat

	Convey("FormatAlpha2", t, func() {
		So(f(1), ShouldBeNil)

		So(f(""), ShouldBeError, `invalid ISO 3166-1 alpha-2 code: ""`)
		So(f("foobar"), ShouldBeError, `invalid ISO 3166-1 alpha-2 code: "foobar"`)

		So(f("US"), ShouldBeNil)
		So(f("HK"), ShouldBeNil)
		So(f("CN"), ShouldBeNil)
		So(f("TW"), ShouldBeNil)
		So(f("JP"), ShouldBeNil)
	})
}

func TestFormatCustomAttributePointer(t *testing.T) {
	f := FormatCustomAttributePointer{}.CheckFormat

	Convey("FormatCustomAttributePointer", t, func() {
		So(f(1), ShouldBeNil)

		So(f(""), ShouldBeError, "custom attribute pointer must be one-level but found 0")
		So(f("/"), ShouldBeError, "custom attribute pointer must not be empty")
		So(f("/0"), ShouldBeError, `invalid character at 0: "0"`)
		So(f("/_"), ShouldBeError, `invalid character at 0: "_"`)
		So(f("/a_"), ShouldBeError, `invalid character at 1: "_"`)
		So(f("/a-a"), ShouldBeError, `invalid character at 1: "-"`)
		So(f("/aあa"), ShouldBeError, `invalid character at 1: "あ"`)

		So(f("/a0"), ShouldBeNil)
		So(f("/aa"), ShouldBeNil)
		So(f("/aaa"), ShouldBeNil)
		So(f("/a_a"), ShouldBeNil)
		So(f("/mystring"), ShouldBeNil)
		So(f("/my_string"), ShouldBeNil)
	})
}

func TestFormatGoogleTagManagerContainerID(t *testing.T) {
	f := FormatGoogleTagManagerContainerID{}.CheckFormat

	Convey("FormatGoogleTagManagerContainerID", t, func() {
		So(f(1), ShouldBeNil)

		So(f(""), ShouldBeError, "expect google tag manager container ID to start with GTM-")
		So(f("GTM-AAAAAA"), ShouldBeNil)
	})
}

func TestFormatHookURI(t *testing.T) {
	Convey("FormatHookURI", t, func() {
		f := FormatHookURI{}.CheckFormat

		So(f(1), ShouldBeNil)
		So(f(""), ShouldBeError, "invalid scheme: ")
		So(f("foobar:"), ShouldBeError, "invalid scheme: foobar")

		So(f("http://example.com"), ShouldBeNil)
		So(f("http://example.com/"), ShouldBeNil)
		So(f("http://example.com/a"), ShouldBeNil)
		So(f("http://example.com/a/"), ShouldBeNil)

		So(f("http://example.com?"), ShouldBeNil)
		So(f("http://example.com?a"), ShouldBeNil)
		So(f("http://example.com?a=b"), ShouldBeNil)

		So(f("http://example.com#"), ShouldBeNil)
		So(f("http://example.com#a"), ShouldBeNil)

		So(f("authgeardeno:///deno/a.ts"), ShouldBeNil)
		So(f("authgeardeno:///../deno/a.ts"), ShouldBeError, "invalid path: /../deno/a.ts")
		So(f("authgeardeno://host/"), ShouldBeError, "authgeardeno URI does not have host")
	})
}

func TestFormatDurationString(t *testing.T) {
	f := FormatDurationString{}.CheckFormat

	Convey("FormatDurationString", t, func() {
		So(f(1), ShouldBeNil)

		So(f(""), ShouldBeError, `time: invalid duration ""`)
		So(f("foobar"), ShouldBeError, `time: invalid duration "foobar"`)

		So(f("0s"), ShouldBeError, `non-positive duration "0s"`)
		So(f("1.1s"), ShouldBeNil)
		So(f("2m3s"), ShouldBeNil)
		So(f("4h5m6s"), ShouldBeNil)
	})
}

func TestFormatBase64URL(t *testing.T) {
	f := FormatBase64URL{}.CheckFormat

	Convey("FormatBase64URL", t, func() {
		So(f(1), ShouldBeNil)
		So(f(""), ShouldBeNil)

		So(f("aA"), ShouldBeNil)
		So(f("aA=="), ShouldBeError, "illegal base64 data at input byte 2")
		So(f("a\nA"), ShouldBeNil)
	})
}

func TestFormatDateTime(t *testing.T) {
	f := FormatDateTime{}.CheckFormat

	Convey("TestFormatDateTime", t, func() {
		So(f(1), ShouldBeNil)
		So(f("2024-05-17T08:08:13.26635Z"), ShouldBeNil)
		So(f("2024-05-17T08:08:13Z"), ShouldBeNil)
		So(f("2024-05-17T08:08:13.26635+08:00"), ShouldBeNil)

		So(f(""), ShouldBeError, `date-time must be in rfc3999 format`)
	})
}
