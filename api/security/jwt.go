package security

// This struct is an abstraction over the different
// libraries used for offline (dev/test) and online
// jwt validation.
// Missing properties can be added on demand

type JWT struct {
	Subject string
}
