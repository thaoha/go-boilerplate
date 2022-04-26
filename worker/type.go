package main

import "time"

type PostCreatedMsg struct {
	ID          string `json:"id"`
	ClasspageID string `json:"classpage_id"`
	UserID      string `json:"user_id"`
	Content     string `json:"content"`
}

type Classpage struct {
	ID         string `json:"id"`
	ClassName  string `json:"class_name"`
	SchoolName string `json:"school_name"`
}

type ClasspageMember struct {
	ClasspageID  string `json:"classpage_id"`
	UserID       string `json:"user_id"`
	UserFullname string `json:"user_fullname"`
	UserAvatar   string `json:"user_avatar"`
	Role         string `json:"role"`
}

type MeetingMember struct {
	ClasspageID  string `json:"classpage_id"`
	MeetingID    string `json:"meeting_id"`
	UserID       string `json:"user_id"`
	UserFullname string `json:"user_fullname"`
	UserAvatar   string `json:"user_avatar"`
	Role         string `json:"role"`
}

type ClasspageService interface {
	OnPostCreated(message PostCreatedMsg) error
}

type ClasspageRepository interface {
	Find(id string) (*Classpage, error)
	FindMember(classpageID string, userID string) (*ClasspageMember, error)
	GetMembers(classpageID string, page int, limit int) []ClasspageMember
	GetMeetingMembers(meetingID string, page int, limit int) []MeetingMember
	CreateNotifyCommand(userID string, title string, body string, data interface{})
}
