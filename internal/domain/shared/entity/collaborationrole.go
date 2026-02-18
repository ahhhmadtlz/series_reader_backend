package entity

type CollaborationRole uint8

const (
	ViewerCollabRole CollaborationRole = iota + 1
	EditorCollabRole
	PublisherCollabRole
	OwnerCollabRole
)

const (
	ViewerCollabRoleStr    = "viewer"
	EditorCollabRoleStr    = "editor"
	PublisherCollabRoleStr = "publisher"
	OwnerCollabRoleStr     = "owner"
)

func (c CollaborationRole) String() string {
	switch c {
	case ViewerCollabRole:
		return ViewerCollabRoleStr
	case EditorCollabRole:
		return EditorCollabRoleStr
	case PublisherCollabRole:
		return PublisherCollabRoleStr
	case OwnerCollabRole:
		return OwnerCollabRoleStr
	}
	return ""
}

func MapToCollaborationRole(roleStr string) CollaborationRole {
	switch roleStr {
	case ViewerCollabRoleStr:
		return ViewerCollabRole
	case EditorCollabRoleStr:
		return EditorCollabRole
	case PublisherCollabRoleStr:
		return PublisherCollabRole
	case OwnerCollabRoleStr:
		return OwnerCollabRole
	}
	return CollaborationRole(0)
}