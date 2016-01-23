package blogalert

// Flag defines article flag
type Flag int

// Article Flags
const (
	Show   Flag = iota // Show article in listing and index
	Hide               // Hide page in listing, but index
	Ignore             // Hide page in listing and dont index
)
