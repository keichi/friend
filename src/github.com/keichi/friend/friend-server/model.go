package main

import (
	"time"
)

type User struct {
	Id        int64 `primaryKey:"yes"`
	Name      string
	Password  string
	PublicKey string
	Sessions  []Session
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TrustRelation struct {
	Id        int64 `primaryKey:"yes"`
	TrusterId int64
	TrusteeId int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	Id        int64 `primaryKey:"yes"`
	UserId    int64
	Token     string
	Expires   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transfer struct {
	Id            int64 `primaryKey:"yes"`
	Token         string
	Sender        User
	Receiver      User
	FileName      string
	FileSize      int
	SenderReady   bool
	ReceiverReady bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
