package utility

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

const pushProgressBarWidth = 24

type pushProgressReporter struct {
	w          io.Writer
	isTTY      bool
	total      int64
	lastDone   int64
	lastPhase  string
	finalized  bool
	started    bool
}

func newPushProgressReporter(w io.Writer, isTTY bool) *pushProgressReporter {
	if w == nil {
		w = io.Discard
	}

	return &pushProgressReporter{w: w, isTTY: isTTY}
}

func (p *pushProgressReporter) Start(total int64) {
	if p.finalized {
		return
	}

	if total <= 0 {
		total = 100
	}

	p.total = total
	p.lastDone = 0
	p.lastPhase = "starting"
	p.started = true
	p.render(0, p.lastPhase)
}

func (p *pushProgressReporter) Update(done int64, phase string) {
	if p.finalized {
		return
	}

	if !p.started {
		p.Start(100)
	}

	if done < 0 {
		done = 0
	}
	if done > p.total {
		done = p.total
	}
	if phase == "" {
		phase = p.lastPhase
	}

	if !p.isTTY && done == p.lastDone && phase == p.lastPhase {
		return
	}

	p.lastDone = done
	p.lastPhase = phase
	p.render(done, phase)
}

func (p *pushProgressReporter) Complete(msg string) {
	if p.finalized {
		return
	}

	if !p.started {
		p.Start(100)
	}

	p.Update(p.total, "complete")
	if p.isTTY {
		fmt.Fprintln(p.w)
	}
	if msg != "" {
		fmt.Fprintln(p.w, msg)
	}
	p.finalized = true
}

func (p *pushProgressReporter) Fail(err error) {
	if p.finalized {
		return
	}

	if p.isTTY {
		fmt.Fprintln(p.w)
	}
	if err != nil {
		fmt.Fprintf(p.w, "push failed: %v\n", err)
	} else {
		fmt.Fprintln(p.w, "push failed")
	}
	p.finalized = true
}

func (p *pushProgressReporter) render(done int64, phase string) {
	percent := int64(0)
	if p.total > 0 {
		percent = done * 100 / p.total
	}

	if p.isTTY {
		filled := int(done * pushProgressBarWidth / p.total)
		if filled < 0 {
			filled = 0
		}
		if filled > pushProgressBarWidth {
			filled = pushProgressBarWidth
		}
		bar := strings.Repeat("=", filled) + strings.Repeat("-", pushProgressBarWidth-filled)
		fmt.Fprintf(p.w, "\r[%s] %3d%% %s", bar, percent, phase)
		return
	}

	fmt.Fprintf(p.w, "push progress: %3d%% %s\n", percent, phase)
}

func isInteractiveTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}