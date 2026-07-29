"""Notification System LLD — Python reference implementation.

Multi-channel notification dispatch via the Strategy pattern, per-user
channel preferences, retry-with-fixed-delay on send failure, and basic
"{placeholder}" template rendering. See ../README.md for the design writeup.
"""
from __future__ import annotations

import threading
import time
from dataclasses import dataclass, field
from enum import Enum
from typing import Dict, List, Optional, Protocol


class Channel(Enum):
    EMAIL = "EMAIL"
    SMS = "SMS"
    PUSH = "PUSH"


class SendFailedError(Exception):
    """Raised by a channel implementation to simulate a delivery failure."""


class Notifier(Protocol):
    """Strategy interface every delivery channel implements."""

    def channel(self) -> Channel: ...

    def send(self, recipient: str, message: str) -> None: ...


@dataclass(frozen=True)
class SentMessage:
    recipient: str
    message: str


class _RecordingChannel:
    """Shared base for the in-memory simulated channels below."""

    def __init__(self) -> None:
        self._lock = threading.Lock()
        self._sent: List[SentMessage] = []

    def send(self, recipient: str, message: str) -> None:
        with self._lock:
            self._sent.append(SentMessage(recipient, message))

    def sent(self) -> List[SentMessage]:
        with self._lock:
            return list(self._sent)


class EmailChannel(_RecordingChannel):
    """Simulates sending email by recording messages in memory."""

    def channel(self) -> Channel:
        return Channel.EMAIL


class SMSChannel(_RecordingChannel):
    """Simulates sending SMS by recording messages in memory."""

    def channel(self) -> Channel:
        return Channel.SMS


class PushChannel(_RecordingChannel):
    """Simulates sending a push notification by recording messages in memory."""

    def channel(self) -> Channel:
        return Channel.PUSH


class FlakyChannel:
    """Wraps a delegate Notifier and fails the first N sends before
    delegating for real. Exists purely so tests can exercise the retry
    mechanism deterministically.
    """

    def __init__(self, delegate: Notifier, fail_count: int) -> None:
        self._delegate = delegate
        self._failures_remaining = fail_count
        self._lock = threading.Lock()
        self._attempts = 0

    def channel(self) -> Channel:
        return self._delegate.channel()

    def send(self, recipient: str, message: str) -> None:
        with self._lock:
            self._attempts += 1
            if self._failures_remaining > 0:
                self._failures_remaining -= 1
                raise SendFailedError("simulated flaky failure")
        self._delegate.send(recipient, message)

    def attempts(self) -> int:
        with self._lock:
            return self._attempts


class AlwaysFailChannel:
    """Always fails; useful for exercising "retries exhausted" behavior."""

    def __init__(self, channel_type: Channel) -> None:
        self._channel_type = channel_type
        self._lock = threading.Lock()
        self._attempts = 0

    def channel(self) -> Channel:
        return self._channel_type

    def send(self, recipient: str, message: str) -> None:
        with self._lock:
            self._attempts += 1
        raise SendFailedError("simulated permanent failure")

    def attempts(self) -> int:
        with self._lock:
            return self._attempts


def render_template(template: str, data: Optional[Dict[str, str]] = None) -> str:
    """Substitutes "{key}" tokens in template with data[key]. Tokens with no
    matching key are left untouched.
    """
    data = data or {}
    out: List[str] = []
    i = 0
    n = len(template)
    while i < n:
        open_idx = template.find("{", i)
        if open_idx == -1:
            out.append(template[i:])
            break
        close_idx = template.find("}", open_idx)
        if close_idx == -1:
            out.append(template[i:])
            break
        out.append(template[i:open_idx])
        key = template[open_idx + 1:close_idx]
        if key in data:
            out.append(data[key])
        else:
            out.append(template[open_idx:close_idx + 1])
        i = close_idx + 1
    return "".join(out)


@dataclass(frozen=True)
class RetryPolicy:
    """Configures the retry-on-failure wrapper around a channel send."""

    max_attempts: int = 3  # total attempts, including the first; must be >= 1
    delay_seconds: float = 0.01  # fixed delay between attempts

    @staticmethod
    def default() -> "RetryPolicy":
        return RetryPolicy(max_attempts=3, delay_seconds=0.01)


@dataclass
class SendResult:
    """Captures the outcome of dispatching to a single channel."""

    channel: Channel
    attempts: int
    error: Optional[Exception] = None

    @property
    def success(self) -> bool:
        return self.error is None


class NotificationService:
    """Dispatches rendered notifications to a user's preferred channels,
    retrying each channel independently on failure.
    """

    def __init__(self, retry: RetryPolicy = RetryPolicy()):
        self._retry = retry
        self._channels: Dict[Channel, Notifier] = {}
        self._preferences: Dict[str, List[Channel]] = {}
        self._lock = threading.Lock()

    def register_channel(self, notifier: Notifier) -> None:
        """Makes a Notifier available for dispatch under its own channel() identity."""
        with self._lock:
            self._channels[notifier.channel()] = notifier

    def set_preferences(self, user_id: str, channels: List[Channel]) -> None:
        """Records which channels a given user wants notifications on, in order."""
        with self._lock:
            self._preferences[user_id] = list(channels)

    def notify(
        self,
        user_id: str,
        recipient: str,
        template: str,
        data: Optional[Dict[str, str]] = None,
    ) -> List[SendResult]:
        """Renders template with data and sends the result to every channel
        preferred by user_id, retrying each channel per the service's
        RetryPolicy. Returns one SendResult per attempted channel and does
        not stop early if one channel ultimately fails.
        """
        with self._lock:
            prefs = self._preferences.get(user_id)
            if prefs is None:
                raise ValueError(f"no channel preferences registered for user {user_id!r}")
            prefs = list(prefs)

        message = render_template(template, data)

        results: List[SendResult] = []
        for ch in prefs:
            with self._lock:
                notifier = self._channels.get(ch)
            if notifier is None:
                results.append(SendResult(ch, 0, SendFailedError(f"no channel registered for {ch}")))
                continue
            results.append(self._send_with_retry(notifier, recipient, message))
        return results

    def _send_with_retry(self, notifier: Notifier, recipient: str, message: str) -> SendResult:
        max_attempts = max(self._retry.max_attempts, 1)

        last_error: Optional[Exception] = None
        for attempt in range(1, max_attempts + 1):
            try:
                notifier.send(recipient, message)
                return SendResult(notifier.channel(), attempt, None)
            except SendFailedError as e:
                last_error = e
                if attempt < max_attempts:
                    time.sleep(self._retry.delay_seconds)
        return SendResult(notifier.channel(), max_attempts, last_error)


def _demo() -> None:
    email = EmailChannel()
    sms = SMSChannel()
    push = PushChannel()

    service = NotificationService(RetryPolicy.default())
    service.register_channel(email)
    service.register_channel(sms)
    service.register_channel(push)

    service.set_preferences("u1", [Channel.EMAIL, Channel.SMS])

    results = service.notify(
        "u1", "ann@example.com", "Hello {name}, your order {orderId} shipped",
        {"name": "Ann", "orderId": "42"},
    )
    for r in results:
        print(f"{r.channel} -> success={r.success} attempts={r.attempts}")
    print(f"email sent: {len(email.sent())}")
    print(f"sms sent: {len(sms.sent())}")


if __name__ == "__main__":
    _demo()
