package cmd

func cannot(err error) {
	if err != nil {
		panic(err)
	}
}
