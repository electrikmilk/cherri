/*
 * Copyright (c) Cherri
 */

package main

type DocCategory string

const (
	Calendars DocCategory = "Calendars"
	Contacts  DocCategory = "Contacts"
	Reminders DocCategory = "Reminders"
	Notes     DocCategory = "Notes"
	Clock     DocCategory = "Clock"
	Camera    DocCategory = "Camera"
	Photos    DocCategory = "Photos"
	Podcasts  DocCategory = "Podcasts"
	Phone     DocCategory = "Phone"
	Maps      DocCategory = "Maps"
	Shortcuts DocCategory = "Shortcuts"
	Weather   DocCategory = "Weather"
	Safari    DocCategory = "Safari"
	Settings  DocCategory = "Settings"

	Stores        DocCategory = "Stores"
	QRCodes       DocCategory = "QR Codes"
	Scripting     DocCategory = "Scripting"
	Alarms        DocCategory = "Alarms"
	Dates         DocCategory = "Dates"
	Images        DocCategory = "Images"
	GIFs          DocCategory = "GIFs"
	Giphy         DocCategory = "Giphy"
	RSS           DocCategory = "RSS"
	Audio         DocCategory = "Audio"
	Video         DocCategory = "Video"
	Media         DocCategory = "Media"
	Music         DocCategory = "Music"
	Device        DocCategory = "Device"
	Disk          DocCategory = "Disk"
	Email         DocCategory = "Email"
	Files         DocCategory = "Files"
	Archives      DocCategory = "Archives"
	Documents     DocCategory = "Documents"
	Articles      DocCategory = "Articles"
	Books         DocCategory = "Books"
	Network       DocCategory = "Network"
	Previewing    DocCategory = "Previewing"
	Printing      DocCategory = "Printing"
	Sharing       DocCategory = "Sharing"
	Clipboard     DocCategory = "Clipboard"
	Messaging     DocCategory = "Messaging"
	TextEditing   DocCategory = "Text Editing"
	Translation   DocCategory = "Translation"
	Location      DocCategory = "Location"
	Travel        DocCategory = "Travel"
	Apps          DocCategory = "Apps"
	ControlFlow   DocCategory = "Control Flow"
	Dictionaries  DocCategory = "Dictionaries"
	Lists         DocCategory = "Lists"
	Math          DocCategory = "Math"
	Measurements  DocCategory = "Measurements"
	Noonce        DocCategory = "No-ops (noonce)"
	Notifications DocCategory = "Notifications"
	Dialogs       DocCategory = "Dialogs"
	Numbers       DocCategory = "Numbers"
	Shell         DocCategory = "Shell"
	Scripts       DocCategory = "Scripts"
	System        DocCategory = "System"
	Windows       DocCategory = "Windows"
	URLs          DocCategory = "URLs"
	Web           DocCategory = "Web"
	HTTP          DocCategory = "HTTP"
	XCallback     DocCategory = "XCallback"
	Crypto        DocCategory = "Crypto"
	BuiltIn       DocCategory = "Built-in"
)

var categories = map[DocCategory][]DocCategory{
	Calendars: {
		Reminders,
		Alarms,
		Clock,
		Dates,
	},
	Contacts: {
		Phone,
		Email,
	},
	Documents: {
		Archives,
		Books,
		Files,
		Notes,
		Previewing,
		Printing,
		QRCodes,
		TextEditing,
		Translation,
	},
	Location: {
		Maps,
		Travel,
		Weather,
	},
	Media: {
		Stores,
		Audio,
		Camera,
		GIFs,
		Images,
		Photos,
		Music,
		Podcasts,
		Video,
	},
	Scripting: {
		ControlFlow,
		Apps,
		Device,
		Disk,
		Dictionaries,
		Crypto,
		Lists,
		Math,
		Numbers,
		Measurements,
		Network,
		Notifications,
		Dialogs,
		Shell,
		Scripts,
		Shortcuts,
		System,
		Windows,
		URLs,
		XCallback,
		Noonce,
	},
	Sharing: {
		Clipboard,
		Messaging,
	},
	Web: {
		URLs,
		HTTP,
		Safari,
		Giphy,
		Articles,
		RSS,
	},
}
