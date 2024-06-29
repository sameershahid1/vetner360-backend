package static_data

var IndexCollection = map[string][]string{
	"users":           {"token", "location", "firstName", "lastName"},
	"token":           {"token"},
	"contactmessages": {"token"},
	"pets":            {"token"},
	"activity":        {"token"},
	"participants":    {"token", "room"},
	"messages":        {"token", "room"},
}

var IndexCollectionAttribute = map[string]map[string]interface{}{
	"users":           {"token": 1, "location": "2dsphere", "firstName": "text", "lastName": "text"},
	"token":           {"token": 1},
	"contactmessages": {"token": 1},
	"pets":            {"token": 1},
	"activity":        {"token": 1},
	"participants":    {"token": 1, "room": 1},
	"messages":        {"token": 1, "room": 1},
}
