package strategy

// ------ balance strategy interface ------

type IBalanceStrategy interface {
	Next() int
}
