package objects

// SingleObjectInfo — данные одного объекта (метка на карте) для ответа API.
type SingleObjectInfo struct {
	ID              int    `db:"id" json:"id"`
	Title           string `db:"title" json:"title"`
	Description     string `db:"description" json:"description"`
	EventDate       string `db:"event_date" json:"eventDate"`
	PreviewUrlImage string `db:"preview_image_url" json:"previewUrlImage"`
}

// ObjectInfo — объект в списке с координатами.
type ObjectInfo struct {
	ID              int    `db:"id" json:"id"`
	Title           string `db:"title" json:"title"`
	Description     string `db:"description" json:"description"`
	Latitude        string `db:"latitude" json:"latitude"`
	Longitude       string `db:"longitude" json:"longitude"`
	EventDate       string `db:"event_date" json:"eventDate"`
	EventTypeID     int    `db:"event_type_id" json:"eventType"`
	PreviewUrlImage string `db:"preview_image_url" json:"previewUrlImage"`
}

// EventTypeInfo — тип события.
type EventTypeInfo struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}
