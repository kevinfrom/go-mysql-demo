package main

func CheckError(err any) {
	if err != nil {
		panic(err)
	}
}
