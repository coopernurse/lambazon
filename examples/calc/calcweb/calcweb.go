package calcweb

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"html/template"
	"os"
	"strconv"
)

func NewRouterForLambda() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = os.Stderr
	return NewRouter()
}

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html", HomePage{}.Render())
	})

	router.POST("/sum", func(c *gin.Context) {
		a, _ := strconv.ParseFloat(c.PostForm("a"), 64)
		b, _ := strconv.ParseFloat(c.PostForm("b"), 64)
		home := HomePage{A: a, B: b, Sum: a + b}
		c.Data(200, "text/html", home.Render())
	})

	return router
}

type HomePage struct {
	A   float64
	B   float64
	Sum float64
}

func (h HomePage) Render() []byte {
	out := &bytes.Buffer{}
	t, err := template.New("home").Parse(`
<html>
<body>
<h1>calc</h1>

{{if .Sum}}
  <h2>Result</h2>
  <pre>{{.A}} + {{.B}} = {{.Sum}}</pre>
{{end}}

  <form method="POST" action="/sum">
    <input type="text" name="a" size="2" value="{{.A}}"> + <input type="text" name="b" size="2" value="{{.B}}">
    <input type="submit">
  </form>

</body>
</html>
`)
	if err != nil {
		panic(err)
	}
	err = t.Execute(out, h)
	if err != nil {
		panic(err)
	}
	return out.Bytes()
}
