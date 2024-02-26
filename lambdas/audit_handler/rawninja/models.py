from typing import Optional, TypedDict, Any

from pydantic import BaseModel


class MessageBody(TypedDict):
    """ Data type for content of a message """

    Message: str


class Record(TypedDict):
    """ Data type for individual messages in an event from SNS """

    Sns: MessageBody


class SnsEvent(TypedDict):
    """ Data type for events coming from AWS SNS """

    Records: list[Record]


class AuditEvent(BaseModel):
    """ Data type for auditable events in RAW.Ninja application

    Attributes
    ----------
    user_id : Optional[str]
        UUID the actor user in the RAW.Ninja application
    action : str
        the name of the action. e.g. "upload", "delete"
    target_id : str
        UUID for the target of the action
    target_type : str
        the type of the event target. e.g. "photo", "album"
    meta : dict
        additional metadata for the action
    outcome : str
        whether the action was successful (default "success")
    """

    user_id: Optional[str]
    action: str
    target_id: str
    target_type: str
    meta: dict
    outcome: str = "success"

    def add_metadata(self, key: str, value: Any):
        self.meta[key] = value


class AuditResponse(TypedDict):
    """ Data type for result of audit event persistence. """

    status: str
