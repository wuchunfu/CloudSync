# RealIP

Go package that can be used to get client's real public IP, which usually useful for logging HTTP server.

## Feature

- Follows the rule of X-Real-IP
- Follows the rule of X-Forwarded-For
- Exclude local or private address

## Example

```go
package main

import (
	"fmt"
	"github.com/wuchunfu/CloudSync/utils/ipx"
)

func (h *Handler) ServeIndexPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  clientIP := ipx.FromRequest(r)
  fmt.Println("GET / from", clientIP)
}
```
