package engine

// 引擎启动入口
func run(symbol string) {
	// 判断是否已经开启
	if isOpen(symbol) {

	}
	// 开启某个交易标的撮合引擎
	go openMatching()
}

func isOpen(symbol string) bool {
	return true
}

func openMatching() {

}
