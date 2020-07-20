package main

func CheckForError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
