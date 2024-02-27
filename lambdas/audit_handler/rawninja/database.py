import boto3
from boto3 import ClientError

from rawninja import logger
from rawninja.models import AuditEvent


class AuditLog:
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
                f"Couldn't store audit event for correlation_id [{event.correlation_id}]. 
                Code: [{err.response["Error"]["Code"]}] 
                Cause: [{err.response["Error"]["Message"]}]",
            )
            raise
