import logging
from rawninja.config import AppConfig


SUPPRESSED_LOGS = (
    "boto",
    "boto3",
    "botocore",
    "urllib3.connectionpool",
    "aiobotocore",
    "uvicorn.access"
)

app_config = AppConfig()

logging.basicConfig(level=app_config.log_level)
for category in SUPPRESSED_LOGS:
    logging.getLogger(category).setLevel(logging.WARNING)

logger = logging.getLogger("audit_handler")
