package main

import "time"

// json format for message

type Message struct{
	Type 		string 	  `json:"type"`
	Username    string 	  `json:"username"`
	Body 	    string    `json:"body"`
	Timestamp   time.Time `json:"timestamp"`
}

