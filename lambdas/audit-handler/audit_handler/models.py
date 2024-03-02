from typing import Optional, TypedDict

from pydantic import BaseModel


class AuditEvent(BaseModel):
    """Data type for auditable events in RAW.Ninja application

    Attributes
    ----------
    correlation_id : str
        UUID identifying all changes related to the audited event
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
    entry_date : int
        unix epoch of the audited event
    outcome : str
        whether the action was successful (default "success")
    """

    correlation_id: str
    user_id: Optional[str]
    action: str
    target_id: str
    target_type: str
    meta: dict
    entry_date: int
    outcome: str = "success"


class AuditResponse(TypedDict):
    """Data type for result of audit event persistence."""

    status: str
