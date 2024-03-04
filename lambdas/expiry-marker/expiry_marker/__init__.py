import logging
from expiry_marker.config import AppConfig
from expiry_marker.database import create_session_maker
from sqlalchemy import sessionmaker, Session

SUPPRESSED_LOGS = (
    "boto",
    "boto3",
    "botocore",
    "urllib3.connectionpool",
    "aiobotocore",
    "uvicorn.access",
)

app_config = AppConfig()  # type: ignore  # parameters filled from ENV

for category in SUPPRESSED_LOGS:
    logging.getLogger(category).setLevel(logging.WARNING)

logger: logging.Logger = logging.getLogger("expiry_marker")
logger.setLevel(app_config.log_level)

session: sessionmaker[Session] = create_session_maker(app_config.get_connection_string())