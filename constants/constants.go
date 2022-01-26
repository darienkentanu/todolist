package constants

import "time"

const PORT = ":8080"

var EXPIRATION_TIME = time.Now().Add(24 * time.Hour)
