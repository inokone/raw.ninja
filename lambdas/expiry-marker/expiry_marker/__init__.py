import logging
from expiry_marker.config import AppConfig


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

logger = logging.getLogger("expiry_marker")
logger.setLevel(app_config.log_level)
