SELECT * FROM app.users WHERE last_modified > :sql_last_value AND last_modified < NOW()
