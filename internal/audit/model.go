package audit

const (
	ActionInsert = "INSERT"
	ActionUpdate = "UPDATE"
	ActionDelete = "DELETE"
)

type Entry struct {
	TableName     string
	RecordID      int64
	Action        string
	OldData       any
	NewData       any
	ChangedFields []string
	CreatedBy     *int64
}
