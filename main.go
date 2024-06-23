package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

//go:embed templates/* assets/*
var t embed.FS

type Data struct {
	Time              string
	TimeSince         time.Duration
	Temperature       string
	DeviceTemperature string
	Humidity          string
	RelativeHumidity  string
	Pressure          string
	Altitude          string
	Dewpoint          string
	WindSpeed         string
	Gusts             string
	WindDirection     string
	Rain              string
	Lux               string
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	r := gin.Default()
	templ := template.Must(template.New("").ParseFS(t, "templates/*.tmpl"))
	r.SetHTMLTemplate(templ)
	r.SetTrustedProxies(nil) // disable trusted proxies

	var contentFS, _ = fs.Sub(t, "assets")

	r.StaticFS("/assets", http.FS(contentFS))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html.tmpl", gin.H{
			"data": getSomeData(),
		})
	})

	r.Run(":8080")
}

func getSomeData() Data {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var timestamp time.Time
	var temperature, deviceTemperature, humidity,
		relativeHumidity, pressure, dewpoint, windSpeed,
		gusts, rain, waterTemperature, lux float32
	var windDirection string

	err = conn.QueryRow(context.Background(), "select * from data ORDER BY timestamp DESC limit 1").Scan(
		&timestamp,
		&temperature,
		&deviceTemperature,
		&humidity,
		&relativeHumidity,
		&pressure,
		&dewpoint,
		&windSpeed,
		&gusts,
		&windDirection,
		&rain,
		&waterTemperature,
		&lux,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(timestamp, temperature, relativeHumidity, windSpeed)

	d := &Data{}

	// format the time and get the time since
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	d.Temperature = fmt.Sprintf("%.2f", (temperature*9/5)+32)
	d.DeviceTemperature = fmt.Sprintf("%.2f", (deviceTemperature*9/5)+32)
	d.Humidity = fmt.Sprintf("%.2f", humidity)
	d.RelativeHumidity = fmt.Sprintf("%.2f", relativeHumidity)
	d.Dewpoint = fmt.Sprintf("%.2f", math.Ceil(float64(dewpoint)*100)/100)
	d.Pressure = fmt.Sprintf("%.2f", math.Ceil(float64(pressure)*100)/100)
	d.WindSpeed = fmt.Sprintf("%.2f", windSpeed)
	d.Gusts = fmt.Sprintf("%.2f", gusts)
	d.WindDirection = windDirection
	d.Rain = fmt.Sprintf("%.2f", rain)
	d.Lux = round(lux)
	d.Time = timestamp.In(loc).Format(time.UnixDate)
	d.TimeSince = time.Now().Sub(timestamp).Round(1 * time.Second)

	return *d
}

func round(val float32) string {
	return fmt.Sprintf("%.2f", val)
}
