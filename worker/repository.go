package main

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/aws/aws-sdk-go-v2/service/chime"
	"gorm.io/gorm"
)

const (
	SendNotificationTopic = "notifications.cmd.send.0"
)

type auroraClasspageRepository struct {
	aurora        *gorm.DB
	awsChime      *chime.Client
	kafkaProducer sarama.SyncProducer
}

func NewClasspageRepository(aurora *gorm.DB, awsChime *chime.Client, kafkaProducer sarama.SyncProducer) ClasspageRepository {
	classpageRepository := &auroraClasspageRepository{aurora, awsChime, kafkaProducer}
	return classpageRepository
}

func (r *auroraClasspageRepository) Find(id string) (*Classpage, error) {
	var classpage Classpage
	result := r.aurora.Table("classpage_classpages").First(&classpage, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &classpage, nil
}

func (r *auroraClasspageRepository) FindMember(classpageID string, userID string) (*ClasspageMember, error) {
	var member ClasspageMember
	result := r.aurora.Table("classpage_members").
		Where("classpage_id = ?", classpageID).
		Where("user_id = ?", userID).
		First(&member)

	if result.Error != nil {
		return nil, result.Error
	}
	// get user info
	var userProfile struct {
		Fullname string
		Avatar   string
	}
	r.aurora.Table("user_profiles").Where("user_id = ?", userID).First(&userProfile)
	member.UserFullname = userProfile.Fullname
	member.UserAvatar = userProfile.Avatar
	return &member, nil
}

func (r *auroraClasspageRepository) GetMembers(classpageID string, page int, limit int) []ClasspageMember {
	var members []ClasspageMember
	r.aurora.Table("classpage_members").
		Where("classpage_id = ?", classpageID).
		Limit(limit).
		Offset(limit * (page - 1)).
		Find(&members)

	return members
}

func (r *auroraClasspageRepository) GetMeetingMembers(meetingID string, page int, limit int) []MeetingMember {
	var members []MeetingMember
	r.aurora.Table("classpage_meeting_members").
		Where("meeting_id = ?", meetingID).
		Limit(limit).
		Offset(limit * (page - 1)).
		Find(&members)

	return members
}

func (r *auroraClasspageRepository) CreateNotifyCommand(userID string, title string, body string, data interface{}) {
	notification := map[string]interface{}{
		"user_id": userID,
		"title":   title,
		"body":    body,
		"data":    data,
	}
	messageValue, _ := json.Marshal(notification)
	r.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Value: sarama.ByteEncoder(messageValue),
		Topic: SendNotificationTopic,
	})
}
