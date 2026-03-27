package ports
 
type LinkResult struct {
    URL        string
    StatusCode int
    OK         bool
} 
type LinkChecker interface {
    CheckLinks(links []string) []LinkResult  
    
    Stop() // gracefully stops workers
}