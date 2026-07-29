package facade

import "testing"

func TestWatchMovieCoordinatesSubsystems(t *testing.T) {
	theater := NewHomeTheaterFacade()
	theater.WatchMovie("Inception")

	if !theater.Lights.Dimmed {
		t.Fatal("expected lights dimmed")
	}
	if !theater.Screen.Down {
		t.Fatal("expected screen down")
	}
	if !theater.Projector.On || theater.Projector.Input != "dvd" {
		t.Fatal("expected projector on with dvd input")
	}
	if !theater.Amp.On || theater.Amp.Volume != 7 {
		t.Fatal("expected amp on at volume 7")
	}
	if theater.Dvd.Movie != "Inception" {
		t.Fatalf("expected movie playing, got %q", theater.Dvd.Movie)
	}
}

func TestEndMovieResetsSubsystems(t *testing.T) {
	theater := NewHomeTheaterFacade()
	theater.WatchMovie("Inception")
	theater.EndMovie()

	if theater.Dvd.Movie != "" || theater.Dvd.On {
		t.Fatal("expected dvd stopped and off")
	}
	if theater.Amp.On {
		t.Fatal("expected amp off")
	}
	if theater.Projector.On {
		t.Fatal("expected projector off")
	}
	if theater.Screen.Down {
		t.Fatal("expected screen raised")
	}
	if theater.Lights.Dimmed {
		t.Fatal("expected lights brightened")
	}
}

func TestIndependentTheaters(t *testing.T) {
	t1 := NewHomeTheaterFacade()
	t2 := NewHomeTheaterFacade()

	t1.WatchMovie("Dune")
	if t2.Dvd.Movie != "" {
		t.Fatal("expected second theater unaffected by first")
	}
}
