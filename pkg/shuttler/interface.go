package shuttler

type Shuttler interface {
	Shuttle(client, server Connection)
}
