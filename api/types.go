package api

// InventoryItem type
type InventoryItem struct {
	Name                          string `json:"name"`
	ItemType                      string `json:"itemType"`
	ApplicationType               string `json:"applicationType"`
	SchoolClassID                 string `json:"schoolClassId"`
	SchoolClassName               string `json:"schoolClassName"`
	SchoolID                      string `json:"schoolId"`
	SchoolName                    string `json:"schoolName"`
	SchoolRegion                  string `json:"schoolRegion"`
	SchoolYearStart               int    `json:"schoolYearStart"`
	PurchasedFeature              string `json:"purchasedFeature"`
	PurchasedUntil                string `json:"purchasedUntil"`
	HasUnreadMessages             bool   `json:"hasUnreadMessages"`
	HasUnreadDiscussions          bool   `json:"hasUnreadDiscussions"`
	MigrationPending              bool   `json:"migrationPending"`
	SchoolClassPictureID          string `json:"schoolClassPictureId"`
	IsPrincipalMessagingActivated *bool  `json:"isPrincipalMessagingActivated"`
	TeacherRole                   string `json:"teacherRole"`
	CanCreateClasses              *bool  `json:"canCreateClasses"`
	IsFoxAdmin                    *bool  `json:"isFoxAdmin"`
	IsSchoolValid                 bool   `json:"isSchoolValid"`
	IsConnectedToPrincipal        bool   `json:"isConnectedToPrincipal"`
	ColorCode                     string `json:"colorCode"`
	HasTeamClass                  *bool  `json:"hasTeamClass"`
	IsTeamClass                   *bool  `json:"isTeamClass"`
	CreatedBy                     string `json:"createdBy"`
	UpdatedBy                     string `json:"updatedBy"`
	IsActive                      bool   `json:"isActive"`
	ID                            string `json:"id"`
	Version                       string `json:"version"`
	CreatedAt                     string `json:"createdAt"`
	UpdatedAt                     string `json:"updatedAt"`
	Deleted                       bool   `json:"deleted"`
}

// FDItem type
type FDItem struct {
	Name                 string  `json:"name"`
	FullPath             string  `json:"fullPath"`
	CreatorName          string  `json:"creatorName"`
	ItemType             string  `json:"itemType"`
	ItemSubType          string  `json:"itemSubType"`
	TeachersAccessType   string  `json:"teachersAccessType"`
	ParentsAccessType    string  `json:"parentsAccessType"`
	NumberOfParticipants int64   `json:"numberOfParticipants"`
	HasPreview           bool    `json:"hasPreview"`
	Size                 int64   `json:"size"`
	LastEditedDate       string  `json:"lastEditedDate"`
	ParentItemID         *string `json:"parentItemId"`
	SchoolClassID        string  `json:"schoolClassId"`
	PupilID              string  `json:"pupilId"`
	AccessType           string  `json:"accessType"`
	ID                   string  `json:"id"`
	CreatedAt            string  `json:"createdAt"`
	CreatedBy            string  `json:"createdBy"`
	UpdatedAt            string  `json:"updatedAt"`
	UpdatedBy            string  `json:"updatedBy"`
	Deleted              bool    `json:"deleted"`
	Version              string  `json:"version"`
	IsActive             bool    `json:"isActive"`
}
