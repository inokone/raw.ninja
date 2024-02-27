import pytest
import json

from unittest.mock import MagicMock

from rawninja.audit_handler import AuditHandler
from rawninja.models import AuditEvent
from rawninja.database import AuditLog


@pytest.fixture()
def mock_audit_event() -> AuditEvent:
    return AuditEvent(
        correlation_id="correlation_id",
        user_id="user_id",
        action="upload",
        target_id="target_id",
        target_type="target_type",
        meta=dict(key="value"),
        outcome="success",
    )


@pytest.fixture()
def mock_sns_event():
    return {
       "Records": [
            {
                "Sns": {
                    "Message": json.dumps(mock_audit_event()) 
                }
            }
        ]
    }


@pytest.fixture()
def mock_audit_log() -> AuditLog:
    return MagicMock()


def test_audit_handler_persisting_event():
    audit_log = mock_audit_log()
    handler = AuditHandler(mock_audit_log)

    handler(event_message=mock_sns_event())
            
    audit_log.insert.assert_called_once()
