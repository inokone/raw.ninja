import json
from typing import Any

from expiry_marker import logger, session
from expiry_marker.models import AlbumEvent, Response

class Status(str, Enum):
    SUCCESS = "SUCCESS"
    FAILED = "FAILED"


class ExpiryMarker:

    def __call__(self, event_message: Any, context: Any) -> Response:
        logger.info("Initializing expiry marker...")
        event: AlbumEvent = AlbumEvent.parse_obj(
            json.loads(event_message["Records"][0]["Sns"]["Message"])
        )
        response: Response = self._handler(event)
        return response

    def _handler(self, event: AlbumEvent) -> Response:
        with session.begin(): # transaction
            response = Response(status=Status.SUCCESS)
            self._update_expiry_for_album(session=session, event=event)
            logger.info("Responding to request with: %s", response)
            return response

    def _update_expiry_for_album(self, event: AlbumEvent):
        try:
            photos = collect_photos(event.album_id)
            expiry_in_days = get_expiration_days(event.album_id)
            recalculate_expiry_dates(photos, expiry_in_days)
        except Exception as e:
            self.logger.error("Insert failed:", exc_info=True)
            raise Exception("Failed to update exiry date") from e


handler = ExpiryMarker()
