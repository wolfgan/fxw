package xwidget

import (
	"math"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type dotsProgressInfinityRenderer struct {
	BaseRenderer
	dots           []*canvas.Circle
	running        bool
	progress       *DotsProgressBarInfinity
	animation      *time.Ticker
	animatedOn     int
	animatedOff    int
	animationSpeed float32
}

// MinSize calculates the minimum size of a progress bar.
// This is simply the "100%" label size plus padding.
func (r *dotsProgressInfinityRenderer) MinSize() fyne.Size {
	if r.progress.minSize.Width == 0 && r.progress.minSize.Height == 0 {
		tsize := fyne.MeasureText("100%", theme.TextSize(), fyne.TextStyle{})
		return fyne.NewSize(tsize.Width+theme.InnerPadding()*2, tsize.Height+theme.InnerPadding()*2)
	}

	return r.progress.minSize
}

func (r *dotsProgressInfinityRenderer) updateBar() {
	r.moveDots()
	canvas.Refresh(r.progress)
}

func (r *dotsProgressInfinityRenderer) moveDots() {
	r.dots[r.animatedOn].FillColor = theme.PrimaryColor()
	r.dots[r.animatedOn].Refresh()
	r.animatedOn++
	if r.animatedOn == len(r.dots) {
		r.animatedOn = 0
	}

	r.animatedOff++
	if r.animatedOff == len(r.dots) {
		r.animatedOff = 0
	}
	r.dots[r.animatedOff].FillColor = primaryFadedColor(4)
	r.dots[r.animatedOff].Refresh()

}

// Layout the components of the check xwidget
func (r *dotsProgressInfinityRenderer) Layout(size fyne.Size) {
	r.dotsResize(size)
}

func (r *dotsProgressInfinityRenderer) dotsResize(size fyne.Size) {
	centerX := float64(size.Width / 2)
	centerY := float64(size.Height / 2)
	length := fyne.Min(size.Width, size.Height) / 2
	dotSize := length * 0.1
	radius := float64(length - dotSize)

	dot := 0
	for angle := 435; angle >= 90; angle -= 15 {
		x := float32(centerX + radius*math.Cos(float64(angle)*2*math.Pi/360))
		y := float32(centerY - radius*math.Sin(float64(angle)*2*math.Pi/360))
		r.dots[dot].Position1 = fyne.Position{X: x - dotSize, Y: y - dotSize}
		r.dots[dot].Position2 = fyne.Position{X: x + dotSize, Y: y + dotSize}
		dot++
	}
}

// applyTheme updates the progress bar to match the current theme
func (r *dotsProgressInfinityRenderer) applyTheme() {
	for _, dot := range r.dots {
		dot.FillColor = primaryFadedColor(2)
	}
}

func (r *dotsProgressInfinityRenderer) Refresh() {
	if r.isRunning() {
		return // we refresh from the goroutine
	}
	r.applyTheme()

	for _, dot := range r.dots {
		dot.Refresh()
	}

	canvas.Refresh(r.progress)
}

func (r *dotsProgressInfinityRenderer) isRunning() bool {
	r.progress.propertyLock.RLock()
	defer r.progress.propertyLock.RUnlock()

	return r.running
}

// Start the infinite progress bar background thread to update it continuously
func (r *dotsProgressInfinityRenderer) start() {
	if r.isRunning() {
		return
	}

	r.progress.propertyLock.Lock()
	defer r.progress.propertyLock.Unlock()

	r.running = true
	r.animatedOn = 0
	r.animatedOff = len(r.dots) - 4
	r.animation = time.NewTicker(time.Duration(int64(time.Millisecond) * int64(r.animationSpeed*100)))
	go func() {
		for range r.animation.C {
			r.updateBar()
		}
	}()
}

// Stop the background thread from updating the infinite progress bar
func (r *dotsProgressInfinityRenderer) stop() {
	r.progress.propertyLock.Lock()
	defer r.progress.propertyLock.Unlock()

	r.running = false
	r.animation.Stop()
}

func (r *dotsProgressInfinityRenderer) Destroy() {
	r.stop()
}

// DotsProgressBarInfinity xwidget creates a horizontal panel that indicates progress
type DotsProgressBarInfinity struct {
	widget.BaseWidget
	propertyLock sync.RWMutex
	minSize      fyne.Size
	renderer     fyne.WidgetRenderer
}

// Show this xwidget, if it was previously hidden
func (p *DotsProgressBarInfinity) Show() {
	p.Start()
	p.BaseWidget.Show()
}

// Hide this xwidget, if it was previously visible
func (p *DotsProgressBarInfinity) Hide() {
	p.Stop()
	p.BaseWidget.Hide()
}

// Start the infinite progress bar animation
func (p *DotsProgressBarInfinity) Start() {
	p.renderer.(*dotsProgressInfinityRenderer).start()
}

// Stop the infinite progress bar animation
func (p *DotsProgressBarInfinity) Stop() {
	p.renderer.(*dotsProgressInfinityRenderer).stop()
}

// Running returns the current state of the infinite progress animation
func (p *DotsProgressBarInfinity) Running() bool {
	return p.renderer.(*dotsProgressInfinityRenderer).isRunning()
}

func (p *DotsProgressBarInfinity) SetSpeed(speed float32) {
	switch {
	case speed < 0.1:
		speed = 0.1
	case speed > 1.0:
		speed = 1.0
	}
	p.renderer.(*dotsProgressInfinityRenderer).animationSpeed = speed
	p.renderer.(*dotsProgressInfinityRenderer).animation.Reset(time.Duration(int64(time.Millisecond) * int64(speed*100)))

	p.Refresh()
}

func (p *DotsProgressBarInfinity) Speed() float32 {
	return p.renderer.(*dotsProgressInfinityRenderer).animationSpeed
}

// MinSize returns the size that this xwidget should not shrink below
func (p *DotsProgressBarInfinity) MinSize() fyne.Size {
	p.ExtendBaseWidget(p)
	return p.BaseWidget.MinSize()
}

func (p *DotsProgressBarInfinity) SetMinSize(size fyne.Size) {
	p.minSize = size
	p.Refresh()
}

// CreateRenderer is a private method to Fyne which links this xwidget to its renderer
func (p *DotsProgressBarInfinity) CreateRenderer() fyne.WidgetRenderer {
	p.ExtendBaseWidget(p)

	var dots []*canvas.Circle
	var objects []fyne.CanvasObject

	for angle := 435; angle >= 90; angle -= 15 {
		dot := canvas.NewCircle(primaryFadedColor(4))
		dot.StrokeWidth = 1
		dot.StrokeColor = primaryFadedColor(2)

		dots = append(dots, dot)
		objects = append(objects, dot)
	}

	p.renderer = &dotsProgressInfinityRenderer{
		BaseRenderer:   NewBaseRenderer(objects),
		dots:           dots,
		progress:       p,
		animationSpeed: 0.5,
	}
	p.renderer.(*dotsProgressInfinityRenderer).start()

	return p.renderer
}

// NewDotsProgressBarInfinity creates a new progress bar xwidget.
// The default Min is 0 and Max is 1, Values set should be between those numbers.
// The display will convert this to a percentage.
func NewDotsProgressBarInfinity() *DotsProgressBarInfinity {
	p := &DotsProgressBarInfinity{minSize: fyne.Size{0, 0}}
	p.ExtendBaseWidget(p)
	return p
}
