package alerts

// MaxLimit is the largest page size the alerts endpoint accepts. The
// endpoint is shared by every integration type, so it is the most alerts
// one call can return.
const MaxLimit = 1000
