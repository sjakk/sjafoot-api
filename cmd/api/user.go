package main


import(
	"net/http"

	"github.com/sjakk/sjafoot/internal/data"
)




func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
id, err := app.readIDParam(r)
if err != nil {
http.NotFound(w, r)
return
}
data := data.User{
ID: id,
}
err = app.writeJSON(w, http.StatusOK, data, nil)
if err != nil {
app.logger.Println(err)
http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
}
}

