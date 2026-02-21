package entity


type ImageKind string

const (
 ImageKindAvatar ImageKind= "avatar"
 ImageKindCover ImageKind="cover"
 ImageKindChapterPage ImageKind="chapter_page"
)

const (
	ImageKindAvatarStr = "avatar"
	ImageKindCoverStr= "cover"
	ImageKindChapterPageStr="chapter_page"
)

func (k ImageKind) String() string{
	return string(k)
}

func ValidImageKinds()[]ImageKind {
	return []ImageKind{
		ImageKindAvatar,ImageKindCover,ImageKindChapterPage,
	}
}
func IsValidImageKind(kind string) bool {
	for _, k := range ValidImageKinds() {
		if string(k) == kind {
			return true
		}
	}
	return false
}