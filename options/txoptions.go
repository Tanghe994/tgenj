package options


/*TxOptions is used to configure a transaction upon creation*/
type TxOptions struct {
	ReadOnly bool
	Writable bool
}
