package main

func NewClasspageService(classpages ClasspageRepository) ClasspageService {
	return &classpageService{classpages}
}

type classpageService struct {
	classpages ClasspageRepository
}

func (s classpageService) OnPostCreated(message PostCreatedMsg) error {
	classpage, err := s.classpages.Find(message.ClasspageID)
	if err != nil {
		return err
	}
	postOwner, err := s.classpages.FindMember(message.ClasspageID, message.UserID)
	if err != nil {
		return err
	}
	// notify to all members
	postContent := message.Content
	if len(postContent) > 100 {
		postContent = postContent[0:100] + "..."
	}
	if postContent == "" {
		postContent = "vừa đăng bài viết mới"
	}
	var (
		notiTitle = classpage.ClassName
		notiBody  = postOwner.UserFullname + ": " + postContent
		notiData  = map[string]string{
			"action":      "classpage-post-created",
			"object_type": "classpage",
			"object_id":   classpage.ID,
		}
	)
	var (
		members []ClasspageMember
		page    int = 1
		limit   int = 100
	)
	for {
		members = s.classpages.GetMembers(classpage.ID, page, limit)
		if len(members) <= 0 {
			break
		}
		for _, member := range members {
			if member.UserID == message.UserID {
				continue
			}
			s.classpages.CreateNotifyCommand(member.UserID, notiTitle, notiBody, notiData)
		}
		page++
	}
	return nil
}
