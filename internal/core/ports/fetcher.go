package ports

type Fetcher interface {
    Fetch(url string) (body string, statusCode int, err error)
}