import os
import subprocess
import time
import utils.config as config


def start_containers():
    import yaml

    with open("../configs/database.yml", "r") as f:
        database_config = yaml.safe_load(f)

    with open(".env", "w") as f:
        f.write(f'POSTGRES_USER={database_config["database"]["POSTGRES_USER"]}\n')
        f.write(
            f'POSTGRES_PASSWORD={database_config["database"]["POSTGRES_PASSWORD"]}\n'
        )
        f.write(f'POSTGRES_DB={database_config["database"]["POSTGRES_DB"]}\n')
        f.write(f'HOST={database_config["database"]["HOST"]}\n')
        f.write(f'PORT={database_config["database"]["PORT"]}\n')
        f.write(f'INSIDE_PORT={database_config["database"]["INSIDE_PORT"]}\n')

    subprocess.run(["docker compose", "up", "-d"])
    print("Containers starting...")
    time.sleep(1)


def init_database_tables(drop_tables=False):
    import psycopg2
    import yaml

    with open("../configs/database.yml", "r") as f:
        database_config = yaml.safe_load(f)
    conn = psycopg2.connect(
        host="localhost",
        port=database_config["database"]["PORT"],
        database=database_config["database"]["POSTGRES_DB"],
        user=database_config["database"]["POSTGRES_USER"],
        password=database_config["database"]["POSTGRES_PASSWORD"],
        sslmode=database_config["database"]["SSLMODE"],
    )
    cur = conn.cursor()
    try:
        if drop_tables:
            with open("./requests/dropDataBase.sql", "r") as f:
                sql_file = f.read()

                cur.execute(sql_file)
                conn.commit()

        with open("./requests/createDataBase.sql", "r") as f:
            sql_file = f.read()

            cur.execute(sql_file)
            conn.commit()
        print("SQL Запросы выполнены успешно!")
    except psycopg2.Error as e:
        print(f"Ошибка при подключении к базе данных: {e}")
    except Exception as e:
        print(f"Ошибка в SQL Запросе: {e}")
    finally:
        conn.close()


if __name__ == "__main__":
    # config.configurate()

    start_containers()
    #if input("Init tables? Yes/No : ").lower() == "yes":
    #    init_database_tables(
    #        input("Drop tables if its exists? Yes/No : ").lower() == "yes"
    #    )
