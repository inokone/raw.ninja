import boto3 # type: ignore
from botocore.exceptions import ClientError # type: ignore

from rawninja import logger
from rawninja.models import AuditEvent


class AuditLog: # pylint: disable=too-few-public-methods
    """ AWS DynamoDB implementation of audit log storage """

    def __init__(self, tablename: str) -> None:
        dynamodb = boto3.resource('dynamodb')
        self._table = dynamodb.Table(tablename)

    def insert(self, event: AuditEvent):
        try:
            self._table.put_item(
                Item={
                    'user': event.user_id,
                    'audit': event
                }
            )
        except ClientError as err:
            logger.error(
                "Couldn't store audit event for correlation_id [%s]." 
                "Code: [%s]." 
                "Cause: [%s].",
                event.correlation_id,
                err.response['Error']['Code'],
                err.response['Error']['Message']
            )
            raise
