import os
import subprocess
import time
import psycopg2
import yaml
import utils.config as utils
import utils.secure as secure


def update_database_config():
    with open("../configs/database.yml", "r") as f:
        database_config = yaml.safe_load(f)

    password = database_config["database"]["POSTGRES_PASSWORD"]
    hashed_password = secure.hash_password(password)

    database_config["database"]["HASH_PASSWORD"] = hashed_password
    with open("../configs/database.yml", "w") as f:
        yaml.dump(database_config, f, default_flow_style=False)
        
    


def start_containers():
    with open("../configs/database.yml", "r") as f:
        database_config = yaml.safe_load(f)

    with open(".env", "w") as f:
        f.write(f'POSTGRES_USER={database_config["database"]["POSTGRES_USER"]}\n')
        f.write(
            f'POSTGRES_PASSWORD={database_config["database"]["POSTGRES_PASSWORD"]}\n'
        )
        f.write(f'POSTGRES_DB={database_config["database"]["POSTGRES_DB"]}\n')
        f.write(f'PORT={database_config["database"]["PORT"]}\n')
        f.write(f'INSIDE_PORT={database_config["database"]["INSIDE_PORT"]}\n')

    subprocess.run(["docker-compose", "up", "-d"])
    print("Containers starting...")
    time.sleep(10)


def init_database_tables(drop_tables=False):
    with open("../configs/database.yml", "r") as f:
        database_config = yaml.safe_load(f)

    conn = psycopg2.connect(
        host="localhost",
        port=int(database_config["database"]["PORT"]),
        database=database_config["database"]["POSTGRES_DB"],
        user=database_config["database"]["POSTGRES_USER"],
        password=database_config["database"]["HASH_PASSWORD"],
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
    except e:
        print(f"Ошибка в SQL Запросе: {e}")
    finally:
        conn.close()


if __name__ == "__main__":
    utils.install_pip()
    utils.install_pip_tools()
    utils.generate_requirements()
    utils.install_requirements()

    update_database_config()
    
    start_containers()
    is_drop_tables = input("Drop tables if they exists? Yes/No : ").lower() == "yes"
    init_database_tables(is_drop_tables)
