import logging
from typing import Optional, List

DEFAULT_SUPPRESSED_LOGS = [
    "boto",
    "boto3",
    "botocore",
    "urllib3.connectionpool",
    "aiobotocore",
    "uvicorn.access"
]

def initialize_logger(
    log_level: str,
    log_categories_to_suppress: Optional[List[str]] = DEFAULT_SUPPRESSED_LOGS,
) -> None:
    logging.basicConfig(level=log_level)

    for category in log_categories_to_suppress:
        logging.getLogger(category).setLevel(logging.WARNING)
