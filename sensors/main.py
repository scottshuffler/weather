import time
import weatherhat
import psycopg2
import datetime
import os

from pytz import timezone
from weatherhat.history import WindSpeedHistory, WindDirectionHistory
from dotenv import load_dotenv

load_dotenv()

sensor = weatherhat.WeatherHAT()
wind_speed_history = WindSpeedHistory()
wind_direction = WindDirectionHistory()

while True:
    sensor.update(interval=5.0)

    if sensor.updated_wind_rain:
        wind_speed_history.append(sensor.wind_speed)
        wind_direction.append(sensor.wind_direction)

        print(f"Average wind speed: {wind_speed_history.average_mph()}mph")
        print(f"Wind gust: {wind_speed_history.gust_mph()}mph")
        print(f"wind_direction: {sensor.wind_direction}")
        print(f"rain: {sensor.rain}")
        print(f"rain_total: {sensor.rain_total}")

        wind_direction_cardinal = wind_direction.average_compass(60)
        print(f"wind_direction: {wind_direction_cardinal}")
        
        print(f"temp: {sensor.temperature}")
        print(f"humidity: {sensor.humidity}")
        print(f"relative_humidity: {sensor.relative_humidity}")
        print(f"pressure: {sensor.pressure}")
        print(f"dewpoint: {sensor.dewpoint}")
        print(f"lux: {sensor.lux}")

        conn = psycopg2.connect(
                database=os.getenv('DATABASE'),
                user=os.getenv('DBUSER'),
                password=os.getenv('PASSWORD'),
                host=os.getenv('HOST'),
                port=os.getenv('PORT'),
            )

        now_time = datetime.datetime.now(timezone("US/Eastern"))
        print(f"timestamp: {now_time}")
        cur = conn.cursor()
        cur.execute(
            "INSERT INTO data (timestamp, temperature, device_temperature, humidity, relative_humidity, pressure, dewpoint, wind_speed, gusts, wind_direction, rain, water_temperature, lux) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)",
            (
                now_time,
                sensor.temperature,
                sensor.device_temperature,
                sensor.humidity,
                sensor.relative_humidity,
                sensor.pressure,
                sensor.dewpoint,
                wind_speed_history.average_mph(),
                wind_speed_history.gust_mph(),
                wind_direction_cardinal,
                sensor.rain_total,
                0,
                sensor.lux,
            ),
        )
        conn.commit()

        cur.close()
        conn.close()
    time.sleep(6)
