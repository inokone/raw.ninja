import boto3  # type: ignore
from botocore.exceptions import ClientError  # type: ignore

from audit_handler import logger
from audit_handler.models import AuditEvent


class AuditPersistenceError(Exception):
    pass


class AuditLog:  # pylint: disable=too-few-public-methods
    """AWS DynamoDB implementation of audit log storage"""

    def __init__(self, region: str, tablename: str) -> None:
        dynamodb = boto3.resource("dynamodb", region_name=region)
        self._table = dynamodb.Table(tablename)

    def insert(self, event: AuditEvent):
        try:
            self._table.put_item(
                Item={
                    "user_id": event.user_id,
                    "entry_date": event.entry_date,
                    "audit": event.model_dump(),
                }
            )
        except ClientError as err:
            logger.error(
                "Couldn't store audit event for correlation_id [%s]. "
                "Code: [%s]. "
                "Cause: [%s]. ",
                event.correlation_id,
                err.response["Error"]["Code"],
                err.response["Error"]["Message"],
            )
            raise AuditPersistenceError("failed to insert audit event") from err
