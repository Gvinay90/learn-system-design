from observer import LoggingDisplay, PhoneDisplay, WeatherStation


def test_all_observers_notified_on_change():
    station = WeatherStation()
    phone = PhoneDisplay()
    logger = LoggingDisplay()
    station.subscribe(phone)
    station.subscribe(logger)

    station.set_measurements(21.5, 60)

    assert phone.last_temp_c == 21.5
    assert phone.last_humidity == 60
    assert len(logger.history) == 1


def test_multiple_updates_accumulate():
    station = WeatherStation()
    phone = PhoneDisplay()
    station.subscribe(phone)

    station.set_measurements(20, 50)
    station.set_measurements(22, 55)

    assert phone.update_count == 2
    assert phone.last_temp_c == 22


def test_unsubscribe_stops_notifications():
    station = WeatherStation()
    phone = PhoneDisplay()
    station.subscribe(phone)
    station.set_measurements(20, 50)

    station.unsubscribe(phone)
    station.set_measurements(99, 99)

    assert phone.last_temp_c == 20
    assert phone.update_count == 1


def test_no_observers_does_not_raise():
    station = WeatherStation()
    station.set_measurements(10, 10)
