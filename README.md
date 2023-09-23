# Chess Reminder

Chess.com currently nudges daily-game players only when they are 8 hours out
from their move timer running out.

This AWS Lambda sends a text message reminder 30 minutes before the expiration
of time, when players actually have to make a move and can't procrastinate.
