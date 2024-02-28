from typing import Any

import json

from rawninja.database import AuditLog
from rawninja import logger, app_config
from rawninja.models import (
    AuditEvent,
    AuditResponse,
)


class AuditPersistenceError(Exception):
    pass

# Plyint disabled as callable class does not need public method
class AuditHandler: # pylint: disable=too-few-public-methods
    """ Handler class for audit lambda """

    def __init__(self) -> None:
        self._audit_log = AuditLog(app_config.dynamo_db)

    def _insert_event(self, event: AuditEvent):
        try:
            self._audit_log.insert(event)
        except Exception as e:
            logger.error("Insert failed:", exc_info=True)
            raise AuditPersistenceError("Failed to insert audit event") from e

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
