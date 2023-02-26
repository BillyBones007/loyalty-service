package db

type Store interface {
	CreateTables() error
	Close()
	User() UniversalUserRepository
	Order() UniversalOrderRepository
}
