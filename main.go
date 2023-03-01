package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"math"
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
	T           float64
	Humidity    float32
	H           float64
	Pressure    float32
	P           float64
	Altitude    float32
	A           float64
	Dewpoint    float32
	D           float64
}

type Wind struct {
	Timestamp time.Time
	Time      string
	TimeSince time.Duration
	Speed     float32
	S         float64
	Gusts     float32
	G         float64
	Direction float32
}

type Rain struct {
	Timestamp time.Time
	Time      string
	TimeSince time.Duration
	Amount    float32
	Total     float32
	T         float64
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
			"air":  getSomeAir(),
			"wind": getSomeWind(),
			"rain": getSomeRain(),
		})
	})

	r.GET("/history", func(c *gin.Context) {
		c.HTML(http.StatusOK, "history.html.tmpl", gin.H{})
	})

	r.Run(":8080")
}

func getSomeRain() Rain {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var rain Rain
	err = conn.QueryRow(context.Background(), `
WITH rain_cte AS (
SELECT 
  *   
FROM
  rain
WHERE timestamp BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW() 
ORDER BY timestamp DESC
)
SELECT SUM(amount) AS total
FROM rain_cte`).Scan(
		&rain.Total,
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

	rain.Time = rain.Timestamp.In(loc).Format(time.UnixDate)
	rain.TimeSince = time.Now().Sub(rain.Timestamp).Round(1 * time.Second)

	t := rain.Total / 25.4
	rain.T = math.Ceil(float64(t)*100) / 100

	return rain
}

func getSomeWind() Wind {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var wind Wind
	err = conn.QueryRow(context.Background(), "select * from Wind ORDER BY timestamp DESC limit 1").Scan(
		&wind.Timestamp,
		&wind.Speed,
		&wind.Gusts,
		&wind.Direction,
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

	wind.Time = wind.Timestamp.In(loc).Format(time.UnixDate)
	wind.TimeSince = time.Now().Sub(wind.Timestamp).Round(1 * time.Second)

	s := (wind.Speed / 1.609)
	wind.S = math.Ceil(float64(s)*100) / 100

	g := (wind.Gusts / 1.609)
	wind.G = math.Ceil(float64(g)*100) / 100

	return wind
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

	t := (air.Temperature * (9 / 5)) + 32
	air.T = math.Ceil(float64(t)*100) / 100

	air.H = math.Ceil(float64(air.Humidity)*100) / 100
	air.D = math.Ceil(float64(air.Dewpoint)*100) / 100
	air.P = math.Ceil(float64(air.Pressure)*100) / 100

	air.Time = air.Timestamp.In(loc).Format(time.UnixDate)
	air.TimeSince = time.Now().Sub(air.Timestamp).Round(1 * time.Second)

	return air
}
