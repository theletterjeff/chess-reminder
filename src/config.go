package main

type chessReminderCfg struct {
	reminderMins []int
}

func newChessReminderCfg() *chessReminderCfg {
	return &chessReminderCfg{
		reminderMins: []int{30},
	}
}
