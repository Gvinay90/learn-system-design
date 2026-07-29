from facade import HomeTheaterFacade


def test_watch_movie_coordinates_subsystems():
    theater = HomeTheaterFacade()
    theater.watch_movie("Inception")

    assert theater.lights.dimmed
    assert theater.screen.down
    assert theater.projector.on
    assert theater.projector.input == "dvd"
    assert theater.amp.on
    assert theater.amp.volume == 7
    assert theater.dvd.movie == "Inception"


def test_end_movie_resets_subsystems():
    theater = HomeTheaterFacade()
    theater.watch_movie("Inception")
    theater.end_movie()

    assert theater.dvd.movie == ""
    assert not theater.dvd.on
    assert not theater.amp.on
    assert not theater.projector.on
    assert not theater.screen.down
    assert not theater.lights.dimmed


def test_independent_theaters():
    t1 = HomeTheaterFacade()
    t2 = HomeTheaterFacade()

    t1.watch_movie("Dune")
    assert t2.dvd.movie == ""
