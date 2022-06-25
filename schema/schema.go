package schema

// album represents data about a record album.
type Persons struct {
    Personid     int64  `json:"personid"`
    Lastname  string  `json:"lastname"`
    Firstname string  `json:"firstname"`
    Address  string `json:"address"`
    City string `json:"city"`
}
