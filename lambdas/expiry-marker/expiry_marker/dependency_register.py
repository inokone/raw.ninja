import inject

from rawninja.utils.config import AppConfig
from sqlalchemy.orm import sessionmaker, Session
from sqlalchemy import create_engine

import logging
from trumonitor.log import TMOAdapter
from rawninja import initialize_logger
from rawninja.database import AuditLog


def register_dependencies() -> None:
    inject.configure(_configure)


def _configure(binder: inject.Binder) -> None:
    app_config = AppConfig()  # type:  ignore[call-args]
    binder.bind(AppConfig, app_config)
    initialize_logger(log_level=app_config.LOG_LEVEL.upper())
    tmo_adapter = TMOAdapter(logging.getLogger("fusion.object_storage.audit"))
    binder.bind(TMOAdapter, tmo_adapter)
    binder.bind_to_provider(
        sessionmaker[Session], _create_session_maker(app_config.get_connection_string())
    )
    binder.bind(AuditLog, AuditLog)


def _create_session_maker(db_connection_string):
    engine = create_engine(db_connection_string)

    def create_session():
        _session = sessionmaker(engine)
        return _session

    return create_session
