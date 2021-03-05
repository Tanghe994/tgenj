package transaction



/*TODO*/
type Transaction interface {

	Rollback() error


	Commit() error

}
