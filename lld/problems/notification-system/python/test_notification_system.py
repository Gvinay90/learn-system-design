import pytest

from notification_system import (
    AlwaysFailChannel,
    Channel,
    EmailChannel,
    FlakyChannel,
    NotificationService,
    PushChannel,
    RetryPolicy,
    SMSChannel,
    render_template,
)


@pytest.mark.parametrize(
    "template,data,expected",
    [
        ("Hello {name}", {"name": "Alice"}, "Hello Alice"),
        (
            "Hello {name}, your order {orderId} shipped",
            {"name": "Bob", "orderId": "42"},
            "Hello Bob, your order 42 shipped",
        ),
        ("Hi {name}, code {code}", {"name": "Carl"}, "Hi Carl, code {code}"),
        ("plain message", {"unused": "x"}, "plain message"),
        ("Hello {name", {"name": "Dee"}, "Hello {name"),
        ("", {}, ""),
    ],
)
def test_render_template(template, data, expected):
    assert render_template(template, data) == expected


def test_notify_dispatches_per_preference():
    email = EmailChannel()
    sms = SMSChannel()
    push = PushChannel()

    service = NotificationService(RetryPolicy(max_attempts=1))
    service.register_channel(email)
    service.register_channel(sms)
    service.register_channel(push)

    service.set_preferences("u1", [Channel.EMAIL, Channel.SMS])
    service.set_preferences("u2", [Channel.PUSH])

    results = service.notify("u1", "u1@example.com", "Hello {name}", {"name": "Ann"})
    assert len(results) == 2
    assert all(r.success for r in results)

    assert len(email.sent()) == 1
    assert email.sent()[0].message == "Hello Ann"
    assert len(sms.sent()) == 1
    assert len(push.sent()) == 0

    service.notify("u2", "device-token", "Ping {name}", {"name": "Bo"})
    assert len(push.sent()) == 1
    assert push.sent()[0].message == "Ping Bo"


def test_notify_unknown_user_raises():
    service = NotificationService(RetryPolicy.default())
    service.register_channel(EmailChannel())

    with pytest.raises(ValueError):
        service.notify("ghost", "x", "hi", {})


def test_retry_succeeds_after_n_failures():
    underlying = EmailChannel()
    flaky = FlakyChannel(underlying, fail_count=2)

    service = NotificationService(RetryPolicy(max_attempts=3, delay_seconds=0.001))
    service.register_channel(flaky)
    service.set_preferences("u1", [Channel.EMAIL])

    results = service.notify("u1", "u1@example.com", "Order {orderId} shipped", {"orderId": "7"})
    assert len(results) == 1
    r = results[0]
    assert r.success
    assert r.attempts == 3
    assert flaky.attempts() == 3
    assert len(underlying.sent()) == 1
    assert underlying.sent()[0].message == "Order 7 shipped"


def test_retry_gives_up_after_max_attempts():
    failing = AlwaysFailChannel(Channel.SMS)

    service = NotificationService(RetryPolicy(max_attempts=4, delay_seconds=0.001))
    service.register_channel(failing)
    service.set_preferences("u1", [Channel.SMS])

    results = service.notify("u1", "5551234", "hi {name}", {"name": "X"})
    assert len(results) == 1
    r = results[0]
    assert not r.success
    assert r.attempts == 4
    assert failing.attempts() == 4


def test_notify_multiple_channels_independent_retry():
    email = EmailChannel()
    flaky_sms = FlakyChannel(SMSChannel(), fail_count=1)

    service = NotificationService(RetryPolicy(max_attempts=2, delay_seconds=0.001))
    service.register_channel(email)
    service.register_channel(flaky_sms)
    service.set_preferences("u1", [Channel.EMAIL, Channel.SMS])

    results = service.notify("u1", "recipient", "msg", {})
    assert len(results) == 2
    assert all(r.success for r in results)
    assert results[0].attempts == 1
    assert results[1].attempts == 2
