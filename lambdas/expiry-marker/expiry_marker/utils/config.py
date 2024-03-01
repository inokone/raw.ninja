from pydantic import BaseSettings
from rawninja.utils.kms import KMSSecret


class AppConfig(BaseSettings):
    LOG_LEVEL: str
    VERSION: str

    # image bucket
    AWS_REGION: str
    TARGET_BUCKET: str

    # database parameters
    POSTGRES_USER: str
    POSTGRES_PASSWORD: KMSSecret
    POSTGRES_HOST: str
    POSTGRES_DBNAME: str

    def get_connection_string(self):
        return f"postgresql+psycopg2://{self.POSTGRES_USER}:{self.POSTGRES_PASSWORD.get_secret_value()}@{self.POSTGRES_HOST}/{self.POSTGRES_DBNAME}"


