import threading

from singleton import AppConfig


def test_get_instance_returns_same_instance():
    AppConfig._reset_for_test()
    a = AppConfig.get_instance()
    b = AppConfig.get_instance()
    assert a is b
    assert a.id == b.id


def test_set_and_get():
    AppConfig._reset_for_test()
    config = AppConfig.get_instance()
    config.set("region", "us-east-1")
    assert config.get("region") == "us-east-1"
    assert config.get("missing") is None


def test_concurrent_first_access():
    AppConfig._reset_for_test()
    ids = [None] * 200

    def worker(idx: int) -> None:
        ids[idx] = AppConfig.get_instance().id

    threads = [threading.Thread(target=worker, args=(i,)) for i in range(len(ids))]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    assert len(set(ids)) == 1


def test_demo():
    AppConfig._reset_for_test()
    config = AppConfig.get_instance()
    config.set("env", "us-east-1")
    assert config.get("env") == "us-east-1"
