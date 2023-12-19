import time
import weatherhat
from weatherhat.history import WindSpeedHistory, WindDirectionHistory

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
    time.sleep(2.5)
