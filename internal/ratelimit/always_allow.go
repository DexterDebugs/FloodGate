package ratelimit

type AlwaysAllow struct{}

func (a *AlwaysAllow) Allow(clientID, route string) bool{
	return true
}

func (aa *AlwaysAllow) SetLimit(limit int){
	
}