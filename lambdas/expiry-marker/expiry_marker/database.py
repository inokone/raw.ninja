import datetime
from enum import Enum
from typing import Dict, List

from sqlalchemy.orm import sessionmaker, Session
from sqlalchemy import create_engine
from sqlalchemy import String, func, DateTime, JSON, Integer
from sqlalchemy.orm import DeclarativeBase, Mapped, Table, Column, ForeignKey
from sqlalchemy.orm import mapped_column, relationship


# Key for albums in collections table
ALBUM_COLLECTION_TYPE = "ALBUM"


def create_session_maker(db_connection_string):
    engine = create_engine(db_connection_string)

    def create_session():
        _session = sessionmaker(engine)
        return _session

    return create_session


class Action(Enum):
    # Delete is an action to delete a photo from storage
    DELETE = 1
    # MoveTo is an action to move photo to the storage specified by the "Action target"
    MOVE_TO = 2


class ActionTarget(Enum):
    # StandardStorage is a target of an action, planned to count 1x from quota
    STANDARD_STORAGE = 1
    # FrozenStorage is a target of an action, planned to count 0.5x from quota
    FROZEN_STORAGE = 2
    # Bin is storage for files marked as deleted (not sure this is needed at all)
    BIN = 3


class Base(DeclarativeBase):
    pass


rule_association_table = Table(
    "ruleset_rules",
    Base.metadata,
    Column("rule_id", ForeignKey("rules.id"), primary_key=True),
    Column("ruleset_id", ForeignKey("rulesets.id"), primary_key=True),
)


class Rule(Base):
    __tablename__ = "rules"
    id: Mapped[str] = mapped_column(String(36), primary_key=True)
    timing: Mapped[int] = mapped_column(Integer)
    action_id: Mapped[int] = mapped_column(Integer, index=True)
    target_id: Mapped[int] = mapped_column(Integer, index=True)
    created_at: Mapped[datetime.datetime] = mapped_column(DateTime, default=func.now())
    deleted_at: Mapped[datetime.datetime] = mapped_column(DateTime)


class RuleSet(Base):
    __tablename__ = "rulesets"
    id: Mapped[str] = mapped_column(String(36), primary_key=True)
    rules: Mapped[List["Rule"]] = relationship(
        secondary=rule_association_table, back_populates="rulesets"
    )
    created_at: Mapped[datetime.datetime] = mapped_column(DateTime, default=func.now())
    deleted_at: Mapped[datetime.datetime] = mapped_column(DateTime)


album_association_table = Table(
    "collection_photos",
    Base.metadata,
    Column("collection_id", ForeignKey("collection.id"), primary_key=True),
    Column("photo_id", ForeignKey("photos.id"), primary_key=True),
)


class Album(Base):
    __tablename__ = "collections"
    id: Mapped[str] = mapped_column(String(36), primary_key=True)
    type: Mapped[str] = mapped_column(String(256), primary_key=True)
    photos: Mapped[List["Photo"]] = relationship(
        secondary=album_association_table, back_populates="albums"
    )
    target_id: Mapped[str] = mapped_column(String, index=True)
    created_at: Mapped[datetime.datetime] = mapped_column(DateTime, default=func.now())
    deleted_at: Mapped[datetime.datetime] = mapped_column(DateTime)


class Photo(Base):
    __tablename__ = "photos"
    id: Mapped[str] = mapped_column(String(36), primary_key=True)
    action: Mapped[str] = mapped_column(String(20))
    target_id: Mapped[str] = mapped_column(String(20), index=True)
    target_type: Mapped[str] = mapped_column(String(20), index=True)
    outcome: Mapped[str] = mapped_column(String(20))
    meta: Mapped[Dict] = mapped_column(JSON())
    created: Mapped[datetime.datetime] = mapped_column(DateTime, default=func.now())

    @staticmethod
    def get_all(session, limit=10) -> List["Photo"]:
        rows = session.query(Photo).limit(limit).all()
        return rows
