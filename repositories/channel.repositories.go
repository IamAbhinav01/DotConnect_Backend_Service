package repositories

type ChannelRepository interface {
	CreateChannel() (bool, error)
}