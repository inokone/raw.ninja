import json
from typing import Any

from sqlalchemy.orm import sessionmaker, Session

from expiry_marker import logger, app_config
from expiry_marker.models import AlbumEvent, Response
from expiry_marker.database import Photo, create_session_maker


class Status(str, Enum):
    SUCCESS = "SUCCESS"
    FAILED = "FAILED"


class ExpiryMarker:

    def __init__(self) -> None:
        self.session = sessionmaker[Session] = create_session_maker(app_config.get_connection_string())

    def _insert_event(self, session: Session, event: AlbumEvent):
        try:
            photos = collect_photos(session, event.album_id)
            recalculate_expiry_dates(session, photos)
        except Exception as e:
            self.logger.error("Insert failed:", exc_info=True)
            raise Exception("Failed to insert into audit log") from e

    def _handler(self, event: AlbumEvent) -> Response:
        with self.session() as session, session.begin():
            response = Response(status=Status.SUCCESS)
            self._insert_event(session=session, event=event)
            logger.info("Responding to request with: %s", response)
            return response

    def __call__(self, event_message: Any, context: Any) -> Response:
        logger.info("Initializing expiry marker")
        event: AlbumEvent = AlbumEvent.parse_obj(
            json.loads(event_message["Records"][0]["Sns"]["Message"])
        )
        response: Response = self._handler(event)
        return response


handler = ExpiryMarker()
