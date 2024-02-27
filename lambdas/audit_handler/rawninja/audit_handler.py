from typing import Any

import json

from rawninja.database import AuditLog
from rawninja import logger, app_config
from rawninja.models import (
    AuditEvent,
    AuditResponse,
)


class AuditHandler:
    
    def __init__(self) -> None:
        self._audit_log = AuditLog(app_config.dynamo_db)

    def _insert_event(self, event: AuditEvent):
        try:
            self._audit_log.insert(event)
        except Exception as e:
            logger.error("Insert failed:", exc_info=True)
            raise Exception("Failed to insert audit event") from e

    def _handler(self, event: AuditEvent) -> AuditResponse:
        response = AuditResponse(
            status="SUCCESS"
        )
        self._insert_event(event)
        logger.info("Responding to request with: %s", response)
        return response

    def __call__(self, event_message: Any, context: Any) -> AuditResponse:
        event: AuditEvent = AuditEvent.parse_obj(
            json.loads(event_message["Records"][0]["Sns"]["Message"])
        )
        return self._handler(event)
