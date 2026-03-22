package errs

import "fmt"

var ErrDiscNotFound error = fmt.Errorf("Discipline with that id not found")
var ErrDocNotFound error = fmt.Errorf("Document with that id not found")
