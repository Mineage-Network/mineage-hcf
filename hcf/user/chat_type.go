package user

// ChatType represents different chats users may access.
type ChatType struct {
	chat int
}

// ChatTypeGlobal returns the global chat type.
func ChatTypeGlobal() ChatType {
	return ChatType{chat: 1}
}

// ChatTypeStaff returns the staff chat type.
func ChatTypeStaff() ChatType {
	return ChatType{chat: 2}
}

// ChatTypeFaction returns the faction chat type.
func ChatTypeFaction() ChatType {
	return ChatType{chat: 3}
}

// ChatTypes returns a list of all chat types.
func ChatTypes() []ChatType {
	return []ChatType{
		ChatTypeGlobal(),
		ChatTypeStaff(),
		ChatTypeFaction(),
	}
}

// String ...
func (d ChatType) String() string {
	switch d.chat {
	case 1:
		return "Global"
	case 2:
		return "Staff"
	case 3:
		return "Faction"
	}
	panic("should never happen")
}
