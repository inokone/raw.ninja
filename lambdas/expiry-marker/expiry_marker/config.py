from pydantic_settings import BaseSettings

from expiry_marker.kms import KMSSecret


class AppConfig(BaseSettings):
    log_level: str

    # image bucket
    aws_region: str
    target_bucket: str

    # database parameters
    db_user: str
    db_password: KMSSecret
    db_host: str
    db_name: str

    def get_connection_string(self):
        return f"postgresql+psycopg2://{self.db_user}:
            {self.db_password.get_secret_value()}@{self.db_host}/{self.db_name}"
