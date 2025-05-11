package notice

type NoticeRepositoryInterface struct {
}

type NoticeRepository struct {
}

func NewNoticeRepository(dsn string) *NoticeRepository {
	noticeRepo := &NoticeRepository{}
	return noticeRepo
}
