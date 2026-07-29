import pytest

from factory import NotificationType, UnknownNotificationTypeError, create_notification


def test_create_email_notification():
    n = create_notification(NotificationType.EMAIL)
    assert n.send("alice@example.com", "hello") == "Email to alice@example.com: hello"


def test_create_sms_notification():
    n = create_notification(NotificationType.SMS)
    assert n.send("+1-555-0100", "hello") == "SMS to +1-555-0100: hello"


def test_create_push_notification():
    n = create_notification(NotificationType.PUSH)
    assert n.send("device-123", "hello") == "Push to device-123: hello"


def test_create_unknown_notification_type():
    with pytest.raises(UnknownNotificationTypeError):
        create_notification("bogus")


def test_demo():
    for nt in (NotificationType.EMAIL, NotificationType.SMS, NotificationType.PUSH):
        n = create_notification(nt)
        assert n.send("user", "test") != ""
