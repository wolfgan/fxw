package xwidget

import (
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type dotsProgressRenderer struct {
	BaseRenderer
	dots     []*canvas.Circle
	label    *canvas.Text
	progress *DotsProgressBar
}

// MinSize calculates the minimum size of a progress bar.
// This is simply the "100%" label size plus padding.
func (r *dotsProgressRenderer) MinSize() fyne.Size {
	if r.progress.minSize.Width == 0 && r.progress.minSize.Height == 0 {
		var tsize fyne.Size
		if text := r.progress.TextFormatter; text != nil {
			tsize = fyne.MeasureText(text(), r.label.TextSize, r.label.TextStyle)
		} else {
			tsize = fyne.MeasureText("100%", r.label.TextSize, r.label.TextStyle)
		}

		return fyne.NewSize(tsize.Width+theme.InnerPadding()*2, tsize.Height+theme.InnerPadding()*2)
	}

	return r.progress.minSize
}

func (r *dotsProgressRenderer) updateBar() {
	if r.progress.Value < r.progress.Min {
		r.progress.Value = r.progress.Min
	}
	if r.progress.Value > r.progress.Max {
		r.progress.Value = r.progress.Max
	}

	delta := float32(r.progress.Max - r.progress.Min)
	ratio := float32(r.progress.Value-r.progress.Min) / delta

	if text := r.progress.TextFormatter; text != nil {
		r.label.Text = text()
	} else {
		r.label.Text = strconv.Itoa(int(ratio*100)) + "%"
	}

	dotNum := int(math.Floor(float64(float32(len(r.dots)) * ratio)))
	for i := 0; i < len(r.dots); i++ {
		if i < dotNum {
			r.dots[i].FillColor = theme.PrimaryColor()
		} else {
			r.dots[i].FillColor = primaryFadedColor(4)
		}
	}
}

// Layout the components of the check xwidget
func (r *dotsProgressRenderer) Layout(size fyne.Size) {
	r.dotsResize(size)
	r.label.TextSize = fyne.Min(size.Width, size.Height) * 0.15
	r.label.Resize(size)
	r.updateBar()
}

func (r *dotsProgressRenderer) dotsResize(size fyne.Size) {
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

// Refresh
func (r *dotsProgressRenderer) Refresh() {
	size := r.progress.Size()
	r.label.Color = theme.PrimaryColor()
	r.label.TextSize = fyne.Min(size.Width, size.Height) * 0.15
	r.label.Refresh()
	r.updateBar()

	for _, dot := range r.dots {
		dot.Refresh()
	}
}

// DotsProgressBar the widget creates a circle of dots indicating progress
type DotsProgressBar struct {
	widget.BaseWidget

	minSize fyne.Size

	Min, Max, Value float64

	// TextFormatter can be used to have a custom format of progress text.
	// If set, it overrides the percentage readout and runs each time the value updates.
	//
	// Since: 1.4
	TextFormatter func() string `json:"-"`

	binder basicBinder
}

// Bind connects the specified data source to this DotsProgressBar.
// The current value will be displayed and any changes in the data will cause the widget to update.
//
// Since: 2.0
func (p *DotsProgressBar) Bind(data binding.Float) {
	p.binder.SetCallback(p.updateFromData)
	p.binder.Bind(data)
}

// SetValue changes the current value of this progress bar (from p.Min to p.Max).
// The widget will be refreshed to indicate the change.
func (p *DotsProgressBar) SetValue(v float64) {
	p.Value = v
	p.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (p *DotsProgressBar) MinSize() fyne.Size {
	p.ExtendBaseWidget(p)
	return p.BaseWidget.MinSize()
}

func (p *DotsProgressBar) SetMinSize(size fyne.Size) {
	p.minSize = size
	p.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (p *DotsProgressBar) CreateRenderer() fyne.WidgetRenderer {
	p.ExtendBaseWidget(p)
	if p.Min == 0 && p.Max == 0 {
		p.Max = 1.0
	}

	size := p.Size()
	label := canvas.NewText("0%", theme.PrimaryColor())
	label.TextSize = fyne.Min(size.Width, size.Height) * 0.15
	label.Alignment = fyne.TextAlignCenter

	var dots []*canvas.Circle
	var objects []fyne.CanvasObject

	for angle := 435; angle >= 90; angle -= 15 {
		dot := canvas.NewCircle(primaryFadedColor(4))
		dot.StrokeWidth = 1
		dot.StrokeColor = primaryFadedColor(2)

		dots = append(dots, dot)
		objects = append(objects, dot)
	}
	objects = append(objects, label)

	renderer := dotsProgressRenderer{
		BaseRenderer: NewBaseRenderer(objects),
		dots:         dots,
		label:        label,
		progress:     p,
	}

	return &renderer
}

// Unbind disconnects any configured data source from this ProgressBar.
// The current value will remain at the last value of the data source.
//
// Since: 2.0
func (p *DotsProgressBar) Unbind() {
	p.binder.Unbind()
}

// NewDotsProgressBar creates a new progress bar xwidget.
// The default Min is 0 and Max is 1, Values set should be between those numbers.
// The display will convert this to a percentage.
func NewDotsProgressBar() *DotsProgressBar {
	p := &DotsProgressBar{Min: 0, Max: 1, minSize: fyne.Size{0, 0}}
	p.ExtendBaseWidget(p)
	return p
}

// NewDotsProgressBarWithData returns a progress bar connected with the specified data source.
//
// Since: 2.0
func NewDotsProgressBarWithData(data binding.Float) *DotsProgressBar {
	p := NewDotsProgressBar()
	p.Bind(data)

	return p
}

func (p *DotsProgressBar) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	floatSource, ok := data.(binding.Float)
	if !ok {
		return
	}

	val, err := floatSource.Get()
	if err != nil {
		fyne.LogError("Error getting current data value", err)
		return
	}
	p.SetValue(val)
}
