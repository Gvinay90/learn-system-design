// Package facade demonstrates the Facade pattern: a HomeTheaterFacade
// coordinates five independent subsystems (amplifier, DVD player, projector,
// screen, lights) behind two simple calls, WatchMovie and EndMovie.
package facade

type Amplifier struct {
	On     bool
	Volume int
}

func (a *Amplifier) TurnOn()         { a.On = true }
func (a *Amplifier) TurnOff()        { a.On = false; a.Volume = 0 }
func (a *Amplifier) SetVolume(v int) { a.Volume = v }

type DvdPlayer struct {
	On    bool
	Movie string
}

func (d *DvdPlayer) TurnOn()           { d.On = true }
func (d *DvdPlayer) TurnOff()          { d.On = false; d.Movie = "" }
func (d *DvdPlayer) Play(movie string) { d.Movie = movie }
func (d *DvdPlayer) Stop()             { d.Movie = "" }

type Projector struct {
	On    bool
	Input string
}

func (p *Projector) TurnOn()               { p.On = true }
func (p *Projector) TurnOff()              { p.On = false; p.Input = "" }
func (p *Projector) SetInput(input string) { p.Input = input }

type Screen struct {
	Down bool
}

func (s *Screen) Lower() { s.Down = true }
func (s *Screen) Raise() { s.Down = false }

type Lights struct {
	Dimmed bool
}

func (l *Lights) Dim()      { l.Dimmed = true }
func (l *Lights) Brighten() { l.Dimmed = false }

// HomeTheaterFacade hides the coordination order of five subsystems behind
// two calls, so a client never needs to know the correct on/off sequence.
type HomeTheaterFacade struct {
	Amp       *Amplifier
	Dvd       *DvdPlayer
	Projector *Projector
	Screen    *Screen
	Lights    *Lights
}

func NewHomeTheaterFacade() *HomeTheaterFacade {
	return &HomeTheaterFacade{
		Amp:       &Amplifier{},
		Dvd:       &DvdPlayer{},
		Projector: &Projector{},
		Screen:    &Screen{},
		Lights:    &Lights{},
	}
}

func (f *HomeTheaterFacade) WatchMovie(movie string) {
	f.Lights.Dim()
	f.Screen.Lower()
	f.Projector.TurnOn()
	f.Projector.SetInput("dvd")
	f.Amp.TurnOn()
	f.Amp.SetVolume(7)
	f.Dvd.TurnOn()
	f.Dvd.Play(movie)
}

func (f *HomeTheaterFacade) EndMovie() {
	f.Dvd.Stop()
	f.Dvd.TurnOff()
	f.Amp.TurnOff()
	f.Projector.TurnOff()
	f.Screen.Raise()
	f.Lights.Brighten()
}
