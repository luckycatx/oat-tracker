package view

import (
	"encoding/hex"
	"io"
	"log"
	"strconv"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/luckycatx/oat-tracker/internal/pkg/info"
)

var Paused bool

func Run(w *app.Window, trk func(io.Writer)) error {
	var running bool
	var run_btn, pause_btn widget.Clickable

	var info_list = &widget.List{List: layout.List{Axis: layout.Vertical}}
	var log_list = &widget.List{List: layout.List{Axis: layout.Vertical}}

	var logger = &logger{list: log_list}
	log.SetOutput(logger)

	var th = material.NewTheme()

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&op.Ops{}, e)

			layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd}.Layout(gtx,

				// Title
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					title := material.H2(th, "Oat-Tracker")
					title.Font.Weight = font.Bold
					title.Font.Style = font.Italic
					title.Alignment = text.Middle
					return title.Layout(gtx)
				}),

				// Tracker Info
				layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: unit.Dp(8), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return widget.Border{Color: th.Palette.ContrastBg, CornerRadius: unit.Dp(4), Width: unit.Dp(1)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {

								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Bottom: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return material.H6(th, "Tracker Info").Layout(gtx)
										})
									}),
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										infos := make([]string, 0, len(info.Infos))
										for room := range info.Infos {
											infos = append(infos, room)
										}
										return material.List(th, info_list).Layout(gtx, len(info.Infos), func(gtx layout.Context, index int) layout.Dimensions {
											room := infos[index]
											info.Update(room)
											text := "Room: " + room + "\n"
											for ih, p := range info.Infos[room] {
												ih = hex.EncodeToString([]byte(ih))
												text += "Torrent InfoHash: " + string(ih) + " | Swarm Peers: " + strconv.Itoa(p) + "\n"
											}
											text += "---"
											return material.Body2(th, text).Layout(gtx)
										})
									}),
								)
							})
						})
					})
				}),

				// Logs
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: unit.Dp(8), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return widget.Border{Color: th.Palette.ContrastBg, CornerRadius: unit.Dp(4), Width: unit.Dp(1)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return material.List(th, log_list).Layout(gtx, len(logger.content), func(gtx layout.Context, index int) layout.Dimensions {
									return material.Body2(th, logger.content[index]).Layout(gtx)
								})
							})
						})
					})
				}),

				// Buttons
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {

						return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &pause_btn, "Pause")
								if !Paused {
									btn.Text = "Pause"
								} else {
									btn.Text = "Resume"
								}
								for pause_btn.Clicked(gtx) {
									if !Paused {
										Paused = true
										log.Print("Server paused...")
									} else {
										Paused = false
										log.Print("Server resumed...")
									}
								}
								return btn.Layout(gtx)
							}),

							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &run_btn, "Run")
								btn.Background = th.Palette.ContrastBg
								if run_btn.Clicked(gtx) {
									if !running {
										running = true
										log.Print("Running the server...")
										go trk(logger)
									}
								}
								return btn.Layout(gtx)
							}),
						)
					})
				}),
			)
			e.Frame(gtx.Ops)
		}
	}

}
