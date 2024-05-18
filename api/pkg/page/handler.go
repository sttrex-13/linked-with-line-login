package page

import "net/http"

type PageHandler struct {
}

func (p *PageHandler) LoginSuccess(w http.ResponseWriter, r *http.Request) {
	html := `
		<!DOCTYPE html>
		<html>
		<head>
		    <title>Login Successfully</title>
		</head>
		<body>
		    <h1>You are logged in!!</h1>
		</body>
		</html>
    `
	w.Write([]byte(html))
}
