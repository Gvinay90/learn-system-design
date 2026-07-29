import threading

from payment_gateway import (
    FakePaymentProcessor,
    PaymentGateway,
    PaymentRequest,
    PaymentStatus,
    RetryPolicy,
)


def make_request(key: str) -> PaymentRequest:
    return PaymentRequest(
        idempotency_key=key,
        payer_id="payer-1",
        payee_id="payee-1",
        amount=100.0,
        currency="INR",
    )


def new_gateway(processor: FakePaymentProcessor, max_attempts: int = 3) -> PaymentGateway:
    return PaymentGateway(processor, RetryPolicy(max_attempts=max_attempts, delay_seconds=0.001))


def test_happy_path_charge_records_ledger_entry():
    gateway = new_gateway(FakePaymentProcessor())
    result = gateway.charge(make_request("key-1"))

    assert result.status == PaymentStatus.SUCCESS
    entries = gateway.ledger.entries()
    assert len(entries) == 1
    assert entries[0].payment_id == result.id
    assert entries[0].amount == 100.0


def test_same_idempotency_key_does_not_reprocess():
    processor = FakePaymentProcessor()
    gateway = new_gateway(processor)

    first = gateway.charge(make_request("key-2"))
    second = gateway.charge(make_request("key-2"))

    assert processor.call_count == 1
    assert first.id == second.id
    assert first.status == second.status


def test_retry_policy_succeeds_after_transient_failures():
    processor = FakePaymentProcessor(fail_times=2)
    gateway = new_gateway(processor)

    result = gateway.charge(make_request("key-3"))
    assert result.status == PaymentStatus.SUCCESS
    assert len(result.attempts) == 3


def test_retry_policy_exhausts_to_failed():
    processor = FakePaymentProcessor(always_fail=True)
    gateway = new_gateway(processor)

    result = gateway.charge(make_request("key-4"))
    assert result.status == PaymentStatus.FAILED
    assert len(result.attempts) == 3

    again = gateway.charge(make_request("key-4"))
    assert again.status == PaymentStatus.FAILED
    assert processor.call_count == 3


def test_concurrent_charges_same_key_process_once():
    processor = FakePaymentProcessor()
    gateway = new_gateway(processor)

    n = 20
    results = [None] * n

    def worker(idx: int) -> None:
        results[idx] = gateway.charge(make_request("shared-key"))

    threads = [threading.Thread(target=worker, args=(i,)) for i in range(n)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    assert processor.call_count == 1
    first_id = results[0].id
    for r in results:
        assert r.id == first_id
        assert r.status == PaymentStatus.SUCCESS
