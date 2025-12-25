package srv

import (
	"fmt"
	"net/http"
)

func (www *WWWServer) Handler(w http.ResponseWriter, r *http.Request) {

	// allowed origins?
	//r.Header

	//switch?
	fmt.Fprintln(w, "requested: ", r.RequestURI)

	www.logger.Info(r.RequestURI)
}
