import utils.config as config

config.configurate()

import yaml

with open("../configs/database.yml", "r") as f:
    database_config = yaml.safe_load(f)

with open(".env", "w") as f:
    f.write(f'POSTGRES_USER={database_config["database"]["POSTGRES_USER"]}\n')
    f.write(f'POSTGRES_PASSWORD={database_config["database"]["POSTGRES_PASSWORD"]}\n')
    f.write(f'POSTGRES_DB={database_config["database"]["POSTGRES_DB"]}\n')
    f.write(f'HOST={database_config["database"]["HOST"]}\n')
    f.write(f'PORT={database_config["database"]["PORT"]}\n')
    f.write(f'INSIDE_PORT={database_config["database"]["INSIDE_PORT"]}\n')
