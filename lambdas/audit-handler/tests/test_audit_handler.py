import pytest

from unittest.mock import MagicMock
from unittest.mock import patch

from dotenv import load_dotenv

load_dotenv("tests/test.env", override=True)

from audit_handler.audit_handler import AuditHandler


MOCK_MESSAGE = """{
    "correlation_id": "correlation_id_1",
    "user_id": "user_id_1",
    "action": "upload",
    "target_ids": [
        "photo_id_1",
        "photo_id_2"
    ],
    "target_type": "photo",
    "meta": {
        "key": "value"
    },
    "entry_date": 1709415102,
    "outcome": "success"
}"""


@pytest.fixture()
def mock_sns_event():
    yield {"Records": [{"Sns": {"Message": MOCK_MESSAGE}}]}


def test_audit_handler_persisting_event(mock_sns_event) -> None:
    mock_insert_audit = MagicMock()

    with patch("audit_handler.database.AuditLog.insert", mock_insert_audit):
        handler = AuditHandler()

        handler(event_message=mock_sns_event, context=None)

        mock_insert_audit.assert_called_once()
