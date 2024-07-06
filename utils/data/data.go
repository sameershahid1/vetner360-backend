package static_data

var IndexCollection = map[string][]string{
	"users":           {"token", "location", "firstName", "lastName", "email", "clinicName", "experience"},
	"contactmessages": {"token"},
	"pets":            {"token", "name", "nickName", "type", "gender"},
	"petsells":        {"token"},
	"activity":        {"token"},
	"participants":    {"token", "room"},
	"messages":        {"token", "room"},
	"reviews":         {"token"},
}

var IndexCollectionAttribute = map[string]map[string]interface{}{
	"users":           {"token": 1, "location": "2dsphere", "firstName": "text", "lastName": "text", "email": "text", "clinicName": "text", "experience": "text"},
	"contactmessages": {"token": 1},
	"pets":            {"token": 1, "name": "text", "nickName": "text", "type": "text", "gender": "text"},
	"petsells":        {"token": 1},
	"activity":        {"token": 1},
	"participants":    {"token": 1, "room": 1},
	"messages":        {"token": 1, "room": 1},
	"reviews":         {"token": 1},
}
