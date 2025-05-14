package models

import "time"

type UserData struct {
	Login       string
	Email       string
	Encrypted   int
	Transcribed int
}

type Counters struct {
	UserEncrypted   int
	UserTranscribed int
	AllUsersActions int
}

type ProfileCreditionals struct {
	Email       string
	Login       string
	Encrypted   int
	Transcribed int
	Createdat   time.Time
}
