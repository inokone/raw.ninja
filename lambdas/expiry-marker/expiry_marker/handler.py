import inject
from sqlalchemy.orm import sessionmaker, Session


class Status(str, Enum):
    SUCCESS = "SUCCESS"
    FAILED = "FAILED"


class ExpiryMarker:
    session: sessionmaker[Session] = inject.attr(sessionmaker[Session])
    audit_log = inject.attr(AuditLog)

    def _insert_event(self, session: Session, event: AuditEvent):
        try:
            self.audit_log.insert(
                session=session,
                action=event.action,
                target_id=event.target_id,
                target_type=event.target_type,
                outcome=event.outcome,
                meta=event.meta,
            )
        except Exception as e:
            self.logger.error("Insert failed:", exc_info=True)
            raise Exception("Failed to insert into audit log") from e

    def _handler(self, event: AuditEvent) -> AuditResponse:
        with self.session() as session, session.begin():
            response = AuditResponse(
                status=Status.SUCCESS
            )
            self._insert_event(session=session, event=event)
            self.logger.info("Responding to request with: %s", response)
            return response


    def __call__(self, event_message: AuditSnsEvent, contex: Any) -> AuditResponse:
        event: AuditEvent = AuditEvent.parse_obj(
            json.loads(event_message["Records"][0]["Sns"]["Message"])
        )
        response: AuditResponse = self._handler(event)
        return response