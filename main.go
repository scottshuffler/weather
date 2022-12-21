package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

//go:embed templates/*
var t embed.FS

type Air struct {
	Timestamp   time.Time
	Time        string
	TimeSince   time.Duration
	Temperature float32
	Humidity    float32
	Pressure    float32
	Altitude    float32
	Dewpoint    float32
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	r := gin.Default()
	templ := template.Must(template.New("").ParseFS(t, "templates/*.tmpl"))
	r.SetHTMLTemplate(templ)
	r.SetTrustedProxies(nil) // disable trusted proxies
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html.tmpl", gin.H{
			"air": getSomeAir(),
		})
	})

	r.Run(":8080")
}

func getSomeAir() Air {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var air Air
	err = conn.QueryRow(context.Background(), "select * from Air ORDER BY timestamp DESC limit 1").Scan(
		&air.Timestamp,
		&air.Temperature,
		&air.Humidity,
		&air.Pressure,
		&air.Altitude,
		&air.Dewpoint,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	// format the time and get the time since
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	air.Time = air.Timestamp.In(loc).Format(time.UnixDate)
	air.TimeSince = time.Now().Sub(air.Timestamp).Round(1 * time.Second)

	return air
}
