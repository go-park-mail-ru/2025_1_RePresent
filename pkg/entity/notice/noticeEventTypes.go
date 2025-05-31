package notice

const (
	LowBalance     int = iota // if balance user too low
	TopUpedBalance            // if user top uped balance and his money became more then critical value
)
