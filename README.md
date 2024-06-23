# weather

```bash

just g

just update-go-repos

just dev

fly deploy

```

                                                       Table "public.air"
   Column    |            Type             | Collation | Nullable | Default | Storage | Compression | Stats target | Description
-------------+-----------------------------+-----------+----------+---------+---------+-------------+--------------+-------------
 timestamp   | timestamp without time zone |           |          |         | plain   |             |              |
 temperature | real                        |           |          |         | plain   |             |              |
 humidity    | real                        |           |          |         | plain   |             |              |
 pressure    | real                        |           |          |         | plain   |             |              |
 altitude    | real                        |           |          |         | plain   |             |              |
 dewpoint    | real                        |           |          |         | plain   |             |              |


CREATE TABLE data (
   timestamp timestamp without time zone PRIMARY KEY,
   temperature real,
   device_temperature real,
   humidity real,
   relative_humidity real,
   pressure real,
   dewpoint real,
   wind_speed real,
   gusts real,
   wind_direction VARCHAR(255),
   rain real,
   water_temperature real,
   lux real
);

CREATE TABLE [IF NOT EXISTS] air (
   timestamp timestamp without time zone PRIMARY KEY,
   temperature real,
   device_temperature real,
   humidity real,
   relative_humidity real,
   pressure real,
   altitude real,
   dewpoint real,
);

CREATE TABLE wind (
   timestamp timestamp without time zone PRIMARY KEY,
   speed real,
   gusts real,
   direction real
);

CREATE TABLE rain (
   timestamp timestamp without time zone PRIMARY KEY,
   amount real
);

CREATE TABLE [IF NOT EXISTS] water (
   timestamp timestamp without time zone PRIMARY KEY,
   temperature real
);

CREATE TABLE accounts (
	user_id serial PRIMARY KEY,
	username VARCHAR ( 50 ) UNIQUE NOT NULL,
	password VARCHAR ( 50 ) NOT NULL,
	email VARCHAR ( 255 ) UNIQUE NOT NULL,
	created_on TIMESTAMP NOT NULL,
        last_login TIMESTAMP 
);