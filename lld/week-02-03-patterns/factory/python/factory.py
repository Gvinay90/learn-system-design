"""Factory Method LLD — a NotificationFactory producing Email/SMS/Push
notifications behind a common Notification protocol. See ../README.md for
the design writeup.
"""
from __future__ import annotations

from enum import Enum, auto
from typing import Protocol


class NotificationType(Enum):
    EMAIL = auto()
    SMS = auto()
    PUSH = auto()


class Notification(Protocol):
    def send(self, recipient: str, message: str) -> str: ...


class EmailNotification:
    def send(self, recipient: str, message: str) -> str:
        return f"Email to {recipient}: {message}"


class SMSNotification:
    def send(self, recipient: str, message: str) -> str:
        return f"SMS to {recipient}: {message}"


class PushNotification:
    def send(self, recipient: str, message: str) -> str:
        return f"Push to {recipient}: {message}"


class UnknownNotificationTypeError(Exception):
    pass


_REGISTRY = {
    NotificationType.EMAIL: EmailNotification,
    NotificationType.SMS: SMSNotification,
    NotificationType.PUSH: PushNotification,
}


def create_notification(notification_type: NotificationType) -> Notification:
    cls = _REGISTRY.get(notification_type)
    if cls is None:
        raise UnknownNotificationTypeError(f"unknown notification type: {notification_type}")
    return cls()


def _demo() -> None:
    for nt in (NotificationType.EMAIL, NotificationType.SMS, NotificationType.PUSH):
        notification = create_notification(nt)
        print(notification.send("user@example.com", "your order has shipped"))


if __name__ == "__main__":
    _demo()
