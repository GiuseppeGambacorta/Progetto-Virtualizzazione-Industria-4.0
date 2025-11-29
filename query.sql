SELECT
  "time",
  CAST(value AS double precision) AS value
FROM mqtt_data
WHERE
  topic = 'home/scrivania/temperature' AND
  "time"::date = CURRENT_DATE
ORDER BY "time" DESC;


