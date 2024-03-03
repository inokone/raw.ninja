from typing import List, Optional

from pydantic import BaseModel


class AlbumEvent(BaseModel):
    user_id: str
    action: str
    album_id: str
    correlation_id: Optional[str]
    photos: List[str]


class Response(BaseModel):
    status: str
