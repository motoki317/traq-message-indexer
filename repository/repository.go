package repository

type Repository interface {
	// チャンネル更新用のロック
	ChannelLock()
	ChannelUnlock()
	SeenChannelRepository
	MessageRepository
}
