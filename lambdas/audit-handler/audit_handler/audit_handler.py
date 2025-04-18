from typing import Any

from audit_handler.database import AuditLog
from audit_handler import logger, app_config
from audit_handler.models import (
    AuditEvent,
    AuditResponse,
)


# Plyint disabled as callable class does not need public method
class AuditHandler:  # pylint: disable=too-few-public-methods
    """Handler class for audit lambda"""

    def __init__(self) -> None:
        self._audit_log = AuditLog(
            region=app_config.aws_region, tablename=app_config.dynamo_db
        )

    def __call__(self, event_message: Any, context: Any) -> AuditResponse:
        logger.info("Initializing audit handler...")
        logger.debug("Event [%s]", event_message)
        event: AuditEvent = AuditEvent.model_validate_json(
            event_message["Records"][0]["Sns"]["Message"]
        )
        logger.debug("Event extracted [%s]", event.model_dump())
        return self._handle_event(event)

    def _handle_event(self, event: AuditEvent) -> AuditResponse:
        response = AuditResponse(status="SUCCESS")
        self._audit_log.insert(event)
        logger.info("Responding to request with: %s", response)
        return response


handler = AuditHandler()
