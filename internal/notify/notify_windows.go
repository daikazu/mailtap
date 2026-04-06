package notify

import "github.com/gen2brain/beeep"

func Send(title, body string) {
	beeep.Notify(title, body, "")
}
