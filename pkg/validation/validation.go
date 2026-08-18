package validation

import "regexp"

// EmailCheck matches the iCRM email pattern, e.g. test@google.com
const EmailCheck = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9\-]+\.[a-zA-Z]{2,}(?:\.[a-zA-Z]{2,})?(?:\.[a-zA-Z]{2,})?$`

// ObjectIDCheck matches a MongoDB ObjectID hex string (same length/charset as bson.ObjectIDFromHex).
const ObjectIDCheck = `^[0-9a-fA-F]{24}$`

var (
	emailPattern    = regexp.MustCompile(EmailCheck)
	objectIDPattern = regexp.MustCompile(ObjectIDCheck)
)

func IsValidEmail(value string) bool {
	return emailPattern.MatchString(value)
}

func IsValidObjectID(value string) bool {
	return objectIDPattern.MatchString(value)
}
