from pydantic import BaseModel

class AlbumEvent(BaseModel):
    action: str
    album_id: str
    target_type: str
