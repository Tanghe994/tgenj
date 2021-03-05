package result

import "tgenj/database"


/*Result of a option.*/
type Result struct {
	Tx database.Transaction
	close bool
}
