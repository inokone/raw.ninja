from typing import Optional, TypedDict, List

from pydantic import BaseModel


class AuditEvent(BaseModel):
    """Data type for auditable events in RAW.Ninja application

    Attributes
    ----------
    user_id : str
        UUID the actor user in the RAW.Ninja application
    entry_date : int
        unix epoch of the audited event
    correlation_id : Optional[str]
        UUID identifying all changes related to the audited event
    action : str
        the name of the action. e.g. "upload", "delete"
    target_ids : List[str]
        UUID for the target of the action
    target_type : str
        the type of the event target. e.g. "photo", "album"
    meta : dict
        additional metadata for the action
    outcome : str
        whether the action was successful (default "success")
    """

    user_id: str
    entry_date: int
    correlation_id: Optional[str]
    action: str
    target_ids: List[str]
    target_type: str
    meta: dict
    outcome: str = "success"


class AuditResponse(TypedDict):
    """Data type for result of audit event persistence."""

    status: str
