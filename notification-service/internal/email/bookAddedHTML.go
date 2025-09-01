package email

import (
	"fmt"

	"github.com/Yarik7610/library-backend-common/broker/event"
	"github.com/Yarik7610/library-backend/notification-service/internal/utils"
)

func FormatBookAddedEmailHTML(addedBook *event.BookAdded) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<body style="font-family: Arial, sans-serif; line-height:1.5; color:#333;">
				<h2>📚 New book added!</h2>
				<p>A new book has arrived in the <b>%s</b> category:</p>
				<ul>
					<li><b>📖 Title:</b> %s</li>
					<li><b>👤 Author:</b> %s</li>
					<li><b>📅 Year:</b> %d</li>
				</ul>
				<p>Enjoy reading!</p>
			</body>
		</html>`,
		utils.Capitalize(addedBook.Category),
		addedBook.Title,
		addedBook.AuthorName,
		addedBook.Year,
	)
}
