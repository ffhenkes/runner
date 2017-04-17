package main

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	t       int
	bg      *sdl.Texture
	runners []*sdl.Texture
}

func newScene(r *sdl.Renderer) (*scene, error) {

	bg, err := img.LoadTexture(r, "res/img/background.png")
	if err != nil {
		return nil, err
	}

	var runners []*sdl.Texture
	for i := 1; i <= 6; i++ {
		path := fmt.Sprintf("res/img/frame_%d.png", i)
		human, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, err
		}
		runners = append(runners, human)
	}

	return &scene{bg: bg, runners: runners}, nil
}

func (s *scene) run(ctx context.Context, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		for range time.Tick(10 * time.Millisecond) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) paint(r *sdl.Renderer) error {
	s.t++

	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return err
	}

	rect := &sdl.Rect{X: 10, Y: 300 - 240/2, W: 176, H: 240}

	i := s.t / 10 % len(s.runners)
	if err := r.Copy(s.runners[i], nil, rect); err != nil {
		return err
	}

	r.Present()
	return nil
}

func (s *scene) finish() {
	s.bg.Destroy()
}
