from pydantic_settings import BaseSettings


class AppConfig(BaseSettings):
    """ Data type for application configuration of qudit handler """

    log_level: str
    version: str
    aws_region: str
    dynamo_db: str
