package main

import (
	"context"
	"time"

	"github.com/NeowayLabs/logger"
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

func main() {

	if err := render(); err != nil {
		logger.Fatal("Oh crap! %v", err)
	}
}

func render() error {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return err
	}

	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return err
	}

	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}
	defer w.Destroy()

	if err := drawTitle(r); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	s, err := newScene(r)
	if err != nil {
		return err
	}
	defer s.finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	select {
	case err := <-s.run(ctx, r):
		return err

	case <-time.After(5 * time.Second):
		return nil
	}

	if err := s.paint(r); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	return nil
}

func drawTitle(r *sdl.Renderer) error {

	r.Clear()

	f, err := ttf.OpenFont("res/fonts/Roboto-Black.ttf", 20)
	if err != nil {
		return err
	}

	defer f.Close()

	s, err := f.RenderUTF8_Solid("Runner!", sdl.Color{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	})

	if err != nil {
		return err
	}

	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return err
	}

	defer t.Destroy()

	if err := r.Copy(t, nil, nil); err != nil {
		return err
	}

	r.Present()

	return nil
}
